package discord

import (
	"log/slog"

	"github.com/DistroByte/gerry/internal/config"
	"github.com/bwmarrin/discordgo"
	"github.com/google/shlex"
)

var Session *discordgo.Session

func InitSession() {
	var err error
	Session, err = discordgo.New("Bot " + config.GetDiscordToken())
	if err != nil {
		slog.Error("failed to create discord session", "error", err)
	}

	Session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildMessageTyping | discordgo.IntentsGuildVoiceStates
}

func InitConnection() {
	if err := Session.Open(); err != nil {
		slog.Error("failed to create websocket connection to discord", "error", err)
		return
	}
}

func DiscordReadyHandler(s *discordgo.Session, event *discordgo.Ready) {
	err := s.UpdateGameStatus(0, config.GetBotStatus())
	if err != nil {
		slog.Warn("failed to update game status", "error", err)
	}
}

func DiscordMessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	slog.Info("processing command", "message", m.Content, "author", m.Author.Username, "channel", m.ChannelID)

	prefix := config.GetBotPrefix()
	guildID := SearchGuildByChannelID(m.ChannelID)
	cmd, err := shlex.Split(m.Content)
	if err != nil {
		slog.Error("failed to split command", "error", err)
		return
	}

	switch cmd[0] {
	case prefix + "ping":
		SendDiscordMessage(m.ChannelID, "Pong!")
	case prefix + "echo":
		SendDiscordMessage(m.ChannelID, cmd[1])
	case prefix + "guild":
		SendDiscordMessage(m.ChannelID, guildID)

	default:
		slog.Info("unknown command", "command", cmd[0])
		return
	}
}
