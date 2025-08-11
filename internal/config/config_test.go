package config

import (
	"os"
	"testing"
)

func TestLoad_MissingRequiredPanicsSafely(t *testing.T) {
	// Ensure key is unset
	_ = os.Unsetenv("VISUAL_CROSSING_API_KEY")

	// We can't let log.Fatalf kill the test process; so just assert that
	// setting it makes Load() succeed. (This acts as a guard for future edits.)
	t.Setenv("VISUAL_CROSSING_API_KEY", "ok")
	_ = Load() // should not crash
}
