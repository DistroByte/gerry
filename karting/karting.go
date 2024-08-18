package karting

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log/slog"
	"os"
	"strings"
	"time"
)

const (
	InitialELO = 1000
)

var colors = []color.Color{
	color.RGBA{R: 255, A: 255},
	color.RGBA{G: 255, A: 255},
	color.RGBA{B: 255, A: 255},
	color.RGBA{R: 255, G: 255, A: 255},
	color.RGBA{R: 255, B: 255, A: 255},
	color.RGBA{G: 255, B: 255, A: 255},
}

type Driver struct {
	Name      string
	ELO       int
	ELOChange int
	Stats     *DriverStats
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

type DriverStats struct {
	TotalRaces           int
	TotalWins            int
	AllTimeAverageFinish float64
	Last5Finish          []int
	PeakELO              int
}

type RaceDiff struct {
	Driver *Driver
	Change int
}

func NewKarting() *Karting {
	return &Karting{
		Drivers: []*Driver{},
		Races:   []Event{},
	}
}

func (k *Karting) Reset() {
	for _, driver := range k.Drivers {
		slog.Info("resetting driver", "driver", driver.Name)
		driver.ELO = InitialELO
		driver.ELOChange = 0
		driver.Stats = &DriverStats{
			TotalRaces:           0,
			TotalWins:            0,
			AllTimeAverageFinish: 0,
			Last5Finish:          []int{},
			PeakELO:              InitialELO,
		}
	}
	k.Races = []Event{}

	save(k)
}

func (k *Karting) Register(name string) (string, error) {
	name = strings.ToLower(name)

	for _, driver := range k.Drivers {
		if strings.ToLower(driver.Name) == name {
			return "", fmt.Errorf("driver %s already exists", name)
		}
	}

	k.Drivers = append(k.Drivers, &Driver{
		Name: name,
		ELO:  InitialELO,
		Stats: &DriverStats{
			TotalRaces:           0,
			TotalWins:            0,
			AllTimeAverageFinish: 0,
			Last5Finish:          []int{},
			PeakELO:              InitialELO,
		},
	})

	err := save(k)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("registered %q", name), nil
}

func (k *Karting) Unregister(name string) (string, error) {
	name = strings.ToLower(name)

	for i, driver := range k.Drivers {
		if strings.ToLower(driver.Name) == name {
			k.Drivers = append(k.Drivers[:i], k.Drivers[i+1:]...)
			err := save(k)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("unregistered %q", name), nil
		}
	}

	return "", fmt.Errorf("driver %s not found", name)
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
