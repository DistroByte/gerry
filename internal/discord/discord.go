package discord

import (
	"log/slog"
)

func SearchGuildByChannelID(textChannelID string) (guildID string) {
	channel, _ := DiscordSession.Channel(textChannelID)
	guildID = channel.GuildID
	return guildID
}

func SendDiscordMessage(channelID string, message string) {
	if message == "" {
		slog.Error("cannot send an empty message")
		return
	}

	_, err := DiscordSession.ChannelMessageSend(channelID, message)
	if err != nil {
		slog.Error("failed to send message", "error", err)
	}
}
