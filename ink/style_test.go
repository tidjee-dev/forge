package ink

import (
	"reflect"
	"strings"
	"sync"
	"testing"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// forceColorMode sets the global colour mode for the duration of the test and
// restores it (along with the env-detection cache) on cleanup.
func forceColorMode(t *testing.T, m outputColorMode) {
	t.Helper()
	t.Cleanup(func() {
		globalColorMode.Store(uint32(colorModeAuto))
		envOnce = sync.Once{}
		envDisabled = false
		envForced = false
	})
	SetGlobalColorMode(m)
}

// containsSGR reports whether s contains at least one ESC[ sequence.
func containsSGR(s string) bool {
	return strings.ContainsRune(s, '\x1b')
}

// ---------------------------------------------------------------------------
// Constructor / zero value
// ---------------------------------------------------------------------------

func TestStyle_New_IsZero(t *testing.T) {
	s := New()
	if s.IsBold() || s.IsItalic() || s.IsUnderline() || s.IsDim() ||
		s.IsStrikethrough() || s.IsInverse() {
		t.Error("New() returned a non-zero style")
	}
	if !s.GetForeground().IsZeroColor() {
		t.Error("New() fg should be zero color")
	}
	if !s.GetBackground().IsZeroColor() {
		t.Error("New() bg should be zero color")
	}
}

// ---------------------------------------------------------------------------
// Attribute setters / getters
// ---------------------------------------------------------------------------

func TestStyle_WithBold(t *testing.T) {
	s := New().WithBold(true)
	if !s.IsBold() {
		t.Error("IsBold() = false after WithBold(true)")
	}
	s2 := s.WithBold(false)
	if s2.IsBold() {
		t.Error("IsBold() = true after WithBold(false)")
	}
	// Original must be unchanged (immutability).
	if !s.IsBold() {
		t.Error("WithBold(false) mutated the original style")
	}
}

func TestStyle_WithItalic(t *testing.T) {
	s := New().WithItalic(true)
	if !s.IsItalic() {
		t.Error("IsItalic() = false after WithItalic(true)")
	}
	if New().WithItalic(false).IsItalic() {
		t.Error("IsItalic() = true after WithItalic(false)")
	}
}

func TestStyle_WithUnderline(t *testing.T) {
	s := New().WithUnderline(true)
	if !s.IsUnderline() {
		t.Error("IsUnderline() = false after WithUnderline(true)")
	}
}

func TestStyle_WithDim(t *testing.T) {
	s := New().WithDim(true)
	if !s.IsDim() {
		t.Error("IsDim() = false after WithDim(true)")
	}
}

func TestStyle_WithStrikethrough(t *testing.T) {
	s := New().WithStrikethrough(true)
	if !s.IsStrikethrough() {
		t.Error("IsStrikethrough() = false after WithStrikethrough(true)")
	}
}

func TestStyle_WithInverse(t *testing.T) {
	s := New().WithInverse(true)
	if !s.IsInverse() {
		t.Error("IsInverse() = false after WithInverse(true)")
	}
}

func TestStyle_WithForeground(t *testing.T) {
	c := RGB(255, 0, 0)
	s := New().WithForeground(c)
	if s.GetForeground() != c {
		t.Errorf("GetForeground() = %v, want %v", s.GetForeground(), c)
	}
}

func TestStyle_WithBackground(t *testing.T) {
	c := RGB(0, 0, 255)
	s := New().WithBackground(c)
	if s.GetBackground() != c {
		t.Errorf("GetBackground() = %v, want %v", s.GetBackground(), c)
	}
}

// ---------------------------------------------------------------------------
// Render — colour disabled
// ---------------------------------------------------------------------------

func TestStyle_Render_ModeNever_PlainText(t *testing.T) {
	forceColorMode(t, colorModeNever)
	s := New().WithForeground(Red).WithBold(true)
	got := s.Render("hello")
	if got != "hello" {
		t.Errorf("Render with NeverColorMode = %q, want %q", got, "hello")
	}
	if containsSGR(got) {
		t.Error("Render with NeverColorMode produced SGR sequences")
	}
}

func TestStyle_Render_ZeroStyle_PlainText(t *testing.T) {
	forceColorMode(t, colorModeAlways)
	got := New().Render("hello")
	if got != "hello" {
		t.Errorf("zero Style.Render = %q, want plain %q", got, "hello")
	}
}

// ---------------------------------------------------------------------------
// Render — empty string must not emit dangling reset
// ---------------------------------------------------------------------------

func TestStyle_Render_EmptyString_NoDanglingReset(t *testing.T) {
	forceColorMode(t, colorModeAlways)
	s := New().WithForeground(Red).WithBold(true)
	got := s.Render("")
	if got != "" {
		t.Errorf("Render(\"\") = %q, want empty string", got)
	}
}

// ---------------------------------------------------------------------------
// Render — SGR sequences present when colour is enabled
// ---------------------------------------------------------------------------

func TestStyle_Render_Bold(t *testing.T) {
	forceColorMode(t, colorModeAlways)
	got := New().WithBold(true).Render("hi")
	if !containsSGR(got) {
		t.Error("bold Render produced no SGR sequence")
	}
	if !strings.Contains(got, "1") {
		t.Errorf("bold Render does not contain SGR code 1, got %q", got)
	}
	if !strings.Contains(got, "\x1b[0m") {
		t.Errorf("bold Render missing reset suffix, got %q", got)
	}
}

func TestStyle_Render_Dim(t *testing.T) {
	forceColorMode(t, colorModeAlways)
	got := New().WithDim(true).Render("hi")
	if !strings.Contains(got, "2") {
		t.Errorf("dim Render does not contain SGR code 2, got %q", got)
	}
}

func TestStyle_Render_Italic(t *testing.T) {
	forceColorMode(t, colorModeAlways)
	got := New().WithItalic(true).Render("hi")
	if !strings.Contains(got, "3") {
		t.Errorf("italic Render does not contain SGR code 3, got %q", got)
	}
}

func TestStyle_Render_Underline(t *testing.T) {
	forceColorMode(t, colorModeAlways)
	got := New().WithUnderline(true).Render("hi")
	if !strings.Contains(got, "4") {
		t.Errorf("underline Render does not contain SGR code 4, got %q", got)
	}
}

func TestStyle_Render_Inverse(t *testing.T) {
	forceColorMode(t, colorModeAlways)
	got := New().WithInverse(true).Render("hi")
	if !strings.Contains(got, "7") {
		t.Errorf("inverse Render does not contain SGR code 7, got %q", got)
	}
}

func TestStyle_Render_Strikethrough(t *testing.T) {
	forceColorMode(t, colorModeAlways)
	got := New().WithStrikethrough(true).Render("hi")
	if !strings.Contains(got, "9") {
		t.Errorf("strikethrough Render does not contain SGR code 9, got %q", got)
	}
}

func TestStyle_Render_ForegroundRGB(t *testing.T) {
	forceColorMode(t, colorModeAlways)
	got := New().WithForeground(RGB(255, 128, 0)).Render("hi")
	// RGB fg: ESC[38;2;255;128;0m
	if !strings.Contains(got, "38;2;255;128;0") {
		t.Errorf("RGB fg Render = %q, expected 38;2;255;128;0", got)
	}
}

func TestStyle_Render_BackgroundRGB(t *testing.T) {
	forceColorMode(t, colorModeAlways)
	got := New().WithBackground(RGB(0, 64, 128)).Render("hi")
	// RGB bg: ESC[48;2;0;64;128m
	if !strings.Contains(got, "48;2;0;64;128") {
		t.Errorf("RGB bg Render = %q, expected 48;2;0;64;128", got)
	}
}

func TestStyle_Render_ForegroundANSI16(t *testing.T) {
	forceColorMode(t, colorModeAlways)
	got := New().WithForeground(ANSI16(1)).Render("hi")
	if !containsSGR(got) {
		t.Error("ANSI16 fg Render produced no SGR sequence")
	}
}

func TestStyle_Render_ForegroundANSI256(t *testing.T) {
	forceColorMode(t, colorModeAlways)
	got := New().WithForeground(ANSI256(200)).Render("hi")
	// 256-colour fg: ESC[38;5;200m
	if !strings.Contains(got, "38;5;200") {
		t.Errorf("ANSI256 fg Render = %q, expected 38;5;200", got)
	}
}

func TestStyle_Render_BackgroundANSI256(t *testing.T) {
	forceColorMode(t, colorModeAlways)
	got := New().WithBackground(ANSI256(100)).Render("hi")
	if !strings.Contains(got, "48;5;100") {
		t.Errorf("ANSI256 bg Render = %q, expected 48;5;100", got)
	}
}

func TestStyle_Render_MultipleAttributes(t *testing.T) {
	forceColorMode(t, colorModeAlways)
	s := New().
		WithBold(true).
		WithItalic(true).
		WithForeground(RGB(200, 100, 50))
	got := s.Render("hi")
	if !strings.Contains(got, "1") {
		t.Errorf("combined Render missing bold code, got %q", got)
	}
	if !strings.Contains(got, "3") {
		t.Errorf("combined Render missing italic code, got %q", got)
	}
	if !strings.Contains(got, "38;2;200;100;50") {
		t.Errorf("combined Render missing fg RGB, got %q", got)
	}
	// Only one leading ESC[ sequence.
	if strings.Count(got, "\x1b[") != 2 { // opening + reset
		t.Errorf("combined Render should have exactly 2 ESC[ sequences, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// Render — Strip round-trip
// ---------------------------------------------------------------------------

func TestStyle_Render_StripRoundTrip(t *testing.T) {
	forceColorMode(t, colorModeAlways)
	texts := []string{
		"hello",
		"hello, world!",
		"line1\nline2\nline3",
		"unicode: café",
	}
	styles := []Style{
		New().WithBold(true),
		New().WithForeground(RGB(255, 0, 128)).WithBackground(RGB(0, 0, 0)),
		New().WithItalic(true).WithUnderline(true),
		New().WithBold(true).WithForeground(ANSI16(3)).WithStrikethrough(true),
	}

	for _, text := range texts {
		for _, s := range styles {
			rendered := s.Render(text)
			stripped := Strip(rendered)
			if stripped != text {
				t.Errorf("Strip(Render(%q)) = %q, want original text", text, stripped)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// String
// ---------------------------------------------------------------------------

func TestStyle_String_Zero(t *testing.T) {
	got := New().String()
	if got != "Style{}" {
		t.Errorf("zero Style.String() = %q, want %q", got, "Style{}")
	}
}

func TestStyle_String_WithAttributes(t *testing.T) {
	s := New().WithBold(true).WithForeground(Red)
	got := s.String()
	if !strings.Contains(got, "bold") {
		t.Errorf("Style.String() missing 'bold', got %q", got)
	}
	if !strings.Contains(got, "fg=") {
		t.Errorf("Style.String() missing 'fg=', got %q", got)
	}
}

func TestStyle_String_AllAttributes(t *testing.T) {
	cases := []struct {
		name  string
		style Style
		want  string
	}{
		{"dim", New().WithDim(true), "dim"},
		{"italic", New().WithItalic(true), "italic"},
		{"underline", New().WithUnderline(true), "underline"},
		{"strikethrough", New().WithStrikethrough(true), "strikethrough"},
		{"inverse", New().WithInverse(true), "inverse"},
		{"background", New().WithBackground(Blue), "bg="},
		{"layout set", New().WithLayout(NewLayout().WithUniformPadding(1)), "layout=<set>"},
		{"border set", New().WithBorder(BorderNormal()), "border=<set>"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.style.String()
			if !strings.Contains(got, tc.want) {
				t.Errorf("Style.String() missing %q, got %q", tc.want, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Render — layout: MaxWidth truncation and JustifyContent alignment
// ---------------------------------------------------------------------------

func TestStyle_Render_Layout_MaxWidth_Truncates(t *testing.T) {
	forceColorMode(t, colorModeNever)

	s := New().WithLayout(NewLayout().WithMaxWidth(5))
	got := s.Render("Hello, world!")

	// Result must not exceed 5 visible columns (4 chars + ellipsis).
	stripped := Strip(got)
	if runeWidth(stripped) > 5 {
		t.Errorf("MaxWidth(5) Render visible width = %d > 5, got %q", runeWidth(stripped), stripped)
	}
	if !strings.Contains(stripped, "…") {
		t.Errorf("MaxWidth truncation missing ellipsis, got %q", stripped)
	}
}

func TestStyle_Render_Layout_MaxWidth_NoTruncation_WhenShort(t *testing.T) {
	forceColorMode(t, colorModeNever)

	s := New().WithLayout(NewLayout().WithMaxWidth(20))
	got := s.Render("Hi")

	if Strip(got) != "Hi" {
		t.Errorf("MaxWidth(20) on short text = %q, want %q", Strip(got), "Hi")
	}
	if strings.Contains(got, "…") {
		t.Errorf("MaxWidth(20) on short text produced unexpected ellipsis")
	}
}

func TestStyle_Render_Layout_JustifyEnd(t *testing.T) {
	forceColorMode(t, colorModeNever)

	s := New().WithLayout(
		NewLayout().WithMinWidth(10).WithJustifyContent(JustifyEnd),
	)
	got := s.Render("Hi")

	// "Hi" is 2 columns; minWidth=10 → 8 spaces prepended.
	if !strings.HasPrefix(got, "        ") {
		t.Errorf("JustifyEnd(minWidth=10) Render = %q, expected leading spaces", got)
	}
	if !strings.Contains(got, "Hi") {
		t.Errorf("JustifyEnd Render lost content, got %q", got)
	}
}

func TestStyle_Render_Layout_JustifyCenter(t *testing.T) {
	forceColorMode(t, colorModeNever)

	s := New().WithLayout(
		NewLayout().WithMinWidth(10).WithJustifyContent(JustifyCenter),
	)
	got := s.Render("Hi")

	// "Hi" is 2 columns; pad=8 → 4 spaces on each side.
	if !strings.HasPrefix(got, "    ") {
		t.Errorf("JustifyCenter(minWidth=10) Render = %q, expected leading spaces", got)
	}
	if !strings.HasSuffix(got, "    ") {
		t.Errorf("JustifyCenter(minWidth=10) Render = %q, expected trailing spaces", got)
	}
	if !strings.Contains(got, "Hi") {
		t.Errorf("JustifyCenter Render lost content, got %q", got)
	}
}

func TestStyle_Render_Layout_JustifyStart(t *testing.T) {
	forceColorMode(t, colorModeNever)

	s := New().WithLayout(
		NewLayout().WithMinWidth(10).WithJustifyContent(JustifyStart),
	)
	got := s.Render("Hi")

	// "Hi" is 2 columns; pad=8 → 8 spaces appended.
	if !strings.HasPrefix(got, "Hi") {
		t.Errorf("JustifyStart(minWidth=10) Render = %q, expected content first", got)
	}
	if !strings.HasSuffix(got, "        ") {
		t.Errorf("JustifyStart(minWidth=10) Render = %q, expected trailing spaces", got)
	}
}

func TestStyle_Render_Layout_Padding(t *testing.T) {
	forceColorMode(t, colorModeNever)

	s := New().WithLayout(NewLayout().WithUniformPadding(1))
	got := s.Render("X")

	lines := strings.Split(got, "\n")
	// uniform padding 1 → 1 blank line top, content line, 1 blank line bottom = 3 lines
	if len(lines) != 3 {
		t.Errorf("UniformPadding(1) Render line count = %d, want 3\noutput: %q", len(lines), got)
	}
	// Content line must have leading and trailing space.
	if !strings.HasPrefix(lines[1], " ") {
		t.Errorf("UniformPadding(1) content line missing leading space: %q", lines[1])
	}
	if !strings.HasSuffix(lines[1], " ") {
		t.Errorf("UniformPadding(1) content line missing trailing space: %q", lines[1])
	}
}

// ---------------------------------------------------------------------------
// Unset
// ---------------------------------------------------------------------------

func TestStyle_Unset_ClearsAll(t *testing.T) {
	s := New().
		WithBold(true).
		WithItalic(true).
		WithForeground(Red).
		WithBackground(Blue)
	cleared := s.Unset()
	if !reflect.DeepEqual(cleared, Style{}) {
		t.Error("Unset() did not return the zero Style")
	}
	// Original must be unchanged.
	if !s.IsBold() {
		t.Error("Unset() mutated the receiver")
	}
}

// ---------------------------------------------------------------------------
// Immutability
// ---------------------------------------------------------------------------

func TestStyle_Immutability(t *testing.T) {
	base := New().WithBold(true)
	_ = base.WithItalic(true)
	if base.IsItalic() {
		t.Error("chaining WithItalic mutated the base style")
	}
}
