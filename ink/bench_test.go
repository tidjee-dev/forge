package ink

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

// ---------------------------------------------------------------------------
// Benchmark helpers
// ---------------------------------------------------------------------------

// initBenchColorMode forces AlwaysColorMode for the duration of the benchmark
// so that Render actually emits SGR sequences and Strip has real work to do.
func initBenchColorMode(b *testing.B) {
	b.Helper()
	old := globalColorMode.Load()
	b.Cleanup(func() {
		globalColorMode.Store(old)
		envOnce = sync.Once{}
		envDisabled = false
		envForced = false
	})
	SetGlobalColorMode(colorModeAlways)
}

// multilineText builds a string of n lines, each containing a short sentence.
func multilineText(n int) string {
	var sb strings.Builder
	for i := range n {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("Line %03d: the quick brown fox jumps over the lazy dog.", i+1))
	}
	return sb.String()
}

// ---------------------------------------------------------------------------
// BenchmarkRender_simple — single-line string, foreground + bold
// ---------------------------------------------------------------------------

func BenchmarkRender_simple(b *testing.B) {
	initBenchColorMode(b)

	s := New().
		WithBold(true).
		WithForeground(RGB(255, 128, 0))

	const text = "Hello, terminal!"

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_ = s.Render(text)
	}
}

// ---------------------------------------------------------------------------
// BenchmarkRender_multiline — 20-line string with padding + border
// ---------------------------------------------------------------------------

func BenchmarkRender_multiline(b *testing.B) {
	initBenchColorMode(b)

	s := New().
		WithBold(true).
		WithForeground(RGB(100, 200, 255)).
		WithBackground(RGB(20, 20, 40)).
		WithBorder(BorderRounded()).
		WithBorderColor(RGB(80, 80, 200)).
		WithLayout(
			NewLayout().
				WithUniformPadding(1).
				WithMinWidth(60).
				WithMaxWidth(80),
		)

	text := multilineText(20)

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_ = s.Render(text)
	}
}

// ---------------------------------------------------------------------------
// BenchmarkStrip — strip a pre-rendered 20-line block
// ---------------------------------------------------------------------------

func BenchmarkStrip(b *testing.B) {
	initBenchColorMode(b)

	s := New().
		WithBold(true).
		WithItalic(true).
		WithForeground(RGB(255, 64, 64)).
		WithBackground(RGB(20, 20, 20))

	rendered := s.Render(multilineText(20))

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_ = Strip(rendered)
	}
}

// ---------------------------------------------------------------------------
// BenchmarkTheme_Render — concurrent reads from a shared Theme
// ---------------------------------------------------------------------------

func BenchmarkTheme_Render(b *testing.B) {
	initBenchColorMode(b)

	th := NewTheme()
	th.Set("title", New().WithBold(true).WithForeground(RGB(100, 200, 255)))
	th.Set("error", New().WithBold(true).WithForeground(Danger))
	th.Set("warning", New().WithForeground(Warning))
	th.Set("success", New().WithForeground(Success))
	th.Set("muted", New().WithForeground(Muted).WithDim(true))

	keys := []string{"title", "error", "warning", "success", "muted"}
	texts := []string{
		"Hello, terminal!",
		"Something went wrong.",
		"Proceed with caution.",
		"Operation succeeded.",
		"This is a note.",
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			idx := i % len(keys)
			_ = th.Render(keys[idx], texts[idx])
			i++
		}
	})
}
