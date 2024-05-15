package gerryfrank

import (
	"math/rand"
	randV2 "math/rand/v2"
	"time"

	"git.dbyte.xyz/distro/gerry/bot"
	"git.dbyte.xyz/distro/gerry/shared"
	"git.dbyte.xyz/distro/gerry/utils"
)

var chain utils.MarkovChainImpl

func Setup(bot bot.Bot, _chain utils.MarkovChainImpl) {
	bot.Register("gerryfrank", Call)
	chain = _chain
	
	chain.LoadModel("./markovChain.json")
	chain.LoadModel("./word_freq_map.json")
}

func Run(bot bot.Bot, stopCh <-chan struct{}) {
	// At most we'll lose the last 4 seconds of markov history
	discFlushTimer := time.NewTicker(5 * time.Second)
	gerryGravityTimer := time.NewTicker(1 * time.Hour)

	go func() {
		for {
			select {
			case <-discFlushTimer.C:
				// Replace this with the actual function to flush the Markov chain history
				chain.Flush()
			case <-gerryGravityTimer.C:
				chain.LoadModel("./word_freq_map.json")
			case <-stopCh:
				discFlushTimer.Stop()
				discFlushTimer.Stop()
			}
		}
	}()
}

func Call(bot bot.Bot, context bot.MessageContext, arguments []string, config shared.Config) error {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	firstWord := chain.ChooseFirstWord()

	var str = firstWord
	var prev = firstWord

	terminals := map[string]struct{}{"!": {}, ".": {}, "?": {}}

	// for ok := true; ok; ok = prev[len(prev)-1] != '!' && prev[len(prev)-1] != '.' && len(str) < 1024 {
	for ok := true; ok; ok = terminals[string(prev[len(prev)-1])] == struct{}{} || randV2.IntN(3) == 0 {
		prev = chain.ChooseNextWord(prev)
		if prev == "EOF" {
			break
		}
		str += " " + prev
	}

	bot.Send(context, str)
	return nil
}
