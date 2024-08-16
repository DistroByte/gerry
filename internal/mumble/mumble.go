package mumble

import "log/slog"

func SendMumbleMessage(channelID uint32, message string) {
	if message == "" {
		slog.Error("cannot send an empty message")
		return
	}

	channel := MumbleSession.Channels[channelID]
	if channel == nil {
		slog.Warn("channel not found", "channel", channelID)
		return
	}

	channel.Send(message, false)
}
