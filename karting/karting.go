package karting

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

const InitialELO = 1000

type Driver struct {
	Name string
	ELO  int
}

type Result struct {
	Position int
	Driver   *Driver
}

type Event struct {
	Results []*Result
	Date    time.Time
}

type Karting struct {
	Drivers []*Driver
	Races   []Event
}

type driverStats struct {
	TotalRaces int
	TotalWins  int
	ELO        int
}

func NewKarting() *Karting {
	return &Karting{
		Drivers: []*Driver{},
		Races:   []Event{},
	}
}

func (k *Karting) Race(results []*Result) error {
	// the result we get passed has no ELO, so we need to look up the driver in our list
	populatedResults := make([]*Result, 0, len(results))

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
			return fmt.Errorf("driver %q not found", result.Driver.Name)
		}

		// update the driver's ELO
		if result.Position == 1 {
			result.Driver.ELO += 32
		}
	}

	event := Event{
		Results: populatedResults,
		Date:    time.Now(),
	}

	// update the global driver elo
	for _, driver := range k.Drivers {
		for _, result := range event.Results {
			if driver.Name == result.Driver.Name {
				driver.ELO = result.Driver.ELO
			}
		}
	}

	k.Races = append(k.Races, event)

	save(k)

	return nil
}

func (k *Karting) Stats() map[string]*driverStats {
	return calculateStats(k)
}

func (k *Karting) Reset() {
	for _, driver := range k.Drivers {
		slog.Info("resetting driver", "driver", driver.Name)
		driver.ELO = InitialELO
	}

	save(k)
}

func (k *Karting) Register(name string) (string, error) {
	name = strings.ToLower(name)

	for _, driver := range k.Drivers {
		if strings.ToLower(driver.Name) == name {
			return "", fmt.Errorf("driver %s already exists", name)
		}
	}

	k.Drivers = append(k.Drivers, &Driver{Name: name, ELO: InitialELO})

	err := save(k)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("registered %q", name), nil
}

func (k *Karting) Load() {
	slog.Info("loading karting data")

	data, err := os.ReadFile("karting.json")
	if err != nil {
		slog.Error("failed to read karting data", "error", err)
	}

	err = json.Unmarshal(data, &k)
	if err != nil {
		slog.Error("failed to unmarshal karting data", "error", err)
	}
}

func save(karting *Karting) error {
	slog.Info("writing karting data")

	out, err := json.Marshal(karting)
	if err != nil {
		return err
	}

	err = os.WriteFile("karting.json", out, 0644)
	if err != nil {
		return err
	}

	return nil
}

func calculateStats(karting *Karting) map[string]*driverStats {
	stats := make(map[string]*driverStats)

	for _, event := range karting.Races {
		for _, result := range event.Results {
			stat, ok := stats[result.Driver.Name]
			if !ok {
				stat = &driverStats{}
				stats[result.Driver.Name] = stat
			}

			stat.TotalRaces++
			if result.Position == 1 {
				stat.TotalWins++
			}

			stat.ELO = result.Driver.ELO
		}
	}

	return stats
}
