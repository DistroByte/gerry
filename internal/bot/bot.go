package bot

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/DistroByte/gerry/internal/config"
	"github.com/DistroByte/gerry/internal/discord"
)

func Start() {
	config.Load()
	if config.IsEnvironment(config.APP_ENVIRONMENT_TEST) {
		fmt.Println("App environment is test, aborting startup")
		return
	}

	discord.InitSession()
	addHandlers()
	discord.InitConnection()

	slog.Info("Bot is running. Press CTRL+C to exit.", "environment", config.GetEnvironment())
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func addHandlers() {
	// register hanlders as callbacks for various events

	// discord
	discord.Session.AddHandler(discord.DiscordReadyHandler)
	discord.Session.AddHandler(discord.DiscordMessageCreateHandler)
}
