package karting

import (
	"testing"
)

func TestKarting_NewKarting(t *testing.T) {
	k := NewKarting()

	if len(k.Drivers) != 0 {
		t.Errorf("expected drivers to be empty, got %v", k.Drivers)
	}

	if len(k.Races) != 0 {
		t.Errorf("expected races to be empty, got %v", k.Races)
	}

	if k.Drivers == nil {
		t.Errorf("expected drivers to be not nil, got nil")
	}
}

func TestKarting_Register(t *testing.T) {
	k := NewKarting()

	_, err := k.Register("Driver1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = k.Register("Driver1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	expectedError := "driver driver1 already exists"
	if err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestKarting_Unregister(t *testing.T) {
	k := NewKarting()

	_, err := k.Register("Driver1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = k.Unregister("Driver1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = k.Unregister("Driver1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	expectedError := "driver driver1 not found"
	if err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestKarting_Reset(t *testing.T) {
	k := NewKarting()

	_, err := k.Register("Driver1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = k.Register("Driver2")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	k.Drivers[0].ELO = 1000
	k.Drivers[0].ELOChange = 100
	k.Drivers[0].Stats.TotalRaces = 10
	k.Drivers[0].Stats.TotalWins = 5
	k.Drivers[0].Stats.AllTimeAverageFinish = 2
	k.Drivers[0].Stats.Last5Finish = []int{1, 2, 3, 4, 5}
	k.Drivers[0].Stats.PeakELO = 1100

	k.Drivers[1].ELO = 1000
	k.Drivers[1].ELOChange = 100
	k.Drivers[1].Stats.TotalRaces = 10
	k.Drivers[1].Stats.TotalWins = 5
	k.Drivers[1].Stats.AllTimeAverageFinish = 2
	k.Drivers[1].Stats.Last5Finish = []int{1, 2, 3, 4, 5}
	k.Drivers[1].Stats.PeakELO = 1100

	k.Reset()

	if k.Drivers[0].ELO != InitialELO {
		t.Errorf("expected ELO to be %d, got %d", InitialELO, k.Drivers[0].ELO)
	}

	if k.Drivers[0].ELOChange != 0 {
		t.Errorf("expected ELOChange to be 0, got %d", k.Drivers[0].ELOChange)
	}

	if k.Drivers[0].Stats.TotalRaces != 0 {
		t.Errorf("expected TotalRaces to be 0, got %d", k.Drivers[0].Stats.TotalRaces)
	}

	if k.Drivers[0].Stats.TotalWins != 0 {
		t.Errorf("expected TotalWins to be 0, got %d", k.Drivers[0].Stats.TotalWins)
	}

	if k.Drivers[0].Stats.AllTimeAverageFinish != 0 {
		t.Errorf("expected AllTimeAverageFinish to be 0, got %f", k.Drivers[0].Stats.AllTimeAverageFinish)
	}

	if len(k.Drivers[0].Stats.Last5Finish) != 0 {
		t.Errorf("expected Last5Finish to be empty, got %v", k.Drivers[0].Stats.Last5Finish)
	}

	if k.Drivers[0].Stats.PeakELO != InitialELO {
		t.Errorf("expected PeakELO to be %d, got %d", InitialELO, k.Drivers[0].Stats.PeakELO)
	}
}
