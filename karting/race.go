package karting

import (
	"fmt"
	"math"
	"time"
)

func (k *Karting) Race(results []*Result) ([]RaceDiff, error) {
	if k == nil || k.Drivers == nil {
		return nil, fmt.Errorf("no drivers registered")
	}

	populatedResults := make([]*Result, 0, len(results))
	raceDiff := make([]RaceDiff, 0, len(results))

	// ensure all drivers are registered and flesh out the results with ELOs
	for _, result := range results {
		found := false
		for _, driver := range k.Drivers {
			if driver.Name == result.Driver.Name {
				result.Driver.ELO = driver.ELO
				populatedResults = append(populatedResults, result)
				found = true
				break
			}
		}

		if !found {
			return []RaceDiff{}, fmt.Errorf("driver %q not found. not recording race", result.Driver.Name)
		}
	}

	// calculate the ELO changes
	n := len(populatedResults)
	kValue := 32 / (n - 1)

	// loop over every result
	for driver, result := range populatedResults {
		curELO := result.Driver.ELO
		curPosition := result.Position

		// loop over every other result
		for opponentDriver, opponentResult := range populatedResults {
			// skip comparing the driver to themselves
			if driver == opponentDriver {
				continue
			}

			opponentELO := opponentResult.Driver.ELO
			opponentPosition := opponentResult.Position

			// calculate the expected score
			var S float64

			// if the driver finished higher than the other driver
			if curPosition < opponentPosition {
				S = 1.0
			} else if curPosition == opponentPosition {
				S = 0.5
			} else {
				S = 0.0
			}

			// calculate the expected score
			E := 1.0 / (1.0 + math.Pow(10, float64(opponentELO-curELO)/400))

			// update the driver's ELO change
			result.Driver.ELOChange += int(math.Round(float64(kValue) * (S - E)))
		}

		// update the driver's ELO
		result.Driver.ELO += result.Driver.ELOChange
		raceDiff = append(raceDiff, RaceDiff{
			Driver: result.Driver,
			Change: result.Driver.ELOChange,
		})
	}

	// update the drivers' ELOs
	for _, result := range populatedResults {
		for _, driver := range k.Drivers {
			if driver.Name == result.Driver.Name {
				driver.ELO = result.Driver.ELO
				driver.Stats.TotalRaces++
				if result.Position == 1 {
					driver.Stats.TotalWins++
				}

				driver.Stats.AllTimeAverageFinish += float64(result.Position)

				driver.Stats.Last5Finish = append(driver.Stats.Last5Finish, result.Position)
				if len(driver.Stats.Last5Finish) > 5 {
					driver.Stats.Last5Finish = driver.Stats.Last5Finish[1:]
				}

				if driver.ELO > driver.Stats.PeakELO {
					driver.Stats.PeakELO = driver.ELO
				}
				break
			}
		}
	}

	// create the event
	event := Event{
		Results: populatedResults,
		Date:    time.Now(),
	}
	k.Races = append(k.Races, event)

	return raceDiff, nil
}
