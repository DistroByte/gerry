package commands

import (
	"git.dbyte.xyz/distro/gerry/bot"
)

type Commands struct {
	*bot.Bot
}

func (c *Commands) Register(name string, cmd bot.Command) {
	c.Bot.CommandList[name] = cmd
}

func (c *Commands) GetCommandsList() map[string]bot.Command {
	return c.Bot.CommandList
}

func (c *Commands) LoadCommands() {
	c.Register("echo", &EchoCommand{})
	c.Register("gerryfrank", &GerryFrankCommand{})

	// start run functions for commands that have them
	for _, cmd := range c.GetCommandsList() {
		if function, ok := cmd.(interface {
			Run(bot *bot.Bot, stopCh <-chan struct{})
		}); ok {
			go function.Run(c.Bot, c.Bot.StopCh)
		}
	}
}
