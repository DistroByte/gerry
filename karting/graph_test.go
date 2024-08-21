package karting

import (
	"testing"

	"gonum.org/v1/plot"
)

func testTicker(t *testing.T, ticker plot.Ticker, start, end float64, expected int) {
	ticks := ticker.Ticks(start, end)
	if len(ticks) != expected {
		t.Errorf("expected %d ticks, got %v", expected, len(ticks))
	}
}

func TestTicks(t *testing.T) {
	tests := []struct {
		name     string
		ticker   plot.Ticker
		start    float64
		end      float64
		expected int
	}{
		{"raceTicker 0", raceTicker{}, 0, 0, 0},
		{"eloTicker 0", eloTicker{}, 0, 0, 0},
		{"raceTicker 100", raceTicker{}, 0, 100, 33},
		{"eloTicker 100", eloTicker{}, 0, 100, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTicker(t, tt.ticker, tt.start, tt.end, tt.expected)
		})
	}
}
