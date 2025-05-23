package cmd

import (
	"os"

	"github.com/distrobyte/gerry/internal/config"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gerry",
		Short: "Gerry is a platform-agnostic bot",
		Long:  "Gerry is a platform-agnostic bot that has a variety of features",

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			root := cmd.Root()
			root.Version = config.GetVersion()
		},

		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				os.Exit(1)
			}
		},

		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}

	addCmd := func(child *cobra.Command) {
		cmd.AddCommand(child)
	}

	addCmd(NewVersionCommand())
	addCmd(NewStartCommand())
	addCmd(NewConfgenCommand())

	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return cmd
}

func Execute() {
	if err := RootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
