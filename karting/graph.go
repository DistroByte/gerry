package karting

import (
	"fmt"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"

	"github.com/rs/zerolog/log"
)

func (k *Karting) Graph() {
	if len(k.Drivers) == 0 {
		log.Error().Msg("no drivers available to generate ELO graph")
		return
	}

	// sort the drivers by ELO
	for i := 0; i < len(k.Drivers); i++ {
		for j := i + 1; j < len(k.Drivers); j++ {
			if k.Drivers[i].ELO < k.Drivers[j].ELO {
				k.Drivers[i], k.Drivers[j] = k.Drivers[j], k.Drivers[i]
			}
		}
	}

	p := plot.New()
	p.Title.Text = "ELO over time"
	p.X.Label.Text = "Races"
	p.Y.Label.Text = "ELO"

	// add a grid
	p.Add(plotter.NewGrid())

	// tick every 5 x values
	p.X.Tick.Marker = raceTicker{}
	p.Y.Tick.Marker = eloTicker{}

	// pad the y axis a bit
	p.Y.Min = float64(k.Drivers[len(k.Drivers)-1].ELO - 50)
	p.Y.Max = float64(k.Drivers[0].ELO + 50)

	// pad the x axis a bit
	p.X.Min = 0

	for j, driver := range k.Drivers {
		xys := make(plotter.XYs, len(k.Races)+1)
		labels := make([]string, len(k.Races)+1)

		var firstRaceIndex int = -1

		// loop over every race
		for i, event := range k.Races {

			// loop over every result
			for _, result := range event.Results {

				// if the result is for the driver we're plotting
				if result.Driver.Name == driver.Name {

					// and this is the first time we've seen them
					if firstRaceIndex < 0 {
						// add the initial ELO to the race before their first
						xys[i].X = float64(i)
						xys[i].Y = float64(InitialELO)
						labels[i] = strconv.Itoa(InitialELO)
						// and remember the index
						firstRaceIndex = i
					}

					// add the ELO for the driver for the current race
					xys[i+1].X = float64(i + 1)
					xys[i+1].Y = float64(result.Driver.ELO)
					break
				} else {
					// if we haven't seen the driver yet, just copy the last value
					xys[i+1].X = float64(i)
					xys[i+1].Y = xys[i].Y
				}
			}

			// add an ELO label every 3 races
			if i%3 == 0 {
				labels[i+1] = strconv.FormatFloat(xys[i+1].Y, 'f', 0, 64)
			}
		}

		if firstRaceIndex < 0 {
			log.Debug().Str("driver", driver.Name).Msg("driver never raced")
			// set the last value to the initial ELO
			xys[len(xys)-1].X = float64(len(xys) - 1)
			xys[len(xys)-1].Y = float64(InitialELO)
			labels[len(labels)-1] = strconv.Itoa(InitialELO)
			firstRaceIndex = len(xys) - 1
		}

		// add the last label
		labels[len(labels)-1] = strconv.FormatFloat(xys[len(xys)-1].Y, 'f', 0, 64)

		// create a line for the driver
		line, points, err := plotter.NewLinePoints(xys[firstRaceIndex:])
		if err != nil {
			log.Error().Err(err).Msg("failed to create plot chart line")
		}

		// create labels for the line
		label, err := plotter.NewLabels(plotter.XYLabels{
			XYs:    xys[firstRaceIndex:],
			Labels: labels[firstRaceIndex:],
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to create plot chart labels")
		}

		// style the line and points
		line.Color = colors[j%len(colors)]
		points.Shape = draw.CircleGlyph{}
		points.Color = colors[j%len(colors)]
		line.StepStyle = plotter.NoStep

		// add the line and labels to the plot
		p.Add(line, points, label)

		// add the driver to the legend
		p.Legend.Add(fmt.Sprintf("%s (%d)", driver.Name, driver.ELO), line)
	}

	if err := p.Save(30*vg.Centimeter, 20*vg.Centimeter, "elo.svg"); err != nil {
		log.Error().Err(err).Str("file", "elo.svg").Msg("failed to save file")
		return
	}

	log.Debug().Str("file", "elo.svg").Msg("saved plot to file")

	if err := p.Save(30*vg.Centimeter, 20*vg.Centimeter, "elo.png"); err != nil {
		log.Error().Err(err).Str("file", "elo.png").Msg("failed to save file")
		return
	}

	log.Debug().Str("file", "elo.png").Msg("saved plot to file")
}

type raceTicker struct{}
type eloTicker struct{}

func (t raceTicker) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	for i := min + 1; i < max; i += 3 {
		ticks = append(ticks, plot.Tick{Value: i, Label: strconv.Itoa(int(i))})
	}
	return ticks
}

func (t eloTicker) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	for i := min; i < max; i += 10 {
		ticks = append(ticks, plot.Tick{Value: i, Label: strconv.Itoa(int(i))})
	}
	return ticks
}
