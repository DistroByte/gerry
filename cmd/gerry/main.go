package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/DistroByte/gerry/internal/bot"
	"github.com/DistroByte/gerry/internal/config"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
	providedArgs := os.Args[1:]

	if len(providedArgs) == 0 {
		log.Fatal().Msg("no arguments provided")
	}

	switch providedArgs[0] {

	case "confgen":
		config.Generate(providedArgs[1])

	case "start":
		if len(providedArgs) < 2 {
			log.Fatal().Msg("no config file provided")
		}

		config.Load(providedArgs[1])
		config.StartTime = time.Now()

		if config.GetEnvironment() == config.APP_ENVIRONMENT_LOCAL {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
			log.Debug().Msg("running locally in debug mode")
		}

		bot.Start()
	}
}
