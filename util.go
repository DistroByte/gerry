package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"git.dbyte.xyz/distro/gerry/shared"
	"git.dbyte.xyz/distro/gerry/symbols"
	"github.com/bwmarrin/discordgo"
	"github.com/google/shlex"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"gopkg.in/fsnotify.v1"
)

// MustEnv returns the value of the environment variable named by the key, or
// logs a fatal error if the variable is not set.
func MustEnv(key string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	log.Fatalf("please provide a %q", key)
	return ""
}

// LoadPlugins loads all plugins from the given paths.
func LoadPlugins(discordConn *discordgo.Session, pluginPaths []string, plugins map[string]*plugin) error {
	for _, path := range pluginPaths {
		source, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read plugin source: %w", err)
		}
		if err := loadPlugin(discordConn, plugins, source); err != nil {
			return fmt.Errorf("load plugin %q: %w", path, err)
		}
	}
	return nil
}

var plugmu sync.Mutex

// loadPlugin loads a single plugin from the given source. It will replace the
// existing plugin if it already exists.
func loadPlugin(discordConn *discordgo.Session, plugins map[string]*plugin, source []byte) error {
	plugmu.Lock()
	defer plugmu.Unlock()

	pluginter := interp.New(interp.Options{
		Env:    os.Environ(),
		GoPath: build.Default.GOPATH,
	})
	_ = pluginter.Use(stdlib.Symbols)
	_ = pluginter.Use(symbols.Symbols)

	prog, err := pluginter.Compile(string(source))
	if err != nil {
		return fmt.Errorf("compile file: %w", err)
	}
	if _, err := pluginter.Execute(prog); err != nil {
		return fmt.Errorf("execute file: %w", err)
	}

	pkg := prog.PackageName()

	// quit old plugin if we're reloading
	if plug, ok := plugins[pkg]; ok && plug.stopCh != nil {
		close(plug.stopCh)
	}
	delete(plugins, pkg)

	log.Printf("loading plugin %s", pkg)

	var plugin plugin
	plugin.name = pkg
	plugin.bot = &Bot{discordConn, &plugin}

	if setup, _ := getFunc[shared.PluginSetupFunc](pluginter, pkg, "Setup"); setup != nil {
		setup(plugin.bot)
	}

	if run, _ := getFunc[shared.PluginRunFunc](pluginter, pkg, "Run"); run != nil {
		plugin.stopCh = make(chan struct{})
		go run(plugin.bot, plugin.stopCh)
	}

	plugins[pkg] = &plugin

	return nil
}

// getFunc returns a function of the given type from the given package.
func getFunc[T any](i *interp.Interpreter, pkg string, key string) (T, error) {
	key = fmt.Sprintf("%s.%s", pkg, key)

	var zero T
	rv, err := i.Eval(key)
	if err != nil {
		return zero, fmt.Errorf("finding an exported `%s` function: %w", key, err)
	}
	funcv, ok := rv.Interface().(T)
	if !ok {
		return zero, fmt.Errorf("exported `%s` func is `%T` not `%T`", key, rv.Interface(), zero)
	}
	return funcv, nil
}

// AddWatchers adds watchers for the given plugin paths. This allows for live
// reloading of plugins.
func AddWatchers(discordConn *discordgo.Session, watcher *fsnotify.Watcher, plugins map[string]*plugin, pluginPaths []string) {
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Printf("modified file: %s", event.Name)
					source, err := os.ReadFile(event.Name)
					if err != nil {
						log.Printf("error reading file: %v", err)
					}
					if err := loadPlugin(discordConn, plugins, source); err != nil {
						log.Printf("error loading plugin: %v", err)
					}
				}
			case err := <-watcher.Errors:
				log.Printf("error: %v", err)
			}
		}
	}()

	for _, path := range pluginPaths {
		if err := watcher.Add(filepath.Dir(path)); err != nil {
			log.Panicf("error adding watcher for %q: %v", path, err)
		}
	}
}

func PluginFromCommand(plugins map[string]*plugin, command string) *plugin {
	for _, p := range plugins {
		if _, ok := p.commands[command]; ok {
			return p
		}
	}
	return nil
}

func PluginCommands(plugins map[string]*plugin) []string {
	var r []string
	for _, p := range plugins {
		for command := range p.commands {
			r = append(r, command)
		}
	}
	sort.Strings(r)
	return r
}

func ParseCommand(message string) (string, []string) {
	parts, _ := shlex.Split(message)
	if len(parts) == 0 {
		return "", nil
	}

	command, ok := strings.CutPrefix(parts[0], "!")
	if !ok {
		return "", nil
	}
	return command, parts[1:]
}
