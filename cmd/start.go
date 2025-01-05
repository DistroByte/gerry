package cmd

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/distrobyte/gerry/internal/bot"
	"github.com/distrobyte/gerry/internal/config"
)

func NewStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the bot",
		Long:  "Start the bot with the provided config file",

		RunE: func(cmd *cobra.Command, args []string) error {
			err := config.Load(cmd.Flag("config").Value.String())
			if err != nil {
				return err
			}

			config.StartTime = time.Now()

			zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
			zerolog.SetGlobalLevel(zerolog.InfoLevel)

			if config.GetEnvironment() == config.APP_ENVIRONMENT_LOCAL {
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
				log.Debug().Msg("running locally in debug mode")
			}

			return bot.Start()
		},
	}

	cmd.Flags().StringP("config", "c", "config.yaml", "config file to use")

	return cmd
}
