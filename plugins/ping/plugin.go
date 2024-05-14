package ping

import (
	"git.dbyte.xyz/distro/gerry/bot"
	"git.dbyte.xyz/distro/gerry/shared"
	"git.dbyte.xyz/distro/gerry/utils"
)

func Setup(bot bot.Bot, _ utils.MarkovChainImpl) {
	bot.Register("ping", Call)
}

func Call(bot bot.Bot, context bot.MessageContext, arguments []string, config shared.Config) error {
	bot.Send(context, "pong")
	return nil
}
