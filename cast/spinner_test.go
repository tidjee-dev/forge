package cast_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/tidjee-dev/forge/cast"
	"github.com/tidjee-dev/forge/ink"
)

func TestSpinner_NotRunningBeforeStart(t *testing.T) {
	s := cast.NewSpinner()
	if s.IsRunning() {
		t.Error("expected spinner to not be running before Start()")
	}
}

func TestSpinner_IsRunningAfterStart(t *testing.T) {
	var buf bytes.Buffer
	s := cast.NewSpinner().
		WithWriter(&buf).
		WithInterval(10 * time.Millisecond)

	s.Start()
	defer s.Stop()

	if !s.IsRunning() {
		t.Error("expected spinner to be running after Start()")
	}
}

func TestSpinner_NotRunningAfterStop(t *testing.T) {
	var buf bytes.Buffer
	s := cast.NewSpinner().
		WithWriter(&buf).
		WithInterval(10 * time.Millisecond)

	s.Start()
	s.Stop()

	if s.IsRunning() {
		t.Error("expected spinner to not be running after Stop()")
	}
}

func TestSpinner_WritesOutput(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	var buf bytes.Buffer
	s := cast.NewSpinner().
		WithFrames(cast.SpinnerDots).
		WithWriter(&buf).
		WithInterval(10 * time.Millisecond)

	s.Start()
	time.Sleep(50 * time.Millisecond)
	s.Stop()

	got := buf.String()
	if got == "" {
		t.Error("expected output after animation, got empty string")
	}
	if !strings.Contains(got, "\r") {
		t.Errorf("expected carriage-return in output, got %q", got)
	}
}

func TestSpinner_WritesFrameGlyphs(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	var buf bytes.Buffer
	s := cast.NewSpinner().
		WithFrames(cast.SpinnerDots).
		WithWriter(&buf).
		WithInterval(10 * time.Millisecond)

	s.Start()
	time.Sleep(60 * time.Millisecond)
	s.Stop()

	got := buf.String()
	foundFrame := false
	for _, frame := range cast.SpinnerDots {
		if strings.Contains(got, frame) {
			foundFrame = true
			break
		}
	}
	if !foundFrame {
		t.Errorf("expected at least one SpinnerDots frame in output, got %q", got)
	}
}

func TestSpinner_WritesLabelInOutput(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	var buf bytes.Buffer
	s := cast.NewSpinner().
		WithWriter(&buf).
		WithLabel("Loading…").
		WithInterval(10 * time.Millisecond)

	s.Start()
	time.Sleep(50 * time.Millisecond)
	s.Stop()

	got := buf.String()
	if !strings.Contains(got, "Loading…") {
		t.Errorf("expected label \"Loading…\" in output, got %q", got)
	}
}

func TestSpinner_StartIdempotent(t *testing.T) {
	var buf bytes.Buffer
	s := cast.NewSpinner().
		WithWriter(&buf).
		WithInterval(10 * time.Millisecond)

	s.Start()
	s.Start() // second Start should be a no-op
	defer s.Stop()

	if !s.IsRunning() {
		t.Error("expected spinner to still be running after double Start()")
	}
}

func TestSpinner_StopIdempotent(t *testing.T) {
	var buf bytes.Buffer
	s := cast.NewSpinner().
		WithWriter(&buf).
		WithInterval(10 * time.Millisecond)

	s.Start()
	s.Stop()
	s.Stop() // second Stop should be a no-op, not panic

	if s.IsRunning() {
		t.Error("expected spinner to not be running after double Stop()")
	}
}

func TestSpinner_StopWithoutStart(t *testing.T) {
	// Stop on a never-started spinner should not panic.
	s := cast.NewSpinner()
	s.Stop()
	if s.IsRunning() {
		t.Error("expected IsRunning() == false on never-started spinner")
	}
}

func TestSpinner_ClearsLineOnStop(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	var buf bytes.Buffer
	s := cast.NewSpinner().
		WithWriter(&buf).
		WithLabel("work").
		WithInterval(10 * time.Millisecond)

	s.Start()
	time.Sleep(30 * time.Millisecond)
	s.Stop()

	got := buf.String()
	// The clear sequence ends with \r, returning the cursor to line start.
	if !strings.HasSuffix(got, "\r") {
		t.Errorf("expected output to end with \\r after Stop (clear sequence), got %q", got)
	}
}

func TestSpinner_CanRestartAfterStop(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	var buf bytes.Buffer
	s := cast.NewSpinner().
		WithWriter(&buf).
		WithInterval(10 * time.Millisecond)

	s.Start()
	time.Sleep(20 * time.Millisecond)
	s.Stop()

	if s.IsRunning() {
		t.Fatal("expected spinner stopped after first cycle")
	}

	buf.Reset()

	s.Start()
	time.Sleep(20 * time.Millisecond)
	s.Stop()

	if s.IsRunning() {
		t.Error("expected spinner stopped after second cycle")
	}
	if buf.Len() == 0 {
		t.Error("expected output on second Start/Stop cycle")
	}
}

func TestSpinner_CustomFrames(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	custom := cast.SpinnerFrames{"A", "B", "C"}
	var buf bytes.Buffer
	s := cast.NewSpinner().
		WithFrames(custom).
		WithWriter(&buf).
		WithInterval(10 * time.Millisecond)

	s.Start()
	time.Sleep(50 * time.Millisecond)
	s.Stop()

	got := buf.String()
	foundCustom := strings.Contains(got, "A") ||
		strings.Contains(got, "B") ||
		strings.Contains(got, "C")
	if !foundCustom {
		t.Errorf("expected custom frame glyphs in output, got %q", got)
	}
}

func TestSpinner_EmptyFramesFallbackToSpinnerDots(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	var buf bytes.Buffer
	s := cast.NewSpinner().
		WithFrames(cast.SpinnerFrames{}). // empty — should retain SpinnerDots
		WithWriter(&buf).
		WithInterval(10 * time.Millisecond)

	s.Start()
	time.Sleep(50 * time.Millisecond)
	s.Stop()

	got := buf.String()
	foundFrame := false
	for _, frame := range cast.SpinnerDots {
		if strings.Contains(got, frame) {
			foundFrame = true
			break
		}
	}
	if !foundFrame {
		t.Errorf("expected SpinnerDots fallback frames in output, got %q", got)
	}
}

func TestSpinner_WithStyle(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	var buf bytes.Buffer
	s := cast.NewSpinner().
		WithWriter(&buf).
		WithStyle(ink.New().WithForeground(ink.Info)).
		WithInterval(10 * time.Millisecond)

	s.Start()
	time.Sleep(30 * time.Millisecond)
	s.Stop()

	got := buf.String()
	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI sequences in styled output, got %q", got)
	}
}

func TestSpinner_WithLabelStyle(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	var buf bytes.Buffer
	s := cast.NewSpinner().
		WithWriter(&buf).
		WithLabel("Working").
		WithStyle(ink.New().WithForeground(ink.Info)).
		WithLabelStyle(ink.New().WithForeground(ink.Muted)).
		WithInterval(10 * time.Millisecond)

	s.Start()
	time.Sleep(30 * time.Millisecond)
	s.Stop()

	got := buf.String()
	if !strings.Contains(got, "Working") {
		t.Errorf("expected label text in output, got %q", got)
	}
	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI sequences in styled output, got %q", got)
	}
}

func TestSpinner_NoColorMode(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	var buf bytes.Buffer
	s := cast.NewSpinner().
		WithWriter(&buf).
		WithLabel("Loading").
		WithStyle(ink.New().WithForeground(ink.Danger)).
		WithInterval(10 * time.Millisecond)

	s.Start()
	time.Sleep(30 * time.Millisecond)
	s.Stop()

	got := buf.String()
	if strings.Contains(got, "\x1b[") {
		t.Errorf("expected no ANSI sequences in NeverColorMode, got %q", got)
	}
	if !strings.Contains(got, "Loading") {
		t.Errorf("expected label text in plain output, got %q", got)
	}
}

func TestSpinner_DefaultInterval(t *testing.T) {
	if cast.DefaultSpinnerInterval <= 0 {
		t.Errorf("DefaultSpinnerInterval must be positive, got %v", cast.DefaultSpinnerInterval)
	}
	if cast.DefaultSpinnerInterval > 500*time.Millisecond {
		t.Errorf("DefaultSpinnerInterval seems too large: %v", cast.DefaultSpinnerInterval)
	}
}

func TestSpinner_BuiltInFrameSetsNonEmpty(t *testing.T) {
	sets := map[string]cast.SpinnerFrames{
		"SpinnerDots":   cast.SpinnerDots,
		"SpinnerLine":   cast.SpinnerLine,
		"SpinnerCircle": cast.SpinnerCircle,
		"SpinnerArrow":  cast.SpinnerArrow,
		"SpinnerBounce": cast.SpinnerBounce,
	}
	for name, frames := range sets {
		if len(frames) == 0 {
			t.Errorf("%s is empty", name)
		}
	}
}

func TestSpinner_AllBuiltInFrameSetsAnimate(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	sets := []struct {
		name   string
		frames cast.SpinnerFrames
	}{
		{"SpinnerDots", cast.SpinnerDots},
		{"SpinnerLine", cast.SpinnerLine},
		{"SpinnerCircle", cast.SpinnerCircle},
		{"SpinnerArrow", cast.SpinnerArrow},
		{"SpinnerBounce", cast.SpinnerBounce},
	}

	for _, tc := range sets {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			s := cast.NewSpinner().
				WithFrames(tc.frames).
				WithWriter(&buf).
				WithInterval(10 * time.Millisecond)

			s.Start()
			time.Sleep(50 * time.Millisecond)
			s.Stop()

			got := buf.String()
			foundFrame := false
			for _, frame := range tc.frames {
				if strings.Contains(got, frame) {
					foundFrame = true
					break
				}
			}
			if !foundFrame {
				t.Errorf("%s: expected at least one frame glyph in output, got %q", tc.name, got)
			}
		})
	}
}
