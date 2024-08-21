package handlers

import (
	"syscall"

	"github.com/DistroByte/gerry/internal/commands"
	"github.com/DistroByte/gerry/internal/config"
	"github.com/DistroByte/gerry/internal/models"
	"github.com/google/shlex"
	"github.com/rs/zerolog/log"
)

func HandleMessage(message *models.Message) (string, error) {
	log.Info().
		Str("platform", message.Platform).
		Str("content", message.Content).
		Str("author", message.Author).
		Str("channel", message.Channel).
		Msg("received message")

	var cmd string
	var args []string

	prefix := config.GetBotPrefix()
	args, err := shlex.Split(message.Content)

	if len(args) > 1 {
		cmd, args = args[0], args[1:]

	} else if len(args) == 1 {
		cmd = args[0]

	} else {
		return "", nil
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to split message")
		return "", err
	}

	if string(cmd[0:len(prefix)]) != prefix {
		return "", nil
	}

	cmd = cmd[len(prefix):]

	switch cmd {
	case "ping":
		return commands.PingCommand(), nil

	case "echo":
		return commands.EchoCommand(args), nil

	case "karting":
		switch args[0] {
		case "race":
			return commands.KartingRaceCommand(args, *message), nil

		case "stats":
			return commands.KartingStatsCommand(args, *message), nil

		default:
			return commands.KartingCommand(args, *message), nil
		}

	case "uptime":
		return commands.UptimeCommand(), nil

	case "shutdown":
		config.ShutdownChannel <- syscall.SIGINT
		return "Shutting down...", nil

	default:
		break
	}

	return "", nil
}
