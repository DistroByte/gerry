package commands

import (
	"sort"
	"strconv"

	"github.com/DistroByte/gerry/internal/models"
	"github.com/DistroByte/gerry/karting"
)

var league *karting.Karting

func init() {
	// Initialize the karting instance
	league = karting.NewKarting()
	league.Load()
}

func KartingCommand(args []string, message models.Message) string {
	if len(args) == 0 {
		response := "Karting command requires arguments"
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

		err := league.Race(results)
		if err != nil {
			return err.Error()
		}

		return "race results recorded"

	case "stats":
		stats := league.Stats()
		response := "Karting stats:\n"
		for driver, stat := range stats {
			response += strconv.Itoa(stat.ELO) + " | " + driver + " - " + strconv.Itoa(stat.TotalWins) + " wins out of " + strconv.Itoa(stat.TotalRaces) + " races\n"
		}

		return response

	case "drivers":
		response := "Drivers:\n"
		// sort the drivers by ELO
		sort.Slice(league.Drivers, func(i, j int) bool {
			return league.Drivers[i].ELO > league.Drivers[j].ELO
		})

		for _, driver := range league.Drivers {
			response += strconv.Itoa(driver.ELO) + " | " + driver.Name + "\n"
		}

		return response

	case "register":
		if len(args) < 2 {
			return "Karting register command requires a driver name"
		}

		res, err := league.Register(args[1])
		if err != nil {
			return err.Error()
		}
		return res

	case "reset":
		league.Reset()
		return "Karting stats have been reset"
	}

	response := "Invalid karting command"

	return response
}
