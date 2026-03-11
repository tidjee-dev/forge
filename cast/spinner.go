package cast

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/tidjee-dev/forge/ink"
)

// SpinnerFrames is the ordered set of glyphs that make up one full cycle of a
// spinner animation. Each element should occupy exactly one terminal column.
type SpinnerFrames []string

// Built-in frame sets.
var (
	// SpinnerDots is the classic braille-dot spinner (10 frames).
	SpinnerDots = SpinnerFrames{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	// SpinnerLine is a simple ASCII line spinner (4 frames).
	SpinnerLine = SpinnerFrames{"-", "\\", "|", "/"}

	// SpinnerCircle is a quarter-circle arc spinner (4 frames).
	SpinnerCircle = SpinnerFrames{"◐", "◓", "◑", "◒"}

	// SpinnerArrow is a rotating arrow spinner (8 frames).
	SpinnerArrow = SpinnerFrames{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}

	// SpinnerBounce is a bouncing-dot braille spinner (8 frames).
	SpinnerBounce = SpinnerFrames{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
)

// DefaultSpinnerInterval is the frame advance interval used by Spinner when no
// explicit interval is set.
const DefaultSpinnerInterval = 80 * time.Millisecond

// Spinner is a self-animating spinner that runs in the background.
// It advances through its frame set at a fixed interval, writing each frame
// to an [io.Writer] (default [os.Stderr]) using a carriage-return so the
// spinner overwrites itself on the same terminal line.
//
// The spinner is started with [Spinner.Start] and stopped with [Spinner.Stop].
// Stop blocks until the background goroutine exits and clears the spinner line
// before returning.
//
// Spinner is safe to call from multiple goroutines; Start and Stop are
// idempotent — calling Start on an already-running spinner or Stop on an
// already-stopped spinner is a no-op.
//
// Basic usage:
//
//	s := cast.NewSpinner().
//	    WithLabel("Loading…").
//	    WithStyle(ink.New().WithForeground(ink.Info))
//
//	s.Start()
//	defer s.Stop()
//
//	// … do work …
type Spinner struct {
	frames     SpinnerFrames
	label      string
	style      ink.Style
	labelStyle ink.Style
	interval   time.Duration
	writer     io.Writer

	mu      sync.Mutex
	running bool
	stop    chan struct{}
	done    chan struct{}
}

// NewSpinner returns a ready-to-use Spinner with SpinnerDots frames,
// [DefaultSpinnerInterval] tick rate, and output directed to [os.Stderr].
// Call [Spinner.Start] to begin animation.
func NewSpinner() *Spinner {
	return &Spinner{
		frames:   SpinnerDots,
		interval: DefaultSpinnerInterval,
		writer:   os.Stderr,
	}
}

// ---------------------------------------------------------------------------
// Fluent configuration
//
// Setters must be called before Start; calling them on a running spinner has
// no effect until the next Start.
// ---------------------------------------------------------------------------

// WithFrames sets the frame set used for animation. If f is empty the
// previous frame set is retained. Returns the receiver for chaining.
func (s *Spinner) WithFrames(f SpinnerFrames) *Spinner {
	if len(f) > 0 {
		cp := make(SpinnerFrames, len(f))
		copy(cp, f)
		s.frames = cp
	}
	return s
}

// WithLabel sets the text displayed alongside the spinner glyph. Pass an
// empty string to remove the label. Returns the receiver for chaining.
func (s *Spinner) WithLabel(text string) *Spinner {
	s.label = text
	return s
}

// WithStyle sets the [ink.Style] applied to the spinner glyph (and to the
// label when no separate [Spinner.WithLabelStyle] is set). Returns the
// receiver for chaining.
func (s *Spinner) WithStyle(st ink.Style) *Spinner {
	s.style = st
	return s
}

// WithLabelStyle sets the [ink.Style] applied to the label text only. When
// not set the label inherits [Spinner.WithStyle]. Returns the receiver for
// chaining.
func (s *Spinner) WithLabelStyle(st ink.Style) *Spinner {
	s.labelStyle = st
	return s
}

// WithInterval sets the frame advance interval. Values ≤ 0 are ignored.
// Returns the receiver for chaining.
func (s *Spinner) WithInterval(d time.Duration) *Spinner {
	if d > 0 {
		s.interval = d
	}
	return s
}

// WithWriter redirects spinner output to w instead of [os.Stderr]. This is
// useful for directing output to [os.Stdout] or a buffer in tests. Returns
// the receiver for chaining.
func (s *Spinner) WithWriter(w io.Writer) *Spinner {
	if w != nil {
		s.writer = w
	}
	return s
}

// ---------------------------------------------------------------------------
// Lifecycle
// ---------------------------------------------------------------------------

// Start begins the animation loop in a background goroutine. If the spinner
// is already running, Start is a no-op.
//
// The goroutine writes a frame to the configured writer on every tick. Each
// write is preceded by "\r" so the frame overwrites the previous one on the
// same terminal line.
func (s *Spinner) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return
	}

	frames := s.frames
	if len(frames) == 0 {
		frames = SpinnerDots
	}

	interval := s.interval
	if interval <= 0 {
		interval = DefaultSpinnerInterval
	}

	// Snapshot configuration so the goroutine is not affected by concurrent
	// setter calls after Start.
	label := s.label
	style := s.style
	labelStyle := s.labelStyle
	writer := s.writer

	s.stop = make(chan struct{})
	s.done = make(chan struct{})
	s.running = true

	go func() {
		defer close(s.done)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		frameIdx := 0

		for {
			select {
			case <-s.stop:
				spinnerClearLine(writer, label)
				return
			case <-ticker.C:
				line := spinnerRenderFrame(frames, frameIdx, label, style, labelStyle)
				fmt.Fprintf(writer, "\r%s", line)
				frameIdx = (frameIdx + 1) % len(frames)
			}
		}
	}()
}

// Stop halts the animation, clears the spinner line, and blocks until the
// background goroutine has fully exited. If the spinner is not running, Stop
// is a no-op.
func (s *Spinner) Stop() {
	s.mu.Lock()

	if !s.running {
		s.mu.Unlock()
		return
	}

	close(s.stop)
	done := s.done
	s.running = false

	s.mu.Unlock()

	// Wait for the goroutine to finish without holding the mutex.
	<-done
}

// IsRunning reports whether the animation goroutine is currently active.
func (s *Spinner) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// spinnerRenderFrame builds the display string for a single animation frame.
func spinnerRenderFrame(frames SpinnerFrames, idx int, label string, style, labelStyle ink.Style) string {
	if len(frames) == 0 {
		return ""
	}

	i := idx % len(frames)
	if i < 0 {
		i += len(frames)
	}

	glyph := frames[i]
	styledGlyph := style.Render(glyph)

	if label == "" {
		return styledGlyph
	}

	var styledLabel string
	if !isStyleSet(labelStyle) {
		styledLabel = style.Render(label)
	} else {
		styledLabel = labelStyle.Render(label)
	}

	return styledGlyph + " " + styledLabel
}

// spinnerClearLine writes a carriage-return followed by enough spaces to
// overwrite the last rendered spinner line, then returns the cursor to the
// start of the line.
func spinnerClearLine(w io.Writer, label string) {
	// One glyph (at most 2 cols) + 1 space + label.
	clearWidth := 2
	if label != "" {
		clearWidth += 1 + visibleWidth(label)
	}
	fmt.Fprintf(w, "\r%s\r", repeatStr(" ", clearWidth))
}
