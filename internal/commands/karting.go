package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/distrobyte/gerry/internal/config"
	"github.com/distrobyte/gerry/internal/models"
	"github.com/distrobyte/multielo"
	"github.com/rs/zerolog/log"
)

var league *multielo.League
var leagues *multielo.MultiLeagueService
var longestPlayerName int

// multielo -> zerolog adapter to surface logs from vendored module
type multieloZerologAdapter struct{}

func (multieloZerologAdapter) Info(msg string, fields map[string]interface{}) {
	log.Info().Fields(fields).Msg(msg)
}

func (multieloZerologAdapter) Error(msg string, err error, fields map[string]interface{}) {
	log.Error().Err(err).Fields(fields).Msg(msg)
}

func (multieloZerologAdapter) Debug(msg string, fields map[string]interface{}) {
	log.Debug().Fields(fields).Msg(msg)
}

// persistedState is a simple snapshot format for players and matches.
// It avoids marshaling internal League fields (which are private) and
// reconstructs the League by re-adding players and replaying matches.
type persistedResult struct {
	Position int    `json:"position"`
	Player   string `json:"player"`
}

type persistedMatch struct {
	Results []persistedResult `json:"results"`
	Date    time.Time         `json:"date"`
}

type persistedState struct {
	Players []string         `json:"players"`
	Matches []persistedMatch `json:"matches"`
}

type persistedConfig struct {
	Config multielo.LeagueConfig `json:"config"`
}

func InitKarting() {
	// Initialize the karting instance
	cfg := multielo.DefaultConfig()
	cfg.OutputDirectory = "assets"
	cfg.DecayEnabled = true
	cfg.DecayInitialPercent = 1.0
	cfg.DecayPerMiss = -0.3

	// Persist config before league is constructed
	err := saveConfig(cfg)
	if err != nil {
		log.Error().Err(err).Msg("failed to save config data")
	}

	league = multielo.NewLeagueWithDependencies(cfg, multielo.LeagueDependencies{Logger: multieloZerologAdapter{}})

	err = load()
	if err != nil {
		log.Error().Err(err).Msg("failed to load karting data")
	}

	// Find the longest driver name
	for _, driver := range league.GetPlayers() {
		if len(driver.Name()) > longestPlayerName {
			longestPlayerName = len(driver.Name())
		}
	}

	_, err = league.GenerateGraph()
	if err != nil {
		log.Error().Err(err).Msg("failed to generate karting graph on startup")
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

		err = load()
		if err != nil {
			return err.Error()
		}

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

	// Get all players and sort by ELO descending
	players := league.GetPlayers()
	sort.Slice(players, func(i, j int) bool {
		return players[i].ELO() > players[j].ELO()
	})

	for _, driver := range players {
		matchesPlayed := driver.MatchesPlayed()
		matchesWon := driver.MatchesWon()
		var winRate float64
		if matchesPlayed > 0 {
			winRate = float64(matchesWon) / float64(matchesPlayed) * 100
		}

		var last5 float64
		last5Finishes := driver.Last5Finish()
		for _, finish := range last5Finishes {
			last5 += float64(finish)
		}
		if len(last5Finishes) > 0 {
			last5 /= float64(len(last5Finishes))
		}

		response += fmt.Sprintf("%6d | %-*s | %3d | %5d | %5.2f%% | %21s | %8d\n", driver.ELO(), longestPlayerName, driver.Name(), matchesWon, matchesPlayed, winRate,
			fmt.Sprintf("%.2f (%.2f)", last5, driver.AllTimeAvgPlace()), driver.PeakELO())
	}

	response += "```"

	return response
}

func KartingRaceCommand(args []string, message models.Message) string {
	if len(args) < 2 {
		return "please provide a list of drivers"
	}

	// Track before state for display
	beforeELOs := make(map[string]int)

	// Build match results with league players
	var results []*multielo.MatchResult
	for i, driverName := range args[1:] {
		// Get or create player
		player, err := league.GetPlayer(driverName)
		if err != nil {
			// Player doesn't exist, add them to the league first
			_ = league.AddPlayer(driverName)
			player, _ = league.GetPlayer(driverName)
		}
		// Capture pre-race ELO (handles newly added players correctly)
		beforeELOs[driverName] = player.ELO()
		results = append(results, &multielo.MatchResult{
			Position: i + 1,
			Player:   player,
		})
	}

	err := league.AddMatch(results, time.Now())
	if err != nil {
		return err.Error()
	}

	response := "# Race results\n"
	response += fmt.Sprintf("```%*s | Change | Cause\n", longestPlayerName, "Driver")

	// Use last changes from multielo to annotate cause (position/decay)
	last := multielo.GetLastChanges(league)
	// Index by name for quick lookup
	changeByName := make(map[string]multielo.LastChange)
	for _, c := range last {
		changeByName[c.Name] = c
	}

	// Participants first (in the order provided)
	for _, driverName := range args[1:] {
		player, _ := league.GetPlayer(driverName)
		if player == nil {
			continue
		}
		before := beforeELOs[driverName]
		after := player.ELO()
		diff := after - before
		cause := "position"
		if c, ok := changeByName[driverName]; ok {
			cause = c.Cause
			if c.Cause == "position" && c.Position > 0 {
				cause = fmt.Sprintf("position (%d)", c.Position)
			}
		}
		response += fmt.Sprintf("%*s | %+d | %s\n", longestPlayerName, driverName, diff, cause)
	}

	// List decays for non-participants
	for name, c := range changeByName {
		if !c.Attended && c.Cause == "decay" {
			response += fmt.Sprintf("%*s | %+d | decay\n", longestPlayerName, name, c.Diff)
		}
	}

	response += "```"

	// Sync player histories to ensure all players have complete history for graph rendering
	multielo.SyncPlayerHistories(league)

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
	// Ensure assets directory exists
	if err := os.MkdirAll(dir+"/assets", 0755); err != nil {
		log.Error().Err(err).Msg("failed to create assets directory")
		return err
	}
	log.Info().
		Str("file", dir+"/assets/karting.json").
		Msg("writing karting data to")

	// Build snapshot state from current league
	players := league.GetPlayers()
	pnames := make([]string, 0, len(players))
	for _, p := range players {
		pnames = append(pnames, p.Name())
	}

	matches := league.GetMatches()
	pmatches := make([]persistedMatch, 0, len(matches))
	for _, m := range matches {
		pr := make([]persistedResult, 0, len(m.Results))
		for _, r := range m.Results {
			pr = append(pr, persistedResult{Position: r.Position, Player: r.Player.Name()})
		}
		pmatches = append(pmatches, persistedMatch{Results: pr, Date: m.Date})
	}

	state := persistedState{Players: pnames, Matches: pmatches}

	out, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal karting state")
		return err
	}

	err = os.WriteFile(dir+"/assets/karting.json", out, 0644)
	if err != nil {
		log.Error().Err(err).Msg("failed to write karting state")
		return err
	}

	return nil
}

func saveConfig(cfg multielo.LeagueConfig) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	// Ensure assets directory exists
	if err := os.MkdirAll(dir+"/assets", 0755); err != nil {
		log.Error().Err(err).Msg("failed to create assets directory")
		return err
	}
	log.Info().
		Str("file", dir+"/assets/karting_config.json").
		Msg("writing karting config to")

	// Write the provided config (do not depend on league being initialized yet)
	configState := persistedConfig{Config: cfg}
	out, err := json.MarshalIndent(configState, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal karting config")
		return err
	}

	err = os.WriteFile(dir+"/assets/karting_config.json", out, 0644)
	if err != nil {
		log.Error().Err(err).Msg("failed to write karting config")
		return err
	}

	return nil
}

func load() error {
	log.Info().Str("file", "assets/karting.json").Msg("loading karting data")
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(dir + "/assets/karting.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read karting state")
		return err
	}

	var state persistedState
	if err := json.Unmarshal(data, &state); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal karting state")
		return err
	}

	// Load config with sane fallback to defaults
	cfg := multielo.DefaultConfig()
	var configState persistedConfig
	if configData, err := os.ReadFile(dir + "/assets/karting_config.json"); err != nil {
		log.Error().Err(err).Msg("failed to read karting config, using defaults")
	} else {
		log.Info().Str("file", "assets/karting_config.json").Msg("loading karting config")
		if err := json.Unmarshal(configData, &configState); err != nil {
			log.Error().Err(err).Msg("failed to unmarshal karting config, using defaults")
		} else {
			cfg = configState.Config
		}
	}

	// Reconstruct league from snapshot
	league = multielo.NewLeagueWithDependencies(cfg, multielo.LeagueDependencies{Logger: multieloZerologAdapter{}})

	// Don't pre-add all players; let them be auto-created when first appearing in matches.
	// This ensures players only get history entries starting from when they first participate.

	// Replay matches (ensures ELO history is rebuilt)
	for _, m := range state.Matches {
		var results []*multielo.MatchResult
		for _, r := range m.Results {
			player, err := league.GetPlayer(r.Player)
			if err != nil {
				// If player wasn't in the players list, create on the fly
				_ = league.AddPlayer(r.Player)
				player, _ = league.GetPlayer(r.Player)
			}
			results = append(results, &multielo.MatchResult{Position: r.Position, Player: player})
		}
		if err := league.AddMatch(results, m.Date); err != nil {
			log.Error().Err(err).Msg("failed to replay match from state")
			return err
		}
	}

	// Sync player histories to backfill entries for players who joined late
	multielo.SyncPlayerHistories(league)

	return nil
}
