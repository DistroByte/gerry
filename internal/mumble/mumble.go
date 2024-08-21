package mumble

import (
	"strings"

	"github.com/rs/zerolog/log"
)

func SendMumbleMessage(channelID uint32, message string) {

	channel := MumbleSession.Channels[channelID]
	if channel == nil {
		log.Warn().Str("platform", "mumble").Msg("channel not found")
		return
	}

	if strings.Contains(message, "\n") {
		message = strings.ReplaceAll(message, "\n", "<br>")
	}

	channel.Send(message, false)
}
