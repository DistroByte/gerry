package cmds

import (
	"runtime"

	"github.com/distrobyte/gerry/internal/config"
	"github.com/spf13/cobra"
)

type versionOptions struct {
	number bool
	commit bool
	all    bool
}

func newVersionCommand() *cobra.Command {
	options := &versionOptions{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print gerry version information",
		Run: func(cmd *cobra.Command, args []string) {
			runVersion(options, cmd.Root())
		},
	}

	flags := cmd.Flags()
	flags.SortFlags = false

	flags.BoolVarP(&options.number, "number", "n", false, "Print the version number")
	flags.BoolVarP(&options.commit, "commit", "c", false, "Print the commit hash")
	flags.BoolVarP(&options.all, "all", "a", false, "Print all version information")

	return cmd
}

func runVersion(options *versionOptions, root *cobra.Command) {
	if options.all {
		root.Printf("%s version: %s\n", root.Name(), root.Version)
		root.Printf("Build date: %s\n", config.BuildDate)
		root.Printf("System version: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		root.Printf("Golang version: %s\n", runtime.Version())
		return
	}

	if options.number {
		root.Println(root.Version)
		return
	}

	if options.commit {
		root.Println(config.GitCommit)
		return
	}

	root.Printf("%s version: %s\n", root.Name(), root.Version)
}
