package mumble

import (
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/distrobyte/gerry/internal/config"
	"github.com/distrobyte/gerry/internal/handlers"
	"github.com/distrobyte/gerry/internal/models"
	"github.com/rs/zerolog/log"
	"layeh.com/gumble/gumble"
)

var MumbleSession *gumble.Client
var MumbleConfig *gumble.Config

func InitSession() {
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
			log.Fatal().Err(err).Msg("failed to create mumble session")
		}

		log.Info().
			Str("host", config.GetMumbleHost()).
			Int("port", config.GetMumblePort()).
			Msg("mumble session created")
	}
}

func ReadyHandler(event *gumble.ConnectEvent) {
	log.Info().
		Str("address", event.Client.Conn.RemoteAddr().String()).
		Msg("connected to mumble server")
}

func MessageCreateHandler(event *gumble.TextMessageEvent) {
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
		Content:  event.TextMessage.Message,
		Author:   event.TextMessage.Sender.Name,
		Channel:  strconv.FormatUint(uint64(event.TextMessage.Sender.Channel.ID), 10),
		ID:       strconv.FormatInt(time.Now().UnixNano(), 10),
		Platform: "mumble",
	}

	response, err := handlers.HandleMessage(message)
	if err == nil && response != "" {
		channelIDInt, _ := strconv.ParseInt(message.Channel, 10, 32)
		SendMessage(uint32(channelIDInt), response)
	}
}
