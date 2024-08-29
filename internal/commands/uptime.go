package commands

import (
	"fmt"
	"time"

	"github.com/distrobyte/gerry/internal/config"
)

// UptimeCommand returns the uptime of the bot
func UptimeCommand() string {
	// uptime should be in the form months weeks days hours minutes seconds, omitting values that are 0
	year, month, day, hour, min, sec := diff(config.StartTime, time.Now())

	uptime := ""
	if year > 0 {
		uptime += fmt.Sprintf("%d years ", year)
	}
	if month > 0 {
		uptime += fmt.Sprintf("%d months ", month)
	}
	if day > 0 {
		uptime += fmt.Sprintf("%d days ", day)
	}
	if hour > 0 {
		uptime += fmt.Sprintf("%d hours ", hour)
	}
	if min > 0 {
		uptime += fmt.Sprintf("%d minutes ", min)
	}
	if sec > 0 {
		uptime += fmt.Sprintf("%d seconds", sec)
	}

	return fmt.Sprintf(
		"```\n"+
			"Uptime: %s\n"+
			"```",
		uptime)
}

func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}
