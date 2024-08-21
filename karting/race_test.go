package karting

import (
	"testing"
)

func TestRace_AllDriversRegistered(t *testing.T) {
	k := &Karting{
		Drivers: []*Driver{
			{Name: "Driver1", ELO: 1000, Stats: &DriverStats{}},
			{Name: "Driver2", ELO: 1000, Stats: &DriverStats{}},
		},
	}

	results := []*Result{
		{Driver: &Driver{Name: "Driver1"}, Position: 1},
		{Driver: &Driver{Name: "Driver2"}, Position: 2},
	}

	expectedELOChanges := []int{16, -15}

	raceDiff, err := k.Race(results)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for i, diff := range raceDiff {
		if diff.Change != expectedELOChanges[i] {
			t.Errorf("expected ELO change %d, got %d", expectedELOChanges[i], diff.Change)
		}
	}
}

func TestRace_DriverNotRegistered(t *testing.T) {
	k := &Karting{
		Drivers: []*Driver{
			{Name: "Driver1", ELO: 1000, Stats: &DriverStats{}},
		},
	}

	results := []*Result{
		{Driver: &Driver{Name: "Driver1"}, Position: 1},
		{Driver: &Driver{Name: "Driver2"}, Position: 2},
	}

	_, err := k.Race(results)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	expectedError := "driver \"Driver2\" not found. not recording race"
	if err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestRace_CorrectELOChanges(t *testing.T) {
	k := &Karting{
		Drivers: []*Driver{
			{Name: "Driver1", ELO: 1000, Stats: &DriverStats{}},
			{Name: "Driver2", ELO: 1000, Stats: &DriverStats{}},
			{Name: "Driver3", ELO: 1000, Stats: &DriverStats{}},
		},
	}

	results := []*Result{
		{Driver: &Driver{Name: "Driver1"}, Position: 1},
		{Driver: &Driver{Name: "Driver2"}, Position: 2},
		{Driver: &Driver{Name: "Driver3"}, Position: 3},
	}

	expectedELOChanges := []int{16, 0, -16}

	raceDiff, err := k.Race(results)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for i, diff := range raceDiff {
		if diff.Change != expectedELOChanges[i] {
			t.Errorf("expected ELO change %d, got %d", expectedELOChanges[i], diff.Change)
		}
	}
}
