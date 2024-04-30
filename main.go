package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"git.dbyte.xyz/distro/gerry/shared"
	"gopkg.in/fsnotify.v1"
)

var (
	// cfgToken       = mustEnv("GERRY_TOKEN")
	cfgPluginsPath = MustEnv("GERRY_PLUGINS_PATH")
	// cfgLogsPath    = mustEnv("GERRY_LOGS_PATH")
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	log.SetFlags(0)

	rand.New(rand.NewSource(time.Now().UnixNano()))
	println("Starting Gerry...")

	conf := shared.Config{
		PluginPath: cfgPluginsPath,
	}

	// find all plugins
	pluginPaths, err := filepath.Glob(filepath.Join(conf.PluginPath, "*", "plugin.go"))
	if err != nil {
		log.Panicf("error finding plugins in plugins dir")
	}

	// load plugins
	plugins := map[string]*plugin{}
	if err := LoadPlugins(pluginPaths, plugins); err != nil {
		log.Panicf("error loading plugins: %v", err)
	}

	// watch for changes in plugins
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panicf("error adding watcher for plugins dir")
	}

	AddWatchers(watcher, plugins, pluginPaths)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-exit
		for _, plug := range plugins { // stop plugins
			if plug.stopCh != nil {
				close(plug.stopCh)
			}
		}
	}()

	wg.Wait()
	log.Println("shutdown complete")
	// <-exit
}

// plugin represents a plugin. It contains the name of the plugin, the bot
// interface, a map of commands, and a stop channel.
type plugin struct {
	name     string
	bot      shared.Bot
	commands map[string]shared.PluginCallFunc
	stopCh   chan struct{}
}

// Bot represents a bot. It contains a pointer to the plugin.
type Bot struct {
	*plugin
}

// Send sends a message to the given recipient. It is platform agnostic.
func (c *Bot) Send(to string, message string) {
	log.Printf("sending message to %q: %q", to, message)
}

// Register registers a command with the bot. The command is the string that
// triggers the command, and the function is the function that is called when
// the command is triggered.
func (c *Bot) Register(command string, f shared.PluginCallFunc) {
	if c.commands == nil {
		c.commands = map[string]shared.PluginCallFunc{}
	}
	c.commands[command] = f
}
