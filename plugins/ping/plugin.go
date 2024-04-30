package ping

import (
	"git.dbyte.xyz/distro/gerry/bot"
	"git.dbyte.xyz/distro/gerry/shared"
)

func Setup(bot bot.Bot) {
	bot.Register("ping", Call)
}

func Call(bot bot.Bot, context bot.MessageContext, arguments []string, config shared.Config) error {
	bot.Send(context, "pong")
	return nil
}
