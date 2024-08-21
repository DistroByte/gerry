package main

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/DistroByte/gerry/internal/config"
)

func init() {
	os.Setenv("TEST", "true")
}

func TestMain_NoArgs(t *testing.T) {
	cmd := exec.Command("go", "run", ".")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	expected := "no arguments provided\n"
	if string(output) != expected {
		t.Errorf("expected %q, got %q", expected, string(output))
	}
}

func TestMain_ConfGen(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "confgen", "testconfig")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, err := os.Stat("testconfig"); os.IsNotExist(err) {
		t.Errorf("expected config file to be created, got %v", err)
	}

	if output := string(output); output != "config file generated successfully\n" {
		t.Errorf("expected %q, got %q", "config file generated successfully\n", output)
	}

	// Clean up
	if err := os.Remove("testconfig"); err != nil {
		t.Fatalf("failed to remove testconfig: %v", err)
	}
}

func TestMain_Start_NoConfig(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "start", "")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	expected := "no config file provided\n"
	if string(output) != expected {
		t.Errorf("expected %q, got %q", expected, string(output))
	}
}

func TestMain_Start_WithConfig(t *testing.T) {
	exec.Command("go", "run", ".", "confgen", "testconfig")
	cmd := exec.Command("go", "run", ".", "start", "testconfig")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Add assertions based on what bot.Start outputs or does

	// app environment is test, aborting startup
	if output := string(output); output != "app environment is test, aborting startup\n" {
		t.Errorf("expected %q, got %q", "app environment is test, aborting startup\n", output)
	}

	// Clean up
	if err := os.Remove("testconfig"); err != nil {
		t.Fatalf("failed to remove testconfig: %v", err)
	}
}

func TestMain_Start_Time(t *testing.T) {
	// Assuming "testconfig" is a valid config file for the bot
	exec.Command("go", "run", ".", "confgen", "testconfig")
	cmd := exec.Command("go", "run", ".", "start", "testconfig")
	start := time.Now()
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if config.StartTime.Before(start) {
		t.Errorf("expected StartTime to be after %v, got %v", start, config.StartTime)
	}
}
