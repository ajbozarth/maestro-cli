package main

import (
	"os"
	"testing"
	"time"
)

func TestNewProgressIndicator(t *testing.T) {
	indicator := NewProgressIndicator("Test message")

	if indicator.message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", indicator.message)
	}

	if len(indicator.spinner) == 0 {
		t.Error("Spinner should have animation frames")
	}

	if indicator.isRunning {
		t.Error("Progress indicator should not be running initially")
	}
}

func TestProgressIndicatorStart(t *testing.T) {
	indicator := NewProgressIndicator("Test message")

	indicator.Start()

	if !indicator.isRunning {
		t.Error("Progress indicator should be running after Start()")
	}

	if indicator.startTime.IsZero() {
		t.Error("Start time should be set")
	}
}

func TestProgressIndicatorStop(t *testing.T) {
	indicator := NewProgressIndicator("Test message")

	indicator.Start()
	time.Sleep(10 * time.Millisecond) // Small delay to ensure different times
	indicator.Stop("Completed")

	if indicator.isRunning {
		t.Error("Progress indicator should not be running after Stop()")
	}
}

func TestProgressIndicatorStopWithError(t *testing.T) {
	indicator := NewProgressIndicator("Test message")

	indicator.Start()
	time.Sleep(10 * time.Millisecond) // Small delay to ensure different times
	indicator.StopWithError("Failed")

	if indicator.isRunning {
		t.Error("Progress indicator should not be running after StopWithError()")
	}
}

func TestNewProgressBar(t *testing.T) {
	bar := NewProgressBar("Test message", 100)

	if bar.message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", bar.message)
	}

	if bar.total != 100 {
		t.Errorf("Expected total 100, got %d", bar.total)
	}

	if bar.current != 0 {
		t.Errorf("Expected current 0, got %d", bar.current)
	}

	if bar.isRunning {
		t.Error("Progress bar should not be running initially")
	}
}

func TestProgressBarStart(t *testing.T) {
	bar := NewProgressBar("Test message", 100)

	bar.Start()

	if !bar.isRunning {
		t.Error("Progress bar should be running after Start()")
	}

	if bar.startTime.IsZero() {
		t.Error("Start time should be set")
	}
}

func TestProgressBarUpdate(t *testing.T) {
	bar := NewProgressBar("Test message", 100)

	bar.Start()
	bar.Update(50)

	if bar.current != 50 {
		t.Errorf("Expected current 50, got %d", bar.current)
	}
}

func TestProgressBarIncrement(t *testing.T) {
	bar := NewProgressBar("Test message", 100)

	bar.Start()
	bar.Increment()

	if bar.current != 1 {
		t.Errorf("Expected current 1, got %d", bar.current)
	}
}

func TestProgressBarStop(t *testing.T) {
	bar := NewProgressBar("Test message", 100)

	bar.Start()
	time.Sleep(10 * time.Millisecond) // Small delay to ensure different times
	bar.Stop("Completed")

	if bar.isRunning {
		t.Error("Progress bar should not be running after Stop()")
	}

	if bar.current != bar.total {
		t.Error("Progress bar should be at 100% when stopped")
	}
}

func TestProgressBarStopWithError(t *testing.T) {
	bar := NewProgressBar("Test message", 100)

	bar.Start()
	time.Sleep(10 * time.Millisecond) // Small delay to ensure different times
	bar.StopWithError("Failed")

	if bar.isRunning {
		t.Error("Progress bar should not be running after StopWithError()")
	}
}

func TestShouldShowProgress(t *testing.T) {
	// Test with test mode enabled
	os.Setenv("MAESTRO_TEST_MODE", "true")
	if ShouldShowProgress() {
		t.Error("Should not show progress in test mode")
	}

	// Test with test mode disabled
	os.Setenv("MAESTRO_TEST_MODE", "false")
	// Note: We can't easily test the terminal check in unit tests
	// The actual behavior will be tested in integration tests
}

func TestProgressIndicatorTick(t *testing.T) {
	indicator := NewProgressIndicator("Test message")

	indicator.Start()
	initialSpinner := indicator.current

	indicator.Tick()

	if indicator.current == initialSpinner {
		t.Error("Spinner should advance after Tick()")
	}
}

func TestProgressIndicatorUpdate(t *testing.T) {
	indicator := NewProgressIndicator("Test message")

	indicator.Start()
	indicator.Update("Updated message")

	// The update should work without errors
	if !indicator.isRunning {
		t.Error("Progress indicator should still be running after Update()")
	}
}

func TestProgressBarUpdateExceedsTotal(t *testing.T) {
	bar := NewProgressBar("Test message", 100)

	bar.Start()
	bar.Update(150) // Exceeds total

	if bar.current != 100 {
		t.Errorf("Current should be capped at total, got %d", bar.current)
	}
}

func TestProgressBarIncrementExceedsTotal(t *testing.T) {
	bar := NewProgressBar("Test message", 1)

	bar.Start()
	bar.Increment() // Should reach 1
	bar.Increment() // Should stay at 1

	if bar.current != 1 {
		t.Errorf("Current should be capped at total, got %d", bar.current)
	}
}
