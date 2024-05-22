package commands

import (
	"math/rand"
	randV2 "math/rand/v2"
	"time"

	"git.dbyte.xyz/distro/gerry/bot"
	"git.dbyte.xyz/distro/gerry/utils"
)

var chain utils.MarkovChainImpl

type GerryFrankCommand struct{}

func (c *GerryFrankCommand) Execute(bot *bot.Bot, msg bot.Message, args []string) error {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	firstWord := bot.MarkovChainImpl.ChooseFirstWord()

	var str = firstWord
	var prev = firstWord

	terminals := map[string]struct{}{"!": {}, ".": {}, "?": {}}

	for ok := true; ok; ok = terminals[string(prev[len(prev)-1])] == struct{}{} || randV2.IntN(3) == 0 {
		prev = chain.ChooseNextWord(prev)
		if prev == "EOF" {
			break
		}
		str += " " + prev
	}

	return bot.SendMessage(msg, str)
}

func (c *GerryFrankCommand) Help() string {
	return "Generate a sentence in the style of Gerry Hannon vs Frank McCourt on the Late Late Show"
}

func (c *GerryFrankCommand) Name() string {
	return "gerryfrank"
}

func (c *GerryFrankCommand) Run(bot *bot.Bot, stopCh <-chan struct{}) {
	discFlushTimer := time.NewTicker(time.Duration(bot.Cfg.Markov.FlushInterval) * time.Second)
	gerryGravityTimer := time.NewTicker(time.Duration(bot.Cfg.Markov.Gravity) * time.Hour)

	go func() {
		for {
			select {
			case <-discFlushTimer.C:
				bot.MarkovChainImpl.Flush()
			case <-gerryGravityTimer.C:
				bot.MarkovChainImpl.LoadModel("./word_freq_map.json")
			case <-stopCh:
				discFlushTimer.Stop()
				discFlushTimer.Stop()
			}
		}
	}()
}
