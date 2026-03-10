package ink

import (
	"os"
	"sync"
	"testing"
)

// resetColorMode restores the global colour mode and clears the env-detection
// cache after each test so that tests don't bleed state into one another.
func resetColorMode(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		globalColorMode.Store(uint32(colorModeAuto))
		// Reset the once so auto-detection re-runs in subsequent tests.
		envOnce = sync.Once{}
		envDisabled = false
		envForced = false
	})
}

func TestColorMode_NeverReturnsFalse(t *testing.T) {
	resetColorMode(t)
	SetGlobalColorMode(colorModeNever)
	if ColorModeEnabled() {
		t.Error("ColorModeEnabled() = true with NeverColorMode, want false")
	}
}

func TestColorMode_AlwaysReturnsTrue(t *testing.T) {
	resetColorMode(t)
	SetGlobalColorMode(colorModeAlways)
	if !ColorModeEnabled() {
		t.Error("ColorModeEnabled() = false with AlwaysColorMode, want true")
	}
}

func TestColorMode_AlwaysIgnoresTTY(t *testing.T) {
	resetColorMode(t)
	// Even if NO_COLOR is set, AlwaysColorMode overrides it.
	t.Setenv("NO_COLOR", "1")
	envOnce = sync.Once{}
	SetGlobalColorMode(colorModeAlways)
	if !ColorModeEnabled() {
		t.Error("ColorModeEnabled() = false with AlwaysColorMode + NO_COLOR, want true")
	}
}

func TestColorMode_AutoWithNOCOLOR(t *testing.T) {
	resetColorMode(t)
	t.Setenv("NO_COLOR", "1")
	// Reset once so initEnv picks up the new env value.
	envOnce = sync.Once{}
	envDisabled = false
	envForced = false
	SetGlobalColorMode(colorModeAuto)
	if ColorModeEnabled() {
		t.Error("ColorModeEnabled() = true with NO_COLOR set, want false")
	}
}

func TestColorMode_AutoWithTERMDumb(t *testing.T) {
	resetColorMode(t)
	// Ensure NO_COLOR is absent so only TERM=dumb takes effect.
	os.Unsetenv("NO_COLOR")
	t.Setenv("TERM", "dumb")
	envOnce = sync.Once{}
	envDisabled = false
	envForced = false
	SetGlobalColorMode(colorModeAuto)
	if ColorModeEnabled() {
		t.Error("ColorModeEnabled() = true with TERM=dumb, want false")
	}
}

func TestColorMode_PublicAliases(t *testing.T) {
	if AutoColorMode != colorModeAuto {
		t.Errorf("AutoColorMode = %v, want %v", AutoColorMode, colorModeAuto)
	}
	if AlwaysColorMode != colorModeAlways {
		t.Errorf("AlwaysColorMode = %v, want %v", AlwaysColorMode, colorModeAlways)
	}
	if NeverColorMode != colorModeNever {
		t.Errorf("NeverColorMode = %v, want %v", NeverColorMode, colorModeNever)
	}
}

// TestColorMode_ConcurrentSetGet verifies that concurrent calls to
// SetGlobalColorMode and ColorModeEnabled do not race.
// Run with: go test -race ./...
func TestColorMode_ConcurrentSetGet(t *testing.T) {
	resetColorMode(t)

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	modes := []outputColorMode{colorModeAuto, colorModeAlways, colorModeNever}

	// Writers
	for i := range goroutines {
		go func(i int) {
			defer wg.Done()
			SetGlobalColorMode(modes[i%len(modes)])
		}(i)
	}

	// Readers
	for range goroutines {
		go func() {
			defer wg.Done()
			_ = ColorModeEnabled()
		}()
	}

	wg.Wait()
}
