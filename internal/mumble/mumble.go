package mumble

import (
	"log/slog"
	"strings"
)

func SendMumbleMessage(channelID uint32, message string) {
	if message == "" {
		slog.Debug("message is empty")
		return
	}

	channel := MumbleSession.Channels[channelID]
	if channel == nil {
		slog.Warn("channel not found", "channel", channelID)
		return
	}

	if strings.Contains(message, "\n") {
		message = strings.ReplaceAll(message, "\n", "<br>")
	}

	channel.Send(message, false)
}
