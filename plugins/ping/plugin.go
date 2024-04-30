package ping

import "git.dbyte.xyz/distro/gerry/shared"

func Setup(bot shared.Bot) {
	bot.Register("ping", Call)
}

func Call(bot shared.Bot, context shared.MessageContext, arguments []string, config shared.Config) error {
	bot.Send(context.Source, context.Target, "pong")
	return nil
}
