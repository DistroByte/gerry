package commands

import (
	"fmt"
	"runtime"

	"github.com/distrobyte/gerry/internal/config"
)

func VersionCommand(args []string) string {
	if len(args) == 0 {
		return fmt.Sprintf("```\nVersion: %s\n```", config.GetVersion())
	} else {
		switch args[0] {
		case "number":
			return fmt.Sprintf("```\nVersion: %s\n```", config.GetVersion())
		case "commit":
			return fmt.Sprintf("```\nCommit: %s\n```", config.GitCommit)
		case "all":
			return fmt.Sprintf("```\nVersion: %s\nBuild date: %s\nSystem version: %s/%s\nGolang version: %s\n```",
				config.GetVersion(), config.BuildDate, runtime.GOOS, runtime.GOARCH, runtime.Version())
		default:
			return "Invalid argument. Use `number`, `commit` or `all`."
		}
	}
}
