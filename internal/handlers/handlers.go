package handlers

import (
	"strings"
	"syscall"

	"github.com/distrobyte/gerry/internal/commands"
	"github.com/distrobyte/gerry/internal/config"
	"github.com/distrobyte/gerry/internal/models"
	"github.com/google/shlex"
	"github.com/rs/zerolog/log"
)

func HandleMessage(message *models.Message) (string, error) {

	var cmd string
	var args []string

	prefix := config.GetBotPrefix()
	args, err := shlex.Split(message.Content)

	if err != nil {
		if err.Error() == "EOF found when expecting closing quote" {
			args = strings.Fields(message.Content)
		} else {
			log.Error().Err(err).Msg("failed to split message")
			return "", err
		}
	}

	if len(args) == 0 {
		return "", nil
	}

	cmd, ok := strings.CutPrefix(args[0], prefix)
	if !ok {
		return "", nil
	}

	args = args[1:]

	switch cmd {
	case "ping":
		return commands.PingCommand(), nil

	case "echo":
		return commands.EchoCommand(args), nil

	case "karting":
		return commands.KartingCommand(args, *message), nil

	case "uptime":
		return commands.UptimeCommand(), nil

	case "version":
		return commands.VersionCommand(args), nil

	case "shutdown":
		config.ShutdownChannel <- syscall.SIGINT
		return "Shutting down...", nil

	default:
		break
	}

	return "", nil
}

func HandleReaction(message *models.Message) (string, error) {
	return "", nil
}
