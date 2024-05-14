package bot

import (
	"git.dbyte.xyz/distro/gerry/shared"
	"git.dbyte.xyz/distro/gerry/utils"
)

type Bot interface {
	Register(command string, f PluginCallFunc)

	Send(context MessageContext, message string)
}

type MessageContext struct {
	Sender  string
	Target  string
	Source  string
	Message string
}

type PluginSetupFunc = func(bot Bot, chain utils.MarkovChainImpl)
type PluginCallFunc = func(bot Bot, context MessageContext, arguments []string, config shared.Config) error
type PluginRunFunc = func(bot Bot, stop <-chan struct{})
