package commands

import (
	"fmt"
	"runtime"
	"time"

	"github.com/DistroByte/gerry/internal/config"
)

// UptimeCommand returns the uptime of the bot
func UptimeCommand() string {
	uptime := time.Since(config.StartTime).Round(time.Second)
	return fmt.Sprintf(
		"```\n"+
			"Uptime: %s\n"+
			"Go version: %s\n"+
			"```",
		uptime, runtime.Version())
}
