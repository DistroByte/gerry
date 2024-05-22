package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"layeh.com/gumble/gumble"
)

type MumbleMessage struct {
	message *gumble.TextMessage
}

func (m *MumbleMessage) GetContent() string {
	return m.message.Message
}

func (m *MumbleMessage) GetSenderID() string {
	return m.message.Sender.Name
}

func (m *MumbleMessage) GetChannelID() string {
	return strconv.FormatUint(uint64(m.message.Sender.Channel.ID), 10)
}

func (m *MumbleMessage) GetPlatform() string {
	return "mumble"
}

func NewMumbleMessage(message *gumble.TextMessage) *MumbleMessage {
	return &MumbleMessage{message: message}
}

func (m *MumbleMessage) String() string {
	return fmt.Sprintf("MumbleMessage{Content: %s, SenderID: %s, Channel: %s}", m.GetContent(), m.GetSenderID(), m.GetChannelID())
}

func (h *Handlers) OnMumbleReady(event *gumble.ConnectEvent) {
	log.Println("mumble bot is ready")
}

func (h *Handlers) OnMumbleTextMessage(event *gumble.TextMessageEvent) {
	if event.Sender == nil || event.Sender.Name == "" {
		return
	}

	specialChars := map[string]string{
		"<": "&lt;",
		">": "&gt;",
	}

	// replace special characters
	for k, v := range specialChars {
		event.TextMessage.Message = strings.ReplaceAll(event.TextMessage.Message, v, k)
	}

	msg := NewMumbleMessage(&event.TextMessage)
	// println(msg.String())

	err := handleMessage(h.Bot, msg)
	if err != nil {
		log.Println("error handling command", err)
	}
}

func (h *Handlers) OnMumbleUserChange(event *gumble.UserChangeEvent) {
	if event.Type != gumble.UserChangeConnected {
		return
	}

	log.Printf("user %s connected\n", event.User.Name)
}

func (h *Handlers) OnMumbleChannelChange(event *gumble.ChannelChangeEvent) {
	if event.Channel == nil {
		return
	}

	log.Printf("channel changed to %s\n", event.Channel.Name)
}

func (h *Handlers) OnMumbleDisconnect(event *gumble.DisconnectEvent) {
	log.Println("disconnected from mumble server")
}
