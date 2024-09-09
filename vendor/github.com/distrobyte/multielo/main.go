package multielo

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

var (
	ErrPlayerAlreadyExists = errors.New("player already exists")
	ErrPlayerNotFound      = errors.New("player not found")
	ErrMatchNotFound       = errors.New("match not found")
	ErrInvalidMatch        = errors.New("invalid match")
	ErrInvalidELOChange    = errors.New("invalid elo change")
	ErrInvalidPlayer       = errors.New("invalid player")
	ErrInvalidPlayerStats  = errors.New("invalid player stats")
	ErrInvalidLeague       = errors.New("invalid league")
	ErrNoPlayers           = errors.New("no players")
	colors                 = []color.Color{
		color.RGBA{R: 255, A: 255},
		color.RGBA{G: 255, A: 255},
		color.RGBA{B: 255, A: 255},
		color.RGBA{R: 255, G: 255, A: 255},
		color.RGBA{R: 255, B: 255, A: 255},
		color.RGBA{G: 255, B: 255, A: 255},
	}
)

const (
	InitialELO = 1000
)

type Player struct {
	Name      string
	ELO       int
	ELOChange int
	Stats     *PlayerStats
}

type PlayerStats struct {
	MatchesPlayed       int
	MatchesWon          int
	AllTimeAveragePlace float64
	Last5Finish         []int
	PeakELO             int
}

type Match struct {
	Results []*MatchResult
	Date    time.Time
}

type MatchResult struct {
	Position int
	Player   *Player
}

type MatchDiff struct {
	Player *Player
	Diff   int
}

type League struct {
	Players []*Player
	Matches []Match
}

func NewLeague() *League {
	return &League{
		Players: []*Player{},
		Matches: []Match{},
	}
}

func (l *League) AddMatch(date time.Time, results ...*MatchResult) ([]MatchDiff, error) {
	if len(l.Players) == 0 {
		return nil, ErrNoPlayers
	}

	populatedResults := make([]*MatchResult, 0, len(results))
	matchDiff := make([]MatchDiff, 0, len(results))

	// ensure all drivers are registered and flesh out the results with ELOs
	for _, result := range results {
		found := false
		for _, player := range l.Players {
			if player == nil {
				return []MatchDiff{}, fmt.Errorf("nil player found. not recording match")
			}

			if result.Player == nil {
				return []MatchDiff{}, fmt.Errorf("nil player found. not recording match")
			}

			if player.Name == result.Player.Name {
				result.Player.ELO = player.ELO
				populatedResults = append(populatedResults, result)
				found = true
				break
			}
		}

		if !found {
			return []MatchDiff{}, fmt.Errorf("player %q not found. not recording match", result.Player.Name)
		}
	}

	// calculate the ELO changes
	n := len(populatedResults)
	kValue := 32 / (n - 1)

	// loop over every result
	for player, result := range populatedResults {
		curELO := result.Player.ELO
		curPosition := result.Position

		// loop over every other result
		for opponentPlayer, opponentResult := range populatedResults {
			// skip comparing the driver to themselves
			if player == opponentPlayer {
				continue
			}

			opponentELO := opponentResult.Player.ELO
			opponentPosition := opponentResult.Position

			// calculate the expected score
			var S float64

			// if the driver finished higher than the other driver
			if curPosition < opponentPosition {
				S = 1.0
			} else {
				S = 0.0
			}

			// calculate the expected score
			E := 1.0 / (1.0 + math.Pow(10, float64(opponentELO-curELO)/400))

			// update the driver's ELO change
			result.Player.ELOChange += int(math.Round(float64(kValue) * (S - E)))
		}

		// update the driver's ELO
		result.Player.ELO += result.Player.ELOChange
		matchDiff = append(matchDiff, MatchDiff{
			Player: result.Player,
			Diff:   result.Player.ELOChange,
		})
	}

	// update the drivers' ELOs
	for _, result := range populatedResults {
		for _, player := range l.Players {
			if player.Name == result.Player.Name {
				player.ELO = result.Player.ELO

				player.Stats.MatchesPlayed += 1
				if result.Position == 1 {
					player.Stats.MatchesWon++
				}

				player.Stats.AllTimeAveragePlace += float64(result.Position)

				player.Stats.Last5Finish = append(player.Stats.Last5Finish, result.Position)
				if len(player.Stats.Last5Finish) > 5 {
					player.Stats.Last5Finish = player.Stats.Last5Finish[1:]
				}

				if player.ELO > player.Stats.PeakELO {
					player.Stats.PeakELO = player.ELO
				}
				break
			}
		}
	}

	// create the event
	err := l.createEvent(populatedResults)
	if err != nil {
		return []MatchDiff{}, err
	}

	return matchDiff, nil
}

func (l *League) createEvent(results []*MatchResult) error {
	event := Match{
		Results: results,
		Date:    time.Now(),
	}

	l.Matches = append(l.Matches, event)

	return nil
}

func (l *League) AddPlayer(name string) error {
	name = strings.ToLower(name)

	for _, p := range l.Players {
		if strings.ToLower(p.Name) == name {
			return ErrPlayerAlreadyExists
		}
	}

	l.Players = append(l.Players, &Player{
		Name: name,
		ELO:  InitialELO,
		Stats: &PlayerStats{
			Last5Finish:         []int{},
			MatchesPlayed:       0,
			MatchesWon:          0,
			AllTimeAveragePlace: 0,
			PeakELO:             InitialELO,
		},
	})

	return nil
}

func (l *League) GetPlayer(name string) (*Player, error) {
	name = strings.ToLower(name)

	for _, p := range l.Players {
		if strings.ToLower(p.Name) == name {
			return p, nil
		}
	}

	return nil, ErrPlayerNotFound
}

func (l *League) RemovePlayer(name string) error {
	name = strings.ToLower(name)

	for i, p := range l.Players {
		if strings.ToLower(p.Name) == name {
			l.Players = append(l.Players[:i], l.Players[i+1:]...)
			return nil
		}
	}

	return ErrPlayerNotFound
}

func (l *League) ResetPlayers() {
	for _, p := range l.Players {
		p.ELO = InitialELO
		p.Stats = &PlayerStats{}
	}
}

func (l *League) ResetMatches() {
	l.Matches = []Match{}
}

func (l *League) GetPlayers() []*Player {
	return l.Players
}

func (l *League) GetMatches() []Match {
	return l.Matches
}

func (l *League) GetPlayerStats(name string) (*PlayerStats, error) {
	p, err := l.GetPlayer(name)
	if err != nil {
		return nil, err
	}

	return p.Stats, nil
}

func (l *League) GetPlayerELO(name string) (int, error) {
	p, err := l.GetPlayer(name)
	if err != nil {
		return 0, ErrInvalidPlayer
	}

	return p.ELO, nil
}

func (l *League) GenerateGraph() (string, error) {
	if len(l.Players) == 0 {
		return "", ErrNoPlayers
	}

	// sort the drivers by ELO
	for i := 0; i < len(l.Players); i++ {
		for j := i + 1; j < len(l.Players); j++ {
			if l.Players[i].ELO < l.Players[j].ELO {
				l.Players[i], l.Players[j] = l.Players[j], l.Players[i]
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
	p.X.Tick.Marker = RaceTicker{}
	p.Y.Tick.Marker = ELOTicker{}

	// pad the y axis a bit
	p.Y.Min = float64(l.Players[len(l.Players)-1].ELO - 50)
	p.Y.Max = float64(l.Players[0].ELO + 50)

	// pad the x axis a bit
	p.X.Min = 0

	for j, player := range l.Players {
		xys := make(plotter.XYs, len(l.Matches)+1)
		labels := make([]string, len(l.Matches)+1)

		var firstRaceIndex int = -1

		// loop over every race
		for i, event := range l.Matches {

			// loop over every result
			for _, result := range event.Results {

				// if the result is for the driver we're plotting
				if result.Player.Name == player.Name {

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
					xys[i+1].Y = float64(result.Player.ELO)
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
			return "", err
		}

		// create labels for the line
		label, err := plotter.NewLabels(plotter.XYLabels{
			XYs:    xys[firstRaceIndex:],
			Labels: labels[firstRaceIndex:],
		})
		if err != nil {
			return "", err
		}

		// style the line and points
		line.Color = colors[j%len(colors)]
		points.Shape = draw.CircleGlyph{}
		points.Color = colors[j%len(colors)]
		line.StepStyle = plotter.NoStep

		// add the line and labels to the plot
		p.Add(line, points, label)

		// add the driver to the legend
		p.Legend.Add(fmt.Sprintf("%s (%d)", player.Name, player.ELO), line)
	}

	if err := p.Save(30*vg.Centimeter, 20*vg.Centimeter, "elo.svg"); err != nil {
		return "", err
	}

	if err := p.Save(30*vg.Centimeter, 20*vg.Centimeter, "elo.png"); err != nil {
		return "", err
	}

	return "elo.png", nil
}

type RaceTicker struct{}
type ELOTicker struct{}

func (t RaceTicker) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	for i := min + 1; i < max; i += 3 {
		ticks = append(ticks, plot.Tick{Value: i, Label: strconv.Itoa(int(i))})
	}
	return ticks
}

func (t ELOTicker) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	for i := min; i < max; i += 10 {
		ticks = append(ticks, plot.Tick{Value: i, Label: strconv.Itoa(int(i))})
	}
	return ticks
}
