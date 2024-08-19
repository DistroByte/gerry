package commands

import (
	"fmt"
	"strings"

	"github.com/DistroByte/gerry/internal/config"
	"github.com/DistroByte/gerry/internal/models"
	"github.com/DistroByte/gerry/karting"
)

var league *karting.Karting
var longestDriverName int

func init() {
	// Initialize the karting instance
	league = karting.NewKarting()
	league.Load()
	league.Graph()

	// Find the longest driver name
	for _, driver := range league.Drivers {
		if len(driver.Name) > longestDriverName {
			longestDriverName = len(driver.Name)
		}
	}
}

func KartingCommand(args []string, message models.Message) string {
	if len(args) == 0 {
		response := "karting command requires arguments"
		return response
	}

	switch args[0] {

	case "race":
		if len(args) < 2 {
			return "please provide a list of drivers"
		}

		var results []*karting.Result
		for i, driver := range args[1:] {
			results = append(results, &karting.Result{
				Position: i + 1,
				Driver:   &karting.Driver{Name: driver},
			})
		}

		raceDiff, err := league.Race(results)
		if err != nil {
			return err.Error()
		}

		response := "# Race results\n"
		response += fmt.Sprintf("```%*s | Rating change\n", longestDriverName, "Driver")
		for _, diff := range raceDiff {
			response += fmt.Sprintf("%*s | %d -> %d (%+d)\n", longestDriverName, diff.Driver.Name, diff.Driver.ELO-diff.Change, diff.Driver.ELO, diff.Change)
		}
		response += "```"

		// update the graph
		league.Graph()

		return response

	case "stats":
		response := fmt.Sprintf("# Karting stats\n```Rating | %-*s | Won | Total | Win %%  | Last 5 avg (all time) | Peak ELO\n", longestDriverName, "Driver")
		response += "------ | " + fmt.Sprintf("%s | --- | ----- | ------ | --------------------- | --------\n", strings.Repeat("-", longestDriverName))

		for i := 0; i < len(league.Drivers); i++ {
			for j := i + 1; j < len(league.Drivers); j++ {
				if league.Drivers[i].ELO < league.Drivers[j].ELO {
					league.Drivers[i], league.Drivers[j] = league.Drivers[j], league.Drivers[i]
				}
			}
		}

		for _, driver := range league.Drivers {
			winRate := float64(driver.Stats.TotalWins) / float64(driver.Stats.TotalRaces) * 100

			var last5 float64
			for _, finish := range driver.Stats.Last5Finish {
				last5 += float64(finish)
			}
			last5 /= 5

			response += fmt.Sprintf("%6d | %-*s | %3d | %5d | %5.2f%% | %21s | %8d\n", driver.ELO, longestDriverName, driver.Name, driver.Stats.TotalWins, driver.Stats.TotalRaces, winRate,
				fmt.Sprintf("%.2f (%.2f)", last5, driver.Stats.AllTimeAverageFinish/float64(driver.Stats.TotalRaces)), driver.Stats.PeakELO)
		}

		response += "```"

		return response

	case "register":
		if len(args) < 2 {
			return "karting register command requires a driver name"
		}

		res, err := league.Register(args[1])
		if err != nil {
			return err.Error()
		}

		// update the longest driver name for formatting
		if len(args[1]) > longestDriverName {
			longestDriverName = len(args[1])
		}

		return res

	case "unregister":
		if len(args) < 2 {
			return "karting unregister command requires a driver name"
		}

		res, err := league.Unregister(args[1])
		if err != nil {
			return err.Error()
		}

		return res

	case "graph":
		league.Graph()
		// return the URL to the graph
		if config.IsEnvironment(config.APP_ENVIRONMENT_LOCAL) {
			return fmt.Sprintf("http://localhost:%d/karting", config.GetHTTPPort())
		} else {
			return fmt.Sprintf("[karting stats](https://%s/karting.png)", config.GetDomain())
		}

	case "reset":
		league.Reset()
		return "karting stats have been reset"

	default:
		break
	}

	response := "invalid karting command"

	return response
}
