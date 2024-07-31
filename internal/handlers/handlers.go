package handlers

import (
	"log/slog"

	"github.com/DistroByte/gerry/internal/commands"
	"github.com/DistroByte/gerry/internal/config"
	"github.com/DistroByte/gerry/internal/models"
	"github.com/google/shlex"
)

const help string = "Current commands are:\n"

func HandleMessage(message *models.Message) (string, error) {
	slog.Info("processing message", "platform", "discord", "content", message.Content, "author", message.Author, "channel", message.Channel)

	prefix := config.GetBotPrefix()
	cmd, err := shlex.Split(message.Content)
	if err != nil {
		slog.Error("failed to split command", "error", err)
		return "", err
	}

	switch cmd[0] {
	case prefix + "ping":
		return commands.PingCommand(), nil

	case prefix + "echo":
		return commands.EchoCommand(cmd), nil

	default:
		return "", nil
	}
}
