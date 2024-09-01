package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/distrobyte/gerry/internal/config"
	"github.com/distrobyte/gerry/internal/handlers"
	"github.com/distrobyte/gerry/internal/models"
	"github.com/rs/zerolog/log"
)

var DiscordSession *discordgo.Session

func InitSession() {
	var err error
	DiscordSession, err = discordgo.New("Bot " + config.GetDiscordToken())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create discord session")
	}

	DiscordSession.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMessageReactions |
		discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessageTyping |
		discordgo.IntentsGuildVoiceStates |
		discordgo.IntentsDirectMessages

	DiscordSession.State.MaxMessageCount = 200
	DiscordSession.State.TrackChannels = true
	DiscordSession.State.TrackThreads = true
	DiscordSession.State.TrackEmojis = true
	DiscordSession.State.TrackMembers = true
	DiscordSession.State.TrackThreadMembers = true
	DiscordSession.State.TrackRoles = true
	DiscordSession.State.TrackVoice = true
	DiscordSession.State.TrackPresences = true
}

func InitConnection() {
	if err := DiscordSession.Open(); err != nil {
		log.Fatal().Err(err).Msg("failed to create websocket connection to discord")
		return
	}
}

func ReadyHandler(s *discordgo.Session, event *discordgo.Ready) {
	err := s.UpdateListeningStatus(config.GetBotStatus())
	if err != nil {
		log.Warn().Err(err).Msg("failed to update game status")
	}

	log.Info().Msg("connected to discord")
}

func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	message := &models.Message{
		Content:    m.Content,
		Author:     m.Author.Username,
		Channel:    m.ChannelID,
		ID:         m.ID,
		RecievedAt: time.Now(),
		Platform:   "discord",
	}

	response, err := handlers.HandleMessage(message)
	if err == nil && response != "" {
		SendMessage(message.Channel, response)
		log.Info().
			Str("platform", message.Platform).
			Str("event", "message").
			Str("content", message.Content).
			Str("author", message.Author).
			Str("channel", message.Channel).
			Str("id", message.ID).
			TimeDiff("duration", time.Now(), message.RecievedAt).
			Msg("handled message")
	} else if err != nil {
		log.Error().Err(err).Msg("failed to handle message")
	}
}

func MessageReactHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	message, err := s.State.Message(m.ChannelID, m.MessageID)
	if err == discordgo.ErrStateNotFound {
		message, err = s.ChannelMessage(m.ChannelID, m.MessageID)
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to get message")
		return
	}

	log.Info().
		Str("platform", "discord").
		Str("event", "reaction").
		Str("message_id", m.MessageID).
		Str("user_id", m.UserID).
		Str("emoji", m.Emoji.Name).
		Str("message_content", message.Content).
		Str("emoji", m.Emoji.Name).
		Msg("reaction added")

	if m.UserID == s.State.User.ID {
		return
	}
}
