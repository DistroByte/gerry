package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/distrobyte/gerry/http"
	"github.com/distrobyte/gerry/internal/config"
	"github.com/distrobyte/gerry/internal/discord"
	"github.com/distrobyte/gerry/internal/handlers"
	"github.com/distrobyte/gerry/internal/mumble"
	"github.com/rs/zerolog/log"
	"layeh.com/gumble/gumbleutil"
)

func Start() error {
	log.Info().Msg("bot initializing...")

	if config.IsEnvironment(config.APP_ENVIRONMENT_TEST) {
		log.Info().Msg("app environment is test, aborting startup")
		return fmt.Errorf("app environment is test")
	}

	if config.IsDiscordEnabled() {
		discord.InitSession()
	}

	if config.IsMumbleEnabled() {
		mumble.InitSession()
	}

	addHandlers()

	handlers.InitCommands()

	if config.IsHTTPEndpointEnabled() {
		go http.ServeHTTP()
	}

	if config.IsDiscordEnabled() {
		discord.InitConnection()
	}

	log.Info().
		Str("environment", config.GetEnvironment()).
		Str("event", "startup").
		Msg("bot is running. press CTRL+C to exit.")

	signal.Notify(config.ShutdownChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-config.ShutdownChannel
	log.Info().Msg("shutting down...")

	if config.IsDiscordEnabled() {
		err := discord.DiscordSession.Close()
		if err != nil {
			log.Error().Err(err).Msg("discord close error")
		}

		log.Info().Msg("discord closed")
	}

	if config.IsMumbleEnabled() {
		err := mumble.MumbleSession.Disconnect()
		if err != nil {
			log.Error().Err(err).Msg("mumble disconnect error")
		}

		log.Info().Msg("mumble disconnected")
	}

	log.Info().Msg("goodbye")
	return nil
}

func addHandlers() {
	// register hanlders as callbacks for various events

	// discord
	if config.IsDiscordEnabled() {
		discord.DiscordSession.AddHandler(discord.ReadyHandler)
		discord.DiscordSession.AddHandler(discord.MessageCreateHandler)
		discord.DiscordSession.AddHandler(discord.MessageReactHandler)
	}
	// mumble
	if config.IsMumbleEnabled() {
		mumble.MumbleSession.Config.Attach(gumbleutil.Listener{
			Connect:     mumble.ReadyHandler,
			TextMessage: mumble.MessageCreateHandler,
		})
	}
}
