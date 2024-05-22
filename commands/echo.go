package commands

import (
	"strings"

	"git.dbyte.xyz/distro/gerry/bot"
)

type EchoCommand struct{}

func (c *EchoCommand) Execute(bot *bot.Bot, msg bot.Message, args []string) error {
	return bot.SendMessage(msg, strings.Join(args, " "))
}

func (c *EchoCommand) Help() string {
	return "Echoes the provided message"
}

func (c *EchoCommand) Name() string {
	return "echo"
}

func (c *EchoCommand) HasRun() bool {
	return false
}

func (c *EchoCommand) Run(bot *bot.Bot, stopCh <-chan struct{}) {}
