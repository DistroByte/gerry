package bot

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/DistroByte/gerry/internal/config"
	"github.com/DistroByte/gerry/internal/discord"
	"github.com/DistroByte/gerry/internal/mumble"
	"layeh.com/gumble/gumbleutil"
)

func init() {
	slog.Info("Bot initializing...")
}

func Start(configPath string) {
	config.Load(configPath)
	if config.IsEnvironment(config.APP_ENVIRONMENT_TEST) {
		fmt.Println("App environment is test, aborting startup")
		return
	}

	if config.IsDiscordEnabled() {
		discord.InitDiscordSession()
	}

	if config.IsMumbleEnabled() {
		mumble.InitMumbleSession()
	}

	addHandlers()

	if config.IsDiscordEnabled() {
		discord.InitDiscordConnection()
	}

	slog.Info("Bot is running. Press CTRL+C to exit.", "environment", config.GetEnvironment())
	signal.Notify(config.ShutdownChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-config.ShutdownChannel
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
