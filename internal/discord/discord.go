package discord

import (
	"github.com/rs/zerolog/log"
)

func SearchGuildByChannelID(textChannelID string) (guildID string) {
	channel, _ := DiscordSession.Channel(textChannelID)
	guildID = channel.GuildID
	return guildID
}

func SendDiscordMessage(channelID string, message string) {
	_, err := DiscordSession.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Error().Err(err).Msg("failed to send message")
	}
}
