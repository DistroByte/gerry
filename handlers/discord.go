package handlers

import (
	"fmt"
	"log"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type DiscordMessage struct {
	message discord.Message
}

func (m *DiscordMessage) GetContent() string {
	return m.message.Content
}

func (m *DiscordMessage) GetSenderID() string {
	return m.message.Author.ID.String()
}

func (m *DiscordMessage) GetChannelID() string {
	return m.message.ChannelID.String()
}

func (m *DiscordMessage) GetPlatform() string {
	return "discord"
}

func NewDiscordMessage(message discord.Message) *DiscordMessage {
	return &DiscordMessage{message: message}
}

func (m *DiscordMessage) String() string {
	return fmt.Sprintf("DiscordMessage{Content: %s, SenderID: %s, Channel: %s}", m.GetContent(), m.GetSenderID(), m.GetChannelID())
}

func (h *Handlers) OnDiscordReady(event *events.Ready) {
	log.Println("discord bot is ready")
}

func (h *Handlers) OnDiscordMessageCreate(event *events.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}

	msg := NewDiscordMessage(event.Message)
	// println(msg.String())

	err := handleMessage(h.Bot, msg)
	if err != nil {
		log.Println("error handling command", err)
	}
}
