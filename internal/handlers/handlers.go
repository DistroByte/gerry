package handlers

import (
	"log/slog"
	"os"
	"syscall"

	"github.com/DistroByte/gerry/internal/commands"
	"github.com/DistroByte/gerry/internal/config"
	"github.com/DistroByte/gerry/internal/models"
	"github.com/google/shlex"
)

func HandleMessage(message *models.Message) (string, error) {
	slog.Info("processing message", "platform", message.Platform, "content", message.Content, "author", message.Author, "channel", message.Channel)

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
		slog.Error("failed to split message", "error", err)
		return "", err
	}

	if string(cmd[0:len(prefix)]) != prefix {
		return "", nil
	}

	switch cmd {
	case prefix + "ping":
		return commands.PingCommand(), nil

	case prefix + "echo":
		return commands.EchoCommand(args), nil

	case prefix + "karting":
		return commands.KartingCommand(args, *message), nil

	case prefix + "shutdown":
		config.ShutdownChannel <- os.Signal(syscall.SIGTERM)
		return "Shutting down...", nil

	default:
		break
	}

	return "", nil
}
