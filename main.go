package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"git.dbyte.xyz/distro/gerry/shared"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/fsnotify.v1"
)

var (
	cfgToken             = MustEnv("GERRY_TOKEN")
	cfgPluginsPath       = MustEnv("GERRY_PLUGINS_PATH")
	cfgWatcherEnabled, _ = strconv.ParseBool(MustEnv("GERRY_WATCHER_ENABLED"))
	// cfgLogsPath    = MustEnv("GERRY_LOGS_PATH")
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// log date and time
	log.SetFlags(0)

	rand.New(rand.NewSource(time.Now().UnixNano()))
	println("Starting Gerry...")

	// create discord client
	discord, err := discordgo.New("Bot " + cfgToken)
	if err != nil {
		log.Panicf("error creating discord client: %v", err)
	}

	// set intents
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	// open websocket connection
	log.Println("opening discord connection")
	if err := discord.Open(); err != nil {
		log.Panicf("error opening connection: %v", err)
	}

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
	if err := LoadPlugins(discord, pluginPaths, plugins); err != nil {
		log.Panicf("error loading plugins: %v", err)
	}

	// watch for changes in plugins
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panicf("error adding watcher for plugins dir")
	}

	// add watchers for plugins
	if cfgWatcherEnabled {
		AddWatchers(discord, watcher, plugins, pluginPaths)
	}

	// add message handler function
	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		args := strings.Split(strings.ToLower(m.Content), " ")

		// small chance to send a message from the gerryfrank plugin
		if rand.Intn(100) < 2 {
			plugin := PluginFromCommand(plugins, "gerryfrank")
			if plugin != nil {
				context := shared.MessageContext{
					Sender:  m.Author.Username,
					Target:  m.ChannelID,
					Source:  "discord",
					Message: m.Content,
				}
				if err := plugin.commands["gerryfrank"](plugin.bot, context, args, conf); err != nil {
					s.ChannelMessageSend(m.ChannelID, err.Error())
				}
			}
		}

		command, arguments := ParseCommand(m.Content)
		if command == "" {
			return
		}

		if command == "load" {
			if len(arguments) == 0 {
				s.ChannelMessageSend(m.ChannelID, "please provide a plugin name")
				return
			}

			pluginName := arguments[0]
			pluginPath := filepath.Join(conf.PluginPath, pluginName, "plugin.go")
			source, err := os.ReadFile(pluginPath)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("error reading plugin: %v", err))
				return
			}

			if err := LoadPlugin(discord, plugins, source); err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}

			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("loaded plugin %q", pluginName))
			return
		}

		plugin := PluginFromCommand(plugins, command)

		if plugin == nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("commands are: %v", strings.Join(PluginCommands(plugins), ", ")))
			return
		}
		context := shared.MessageContext{
			Sender:  m.Author.Username,
			Target:  m.ChannelID,
			Source:  "discord",
			Message: m.Content,
		}

		if err := plugin.commands[command](plugin.bot, context, arguments, conf); err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-exit
		log.Println("shutting down...")

		log.Println("closing discord connection")
		if err := discord.Close(); err != nil { // close discord connection
			log.Printf("error closing connection: %v", err)
		}

		for _, plug := range plugins { // stop plugins
			if plug.stopCh != nil {
				close(plug.stopCh)
			}
		}
	}()

	// wait for all go routines to finish
	wg.Wait()
	log.Println("shutdown complete")
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
	*discordgo.Session
	*plugin
}

// Send sends a message to the given recipient. It is platform agnostic.
func (c *Bot) Send(context shared.MessageContext, message string) {
	log.Printf("sending message to %s - %q: %q", context.Source, context.Target, message)
	if context.Source == "discord" {
		_, err := c.ChannelMessageSend(context.Target, message)
		if err != nil {
			log.Printf("error sending message: %v", err)
		}
	}
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
