package cmds

import (
	"github.com/distrobyte/gerry/internal/config"
	"github.com/spf13/cobra"
)

func newConfgenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confgen",
		Short: "Generate a config file",

		Run: func(cmd *cobra.Command, args []string) {
			config.Generate(cmd.Flag("output").Value.String())
		},
	}

	cmd.Flags().StringP("output", "o", "config.yaml", "config file to generate")

	return cmd
}
