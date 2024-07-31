package mumble

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"

	"github.com/DistroByte/gerry/internal/config"
	"github.com/DistroByte/gerry/internal/handlers"
	"github.com/DistroByte/gerry/internal/models"
	"layeh.com/gumble/gumble"
)

var MumbleSession *gumble.Client
var MumbleConfig *gumble.Config

func InitMumbleSession() {
	if config.GetMumbleHost() != "" {
		var tlsConfig tls.Config
		if !config.GetMumbleTLS() {
			tlsConfig.InsecureSkipVerify = true
		}

		MumbleConfig = gumble.NewConfig()
		MumbleConfig.Username = config.GetMumbleUsername()

		var err error
		MumbleSession, err = gumble.DialWithDialer(new(net.Dialer),
			fmt.Sprintf("%s:%v", config.GetMumbleHost(), config.GetMumblePort()),
			MumbleConfig,
			&tlsConfig)
		if err != nil {
			slog.Error("failed to create mumble session", "error", err)
		}
	}
}

func MumbleReadyHandler(event *gumble.ConnectEvent) {
	slog.Info("connected to mumble server", "address", event.Client.Conn.RemoteAddr())
}

func MumbleMessageCreateHandler(event *gumble.TextMessageEvent) {
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

	message := &models.Message{
		Content: event.TextMessage.Message,
		Author:  event.TextMessage.Sender.Name,
		Channel: strconv.FormatUint(uint64(event.TextMessage.Sender.Channel.ID), 10),
	}

	response, err := handlers.HandleMessage(message)
	if err == nil {
		channelIDInt, _ := strconv.ParseInt(message.Channel, 10, 32)
		SendMumbleMessage(uint32(channelIDInt), response)
	}
}
