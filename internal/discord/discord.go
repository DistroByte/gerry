package discord

import (
	"log/slog"
)

func SearchGuildByChannelID(textChannelID string) (guildID string) {
	channel, _ := Session.Channel(textChannelID)
	guildID = channel.GuildID
	return guildID
}

func SendDiscordMessage(channelID string, message string) {
	_, err := Session.ChannelMessageSend(channelID, message)
	if err != nil {
		slog.Error("failed to send message", "error", err)
	}
}
