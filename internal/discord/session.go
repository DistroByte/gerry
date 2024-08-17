package discord

import (
	"log/slog"

	"github.com/DistroByte/gerry/internal/config"
	"github.com/DistroByte/gerry/internal/handlers"
	"github.com/DistroByte/gerry/internal/models"
	"github.com/bwmarrin/discordgo"
)

var DiscordSession *discordgo.Session

func InitDiscordSession() {
	var err error
	DiscordSession, err = discordgo.New("Bot " + config.GetDiscordToken())
	if err != nil {
		slog.Error("failed to create discord session", "error", err)
	}

	DiscordSession.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildMessageTyping | discordgo.IntentsGuildVoiceStates
}

func InitDiscordConnection() {
	if err := DiscordSession.Open(); err != nil {
		slog.Error("failed to create websocket connection to discord", "error", err)
		return
	}
}

func DiscordReadyHandler(s *discordgo.Session, event *discordgo.Ready) {
	err := s.UpdateListeningStatus(config.GetBotStatus())
	if err != nil {
		slog.Warn("failed to update game status", "error", err)
	}

	slog.Info("connected to discord")
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
	if err == nil {
		SendDiscordMessage(message.Channel, response)
	}
}
