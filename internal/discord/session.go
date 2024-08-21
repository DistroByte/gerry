package discord

import (
	"github.com/DistroByte/gerry/internal/config"
	"github.com/DistroByte/gerry/internal/handlers"
	"github.com/DistroByte/gerry/internal/models"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var DiscordSession *discordgo.Session

func InitDiscordSession() {
	var err error
	DiscordSession, err = discordgo.New("Bot " + config.GetDiscordToken())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create discord session")
	}

	DiscordSession.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessageTyping |
		discordgo.IntentsGuildVoiceStates |
		discordgo.IntentsDirectMessages
}

func InitDiscordConnection() {
	if err := DiscordSession.Open(); err != nil {
		log.Fatal().Err(err).Msg("failed to create websocket connection to discord")
		return
	}
}

func DiscordReadyHandler(s *discordgo.Session, event *discordgo.Ready) {
	err := s.UpdateListeningStatus(config.GetBotStatus())
	if err != nil {
		log.Warn().Err(err).Msg("failed to update game status")
	}

	log.Info().Msg("connected to discord")
}

func DiscordMessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	message := &models.Message{
		Content:  m.Content,
		Author:   m.Author.Username,
		Channel:  m.ChannelID,
		Platform: "discord",
	}

	response, err := handlers.HandleMessage(message)
	if err == nil && response != "" {
		SendDiscordMessage(message.Channel, response)
	}
}
