package gerryfrank

import (
	"math/rand"
	"time"

	"git.dbyte.xyz/distro/gerry/bot"
	"git.dbyte.xyz/distro/gerry/shared"
	"git.dbyte.xyz/distro/gerry/utils"
)

var chain utils.MarkovChainImpl

func Setup(bot bot.Bot, _chain utils.MarkovChainImpl) {
	bot.Register("gerryfrank", Call)
	chain = _chain
}

func Call(bot bot.Bot, context bot.MessageContext, arguments []string, config shared.Config) error {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	firstWord := chain.ChooseFirstWord()

	var str = firstWord
	var prev = firstWord

	for ok := true; ok; ok = prev[len(prev)-1] != '!' && prev[len(prev)-1] != '.' && len(str) < 1024 {
		prev = chain.ChooseNextWord(prev)
		if prev == "EOF" {
			break
		}
		str += " " + prev
	}

	bot.Send(context, str)
	return nil
}
