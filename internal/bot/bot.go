package bot

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DistroByte/gerry/http"
	"github.com/DistroByte/gerry/internal/config"
	"github.com/DistroByte/gerry/internal/discord"
	"github.com/DistroByte/gerry/internal/mumble"
	"github.com/rs/zerolog/log"
	"layeh.com/gumble/gumbleutil"
)

func init() {
	log.Info().Msg("bot initializing...")
}

func Start() {
	if config.IsEnvironment(config.APP_ENVIRONMENT_TEST) {
		log.Info().Msg("app environment is test, aborting startup")
		return
	}

	if config.IsDiscordEnabled() {
		discord.InitDiscordSession()
	}

	if config.IsMumbleEnabled() {
		mumble.InitMumbleSession()
	}

	addHandlers()

	if config.IsHTTPEndpointEnabled() {
		go http.ServeHTTP()
	}

	if config.IsDiscordEnabled() {
		discord.InitDiscordConnection()
	}

	log.Info().
		Str("environment", config.GetEnvironment()).
		Str("event", "startup").
		Msg("bot is running. press CTRL+C to exit.")

	signal.Notify(config.ShutdownChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-config.ShutdownChannel
	log.Info().Msg("shutting down...")

	if config.IsDiscordEnabled() {
		discord.DiscordSession.Close()
		log.Info().Msg("discord closed")
	}

	if config.IsMumbleEnabled() {
		mumble.MumbleSession.Disconnect()
		log.Info().Msg("mumble disconnected")
	}

	log.Info().Msg("goodbye")
	os.Exit(0)
}

func addHandlers() {
	// register hanlders as callbacks for various events

	// discord
	if config.IsDiscordEnabled() {
		discord.DiscordSession.AddHandler(discord.DiscordReadyHandler)
		discord.DiscordSession.AddHandler(discord.DiscordMessageCreateHandler)
	}
	// mumble
	if config.IsMumbleEnabled() {
		mumble.MumbleSession.Config.Attach(gumbleutil.Listener{
			Connect:     mumble.MumbleReadyHandler,
			TextMessage: mumble.MumbleMessageCreateHandler,
		})
	}
}
