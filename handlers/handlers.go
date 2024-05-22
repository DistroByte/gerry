package handlers

import (
	"fmt"
	"log"
	"strings"

	"git.dbyte.xyz/distro/gerry/bot"
)

type Handlers struct {
	*bot.Bot
}

func handleMessage(b *bot.Bot, msg bot.Message) error {
	splitString := strings.Split(msg.GetContent(), " ")
	cmd, args := splitString[0], splitString[1:]
	if cmd == "" {
		b.MarkovChainImpl.ImportMessage(args)
		return nil
	}

	if strings.HasPrefix(cmd, b.Cfg.Prefix) {
		cmd = strings.TrimPrefix(cmd, b.Cfg.Prefix)
	} else {
		return nil
	}

	log.Printf("command: %s", cmd)
	command, ok := b.CommandList[cmd]
	if !ok {
		return nil
	}

	err := command.Execute(b, msg, args)
	if err != nil {
		return fmt.Errorf("error executing command: %w", err)
	}

	return nil
}
