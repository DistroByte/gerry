package karting

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
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
		log.Debug().Str("driver", driver.Name).Msg("resetting driver")
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

	return fmt.Sprintf("registered %q", name), nil
}

func (k *Karting) Unregister(name string) (string, error) {
	name = strings.ToLower(name)

	for i, driver := range k.Drivers {
		if strings.ToLower(driver.Name) == name {
			k.Drivers = append(k.Drivers[:i], k.Drivers[i+1:]...)
			return fmt.Sprintf("unregistered %q", name), nil
		}
	}

	return "", fmt.Errorf("driver %s not found", name)
}
