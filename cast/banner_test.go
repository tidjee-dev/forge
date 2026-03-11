package cast_test

import (
	"strings"
	"testing"

	"github.com/tidjee-dev/forge/cast"
	"github.com/tidjee-dev/forge/ink"
)

func TestBanner_RenderPlain(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewBanner("Hello, world!").Render()
	if got != "Hello, world!" {
		t.Errorf("expected plain text, got %q", got)
	}
}

func TestBanner_RenderEmpty(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewBanner("").Render()
	if got != "" {
		t.Errorf("expected empty string for empty banner, got %q", got)
	}
}

func TestBanner_WithBorder(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewBanner("Hello").Border(ink.BorderRounded()).Render()
	lines := strings.Split(got, "\n")

	if len(lines) != 3 {
		t.Fatalf("expected 3 lines with border, got %d: %q", len(lines), got)
	}
	// Top line starts with ╭ and ends with ╮
	if !strings.HasPrefix(lines[0], "╭") || !strings.HasSuffix(lines[0], "╮") {
		t.Errorf("top border malformed: %q", lines[0])
	}
	// Content line starts with │ and ends with │
	if !strings.HasPrefix(lines[1], "│") || !strings.HasSuffix(lines[1], "│") {
		t.Errorf("content border malformed: %q", lines[1])
	}
	// Bottom line starts with ╰ and ends with ╯
	if !strings.HasPrefix(lines[2], "╰") || !strings.HasSuffix(lines[2], "╯") {
		t.Errorf("bottom border malformed: %q", lines[2])
	}
	// Content line contains the text
	if !strings.Contains(lines[1], "Hello") {
		t.Errorf("content line missing text: %q", lines[1])
	}
}

func TestBanner_WithBorderNormal(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewBanner("Test").Border(ink.BorderNormal()).Render()
	lines := strings.Split(got, "\n")

	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "┌") || !strings.HasSuffix(lines[0], "┐") {
		t.Errorf("top border malformed: %q", lines[0])
	}
	if !strings.HasPrefix(lines[2], "└") || !strings.HasSuffix(lines[2], "┘") {
		t.Errorf("bottom border malformed: %q", lines[2])
	}
}

func TestBanner_Width(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Without border: content is padded to 20 columns.
	got := cast.NewBanner("Hi").Width(20).Render()
	w := visibleWidthTest(got)
	if w != 20 {
		t.Errorf("expected width 20, got %d: %q", w, got)
	}
}

func TestBanner_WidthWithBorder(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// With border: overall width (border + content) should equal Width.
	got := cast.NewBanner("Hi").Border(ink.BorderRounded()).Width(20).Render()
	lines := strings.Split(got, "\n")
	if len(lines) < 1 {
		t.Fatal("no output lines")
	}
	topW := visibleWidthTest(lines[0])
	if topW != 20 {
		t.Errorf("expected top border width 20, got %d: %q", topW, lines[0])
	}
}

func TestBanner_AlignCenter(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewBanner("Hi").Width(10).Align(ink.JustifyCenter).Render()
	// "Hi" is 2 wide; 10-2=8 padding; 4 left, 4 right → "    Hi    "
	w := visibleWidthTest(got)
	if w != 10 {
		t.Errorf("expected width 10, got %d: %q", w, got)
	}
	// Should be centred: roughly equal whitespace on both sides
	stripped := got
	leftSpaces := len(stripped) - len(strings.TrimLeft(stripped, " "))
	rightSpaces := len(stripped) - len(strings.TrimRight(stripped, " "))
	diff := leftSpaces - rightSpaces
	if diff < -1 || diff > 1 {
		t.Errorf("text does not appear centred: left=%d right=%d in %q", leftSpaces, rightSpaces, got)
	}
}

func TestBanner_AlignEnd(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewBanner("Hi").Width(10).Align(ink.JustifyEnd).Render()
	w := visibleWidthTest(got)
	if w != 10 {
		t.Errorf("expected width 10, got %d: %q", w, got)
	}
	if !strings.HasSuffix(got, "Hi") {
		t.Errorf("expected text at end, got %q", got)
	}
}

func TestBanner_AlignStart(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewBanner("Hi").Width(10).Align(ink.JustifyStart).Render()
	w := visibleWidthTest(got)
	if w != 10 {
		t.Errorf("expected width 10, got %d: %q", w, got)
	}
	if !strings.HasPrefix(got, "Hi") {
		t.Errorf("expected text at start, got %q", got)
	}
}

func TestBanner_Style(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	s := ink.New().WithForeground(ink.BrightWhite).WithBackground(ink.Blue).WithBold(true)
	got := cast.NewBanner("Server started").Style(s).Render()

	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI sequences, got %q", got)
	}
	plain := ink.Strip(got)
	if plain != "Server started" {
		t.Errorf("expected plain text %q, got %q", "Server started", plain)
	}
}

func TestBanner_StyleWithBorder(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	s := ink.New().WithForeground(ink.White).WithBold(true)
	got := cast.NewBanner("Test").Style(s).Border(ink.BorderRounded()).Render()
	lines := strings.Split(got, "\n")

	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %q", len(lines), got)
	}
	// Top border should have no ANSI (border itself is unstyled unless a
	// border colour is set separately).
	if strings.Contains(lines[0], "\x1b[") {
		t.Errorf("top border should be unstyled, got %q", lines[0])
	}
	// Content line should contain ANSI.
	if !strings.Contains(lines[1], "\x1b[") {
		t.Errorf("content line should contain ANSI styling, got %q", lines[1])
	}
}

func TestBanner_NoColorMode(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	s := ink.New().WithForeground(ink.Blue).WithBold(true)
	got := cast.NewBanner("No color").Style(s).Render()

	if strings.Contains(got, "\x1b[") {
		t.Errorf("expected no ANSI in NeverColorMode, got %q", got)
	}
	if got != "No color" {
		t.Errorf("expected plain text, got %q", got)
	}
}

func TestBanner_Immutability(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	base := cast.NewBanner("Base")
	withBorder := base.Border(ink.BorderRounded())

	// base should still render without a border (single line)
	baseLines := strings.Split(base.Render(), "\n")
	if len(baseLines) != 1 {
		t.Errorf("base banner should have no border (1 line), got %d lines", len(baseLines))
	}

	// withBorder should have 3 lines
	borderLines := strings.Split(withBorder.Render(), "\n")
	if len(borderLines) != 3 {
		t.Errorf("bordered banner should have 3 lines, got %d", len(borderLines))
	}
}

func TestBanner_WidthZero(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Width(0) → auto: content is as wide as the text.
	got := cast.NewBanner("Hello").Width(0).Render()
	if got != "Hello" {
		t.Errorf("expected \"Hello\", got %q", got)
	}
}

// ---------------------------------------------------------------------------
// Helper
// ---------------------------------------------------------------------------

// visibleWidthTest is a test-local visible-width helper that avoids importing
// cast internals. It counts visible columns after stripping ANSI sequences.
func visibleWidthTest(s string) int {
	plain := ink.Strip(s)
	w := 0
	for _, r := range plain {
		switch {
		case r < 0x20, r == 0x7F:
			// control
		case r >= 0x1100 && r <= 0x115F,
			r >= 0x2E80 && r <= 0x303E,
			r >= 0x3040 && r <= 0x33FF,
			r >= 0x4E00 && r <= 0x9FFF,
			r >= 0xAC00 && r <= 0xD7AF,
			r >= 0xFF01 && r <= 0xFF60:
			w += 2
		default:
			w++
		}
	}
	return w
}
