package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/distrobyte/gerry/internal/config"
	"github.com/distrobyte/gerry/internal/models"
	"github.com/distrobyte/multielo"
	"github.com/rs/zerolog/log"
)

var league *multielo.League
var longestPlayerName int

func InitKarting() {
	// Initialize the karting instance
	league = multielo.NewLeague()

	err := load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load karting data")
	}

	// Find the longest driver name
	for _, driver := range league.Players {
		if len(driver.Name) > longestPlayerName {
			longestPlayerName = len(driver.Name)
		}
	}
}

func KartingCommand(args []string, message models.Message) string {
	if len(args) == 0 {
		response := "karting command requires arguments"
		return response
	}

	switch args[0] {
	case "register":
		if len(args) < 2 {
			return "karting register command requires a driver name"
		}

		err := league.AddPlayer(args[1])
		if err != nil {
			return err.Error()
		}

		// update the longest driver name for formatting
		if len(args[1]) > longestPlayerName {
			longestPlayerName = len(args[1])
		}

		err = save()
		if err != nil {
			return err.Error()
		}
		return "driver registered"

	case "unregister":
		if len(args) < 2 {
			return "karting unregister command requires a driver name"
		}

		err := league.RemovePlayer(args[1])
		if err != nil {
			return err.Error()
		}

		err = save()
		if err != nil {
			return err.Error()
		}
		return "driver unregistered"

	case "graph":
		_, err := league.GenerateGraph()
		if err != nil {
			return err.Error()
		}

		// return the URL to the graph
		if config.IsEnvironment(config.APP_ENVIRONMENT_LOCAL) {
			return fmt.Sprintf("http://localhost:%d/karting", config.GetHTTPPort())
		} else {
			return fmt.Sprintf("https://%s/karting.png", config.GetDomain())
		}

	case "stats":
		return KartingStatsCommand(args, message)

	case "race":
		return KartingRaceCommand(args, message)

	case "reset":
		league.ResetPlayers()
		league.ResetMatches()

		err := save()
		if err != nil {
			return err.Error()
		}

		load()

		return "karting stats have been reset"

	default:
		break
	}

	response := "invalid karting command"

	return response
}

func KartingStatsCommand(args []string, message models.Message) string {
	response := fmt.Sprintf("# Karting stats\n```Rating | %-*s | Won | Total | Win %%  | Last 5 avg (all time) | Peak ELO\n", longestPlayerName, "Driver")
	response += "------ | " + fmt.Sprintf("%s | --- | ----- | ------ | --------------------- | --------\n", strings.Repeat("-", longestPlayerName))

	for i := 0; i < len(league.Players); i++ {
		for j := i + 1; j < len(league.Players); j++ {
			if league.Players[i].ELO < league.Players[j].ELO {
				league.Players[i], league.Players[j] = league.Players[j], league.Players[i]
			}
		}
	}

	for _, driver := range league.Players {
		winRate := float64(driver.Stats.MatchesWon) / float64(driver.Stats.MatchesPlayed) * 100

		var last5 float64
		for _, finish := range driver.Stats.Last5Finish {
			last5 += float64(finish)
		}
		last5 /= 5

		response += fmt.Sprintf("%6d | %-*s | %3d | %5d | %5.2f%% | %21s | %8d\n", driver.ELO, longestPlayerName, driver.Name, driver.Stats.MatchesWon, driver.Stats.MatchesPlayed, winRate,
			fmt.Sprintf("%.2f (%.2f)", last5, driver.Stats.AllTimeAveragePlace/float64(driver.Stats.MatchesPlayed)), driver.Stats.PeakELO)
	}

	response += "```"

	return response
}

func KartingRaceCommand(args []string, message models.Message) string {
	if len(args) < 2 {
		return "please provide a list of drivers"
	}

	var results []*multielo.MatchResult
	for i, driver := range args[1:] {
		results = append(results, &multielo.MatchResult{
			Position: i + 1,
			Player:   &multielo.Player{Name: driver},
		})
	}

	raceDiff, err := league.AddMatch(time.Now(), results...)
	if err != nil {
		return err.Error()
	}

	response := "# Race results\n"
	response += fmt.Sprintf("```%*s | Rating change\n", longestPlayerName, "Driver")
	for _, diff := range raceDiff {
		response += fmt.Sprintf("%*s | %d -> %d (%+d)\n", longestPlayerName, diff.Player.Name, diff.Player.ELO-diff.Player.ELOChange, diff.Player.ELO, diff.Diff)
	}
	response += "```"

	// update the graph
	_, err = league.GenerateGraph()
	if err != nil {
		return err.Error()
	}

	err = save()
	if err != nil {
		return err.Error()
	}

	return response
}

func save() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	log.Info().
		Str("file", dir+"/karting.json").
		Msg("writing karting data to")

	out, err := json.Marshal(league)
	if err != nil {
		return err
	}

	err = os.WriteFile(dir+"/karting.json", out, 0644)
	if err != nil {
		return err
	}

	return nil
}

func load() error {
	log.Info().Msg("loading karting data")
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(dir + "/karting.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read karting data")
		return err
	}

	err = json.Unmarshal(data, &league)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal karting data")
		return err
	}

	return nil
}
