package ink

import (
	"strings"
	"sync"
	"testing"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// enableColor forces AlwaysColorMode for the duration of t and restores state
// on cleanup. Reused across border tests that check coloured glyphs.
func enableColorForBorder(t *testing.T) {
	t.Helper()
	old := globalColorMode.Load()
	t.Cleanup(func() {
		globalColorMode.Store(old)
		envOnce = sync.Once{}
		envDisabled = false
		envForced = false
	})
	SetGlobalColorMode(colorModeAlways)
}

// renderBorder is a thin wrapper that calls Style.Render with colour always
// enabled, so border tests are not affected by the test runner's TTY state.
func renderBorder(t *testing.T, s Style, text string) string {
	t.Helper()
	enableColorForBorder(t)
	return s.Render(text)
}

// ---------------------------------------------------------------------------
// BorderStyle.IsZero
// ---------------------------------------------------------------------------

func TestBorderStyle_IsZero_ZeroValue(t *testing.T) {
	var bs BorderStyle
	if !bs.IsZero() {
		t.Error("zero-value BorderStyle.IsZero() = false, want true")
	}
}

func TestBorderStyle_IsZero_NoBorder(t *testing.T) {
	if !NoBorder().IsZero() {
		t.Error("NoBorder().IsZero() = false, want true")
	}
}

func TestBorderStyle_IsZero_FalseWhenAnyGlyphSet(t *testing.T) {
	presets := []struct {
		name string
		bs   BorderStyle
	}{
		{"BorderNormal", BorderNormal()},
		{"BorderRounded", BorderRounded()},
		{"BorderThick", BorderThick()},
		{"BorderDouble", BorderDouble()},
		{"BorderASCII", BorderASCII()},
		{"BorderDashed", BorderDashed()},
		{"BorderBlock", BorderBlock()},
		{"BorderInnerHalfBlock", BorderInnerHalfBlock()},
	}
	for _, p := range presets {
		t.Run(p.name, func(t *testing.T) {
			if p.bs.IsZero() {
				t.Errorf("%s.IsZero() = true, want false", p.name)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Preset glyph correctness
// ---------------------------------------------------------------------------

func TestBorderNormal_Glyphs(t *testing.T) {
	bs := BorderNormal()
	assertGlyph(t, "TopLeft", bs.TopLeft, "┌")
	assertGlyph(t, "TopRight", bs.TopRight, "┐")
	assertGlyph(t, "BottomLeft", bs.BottomLeft, "└")
	assertGlyph(t, "BottomRight", bs.BottomRight, "┘")
	assertGlyph(t, "Top", bs.Top, "─")
	assertGlyph(t, "Bottom", bs.Bottom, "─")
	assertGlyph(t, "Left", bs.Left, "│")
	assertGlyph(t, "Right", bs.Right, "│")
}

func TestBorderRounded_Glyphs(t *testing.T) {
	bs := BorderRounded()
	assertGlyph(t, "TopLeft", bs.TopLeft, "╭")
	assertGlyph(t, "TopRight", bs.TopRight, "╮")
	assertGlyph(t, "BottomLeft", bs.BottomLeft, "╰")
	assertGlyph(t, "BottomRight", bs.BottomRight, "╯")
	assertGlyph(t, "Top", bs.Top, "─")
	assertGlyph(t, "Bottom", bs.Bottom, "─")
	assertGlyph(t, "Left", bs.Left, "│")
	assertGlyph(t, "Right", bs.Right, "│")
}

func TestBorderThick_Glyphs(t *testing.T) {
	bs := BorderThick()
	assertGlyph(t, "TopLeft", bs.TopLeft, "┏")
	assertGlyph(t, "TopRight", bs.TopRight, "┓")
	assertGlyph(t, "BottomLeft", bs.BottomLeft, "┗")
	assertGlyph(t, "BottomRight", bs.BottomRight, "┛")
	assertGlyph(t, "Top", bs.Top, "━")
	assertGlyph(t, "Bottom", bs.Bottom, "━")
	assertGlyph(t, "Left", bs.Left, "┃")
	assertGlyph(t, "Right", bs.Right, "┃")
}

func TestBorderDouble_Glyphs(t *testing.T) {
	bs := BorderDouble()
	assertGlyph(t, "TopLeft", bs.TopLeft, "╔")
	assertGlyph(t, "TopRight", bs.TopRight, "╗")
	assertGlyph(t, "BottomLeft", bs.BottomLeft, "╚")
	assertGlyph(t, "BottomRight", bs.BottomRight, "╝")
	assertGlyph(t, "Top", bs.Top, "═")
	assertGlyph(t, "Bottom", bs.Bottom, "═")
	assertGlyph(t, "Left", bs.Left, "║")
	assertGlyph(t, "Right", bs.Right, "║")
}

func TestBorderASCII_Glyphs(t *testing.T) {
	bs := BorderASCII()
	assertGlyph(t, "TopLeft", bs.TopLeft, "+")
	assertGlyph(t, "TopRight", bs.TopRight, "+")
	assertGlyph(t, "BottomLeft", bs.BottomLeft, "+")
	assertGlyph(t, "BottomRight", bs.BottomRight, "+")
	assertGlyph(t, "Top", bs.Top, "-")
	assertGlyph(t, "Bottom", bs.Bottom, "-")
	assertGlyph(t, "Left", bs.Left, "|")
	assertGlyph(t, "Right", bs.Right, "|")
}

func assertGlyph(t *testing.T, field, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("BorderStyle.%s = %q, want %q", field, got, want)
	}
}

func TestBorderDashed_Glyphs(t *testing.T) {
	bs := BorderDashed()
	assertGlyph(t, "TopLeft", bs.TopLeft, "┌")
	assertGlyph(t, "TopRight", bs.TopRight, "┐")
	assertGlyph(t, "BottomLeft", bs.BottomLeft, "└")
	assertGlyph(t, "BottomRight", bs.BottomRight, "┘")
	assertGlyph(t, "Top", bs.Top, "┄")
	assertGlyph(t, "Bottom", bs.Bottom, "┄")
	assertGlyph(t, "Left", bs.Left, "┆")
	assertGlyph(t, "Right", bs.Right, "┆")
}

func TestBorderBlock_Glyphs(t *testing.T) {
	bs := BorderBlock()
	for _, field := range []struct {
		name  string
		value string
	}{
		{"TopLeft", bs.TopLeft},
		{"TopRight", bs.TopRight},
		{"BottomLeft", bs.BottomLeft},
		{"BottomRight", bs.BottomRight},
		{"Top", bs.Top},
		{"Bottom", bs.Bottom},
		{"Left", bs.Left},
		{"Right", bs.Right},
	} {
		assertGlyph(t, field.name, field.value, "█")
	}
}

func TestBorderHidden_Glyphs(t *testing.T) {
	bs := BorderHidden()
	for _, field := range []struct {
		name  string
		value string
	}{
		{"TopLeft", bs.TopLeft},
		{"TopRight", bs.TopRight},
		{"BottomLeft", bs.BottomLeft},
		{"BottomRight", bs.BottomRight},
		{"Top", bs.Top},
		{"Bottom", bs.Bottom},
		{"Left", bs.Left},
		{"Right", bs.Right},
	} {
		assertGlyph(t, field.name, field.value, " ")
	}
}

func TestBorderInnerHalfBlock_Glyphs(t *testing.T) {
	bs := BorderInnerHalfBlock()
	assertGlyph(t, "TopLeft", bs.TopLeft, "▄")
	assertGlyph(t, "TopRight", bs.TopRight, "▄")
	assertGlyph(t, "BottomLeft", bs.BottomLeft, "▀")
	assertGlyph(t, "BottomRight", bs.BottomRight, "▀")
	assertGlyph(t, "Top", bs.Top, "▄")
	assertGlyph(t, "Bottom", bs.Bottom, "▀")
	assertGlyph(t, "Left", bs.Left, "▌")
	assertGlyph(t, "Right", bs.Right, "▐")
}

// ---------------------------------------------------------------------------
// Style border setters
// ---------------------------------------------------------------------------

func TestStyle_WithBorder_SetsBorderSideAll(t *testing.T) {
	s := New().WithBorder(BorderNormal())
	if s.GetBorderSide() != BorderSideAll {
		t.Errorf("WithBorder sets sides = %v, want BorderSideAll", s.GetBorderSide())
	}
	if s.GetBorderStyle() != BorderNormal() {
		t.Error("WithBorder did not store the BorderStyle")
	}
}

func TestStyle_WithBorderStyle_SetsSideAllWhenNoneWasPreviouslySet(t *testing.T) {
	s := New().WithBorderStyle(BorderRounded())
	if s.GetBorderSide() != BorderSideAll {
		t.Errorf("WithBorderStyle (fresh) sides = %v, want BorderSideAll", s.GetBorderSide())
	}
}

func TestStyle_WithBorderStyle_PreservesExistingSides(t *testing.T) {
	s := New().
		WithBorderSide(BorderSideTop | BorderSideBottom).
		WithBorderStyle(BorderThick())
	// Sides were already set; WithBorderStyle must not override them.
	if s.GetBorderSide() != BorderSideTop|BorderSideBottom {
		t.Errorf("WithBorderStyle changed sides to %v, want Top|Bottom", s.GetBorderSide())
	}
}

func TestStyle_WithBorderColor(t *testing.T) {
	c := RGB(255, 0, 0)
	s := New().WithBorderColor(c)
	if s.GetBorderColor() != c {
		t.Errorf("GetBorderColor() = %v, want %v", s.GetBorderColor(), c)
	}
}

func TestStyle_WithBorderSide_Subsets(t *testing.T) {
	cases := []struct {
		name  string
		sides BorderSide
	}{
		{"None", BorderSideNone},
		{"Top", BorderSideTop},
		{"Right", BorderSideRight},
		{"Bottom", BorderSideBottom},
		{"Left", BorderSideLeft},
		{"All", BorderSideAll},
		{"Vertical", BorderSideVertical},
		{"Horizontal", BorderSideHorizontal},
		{"Top+Bottom", BorderSideTop | BorderSideBottom},
		{"Left+Right", BorderSideLeft | BorderSideRight},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := New().WithBorderSide(tc.sides)
			if s.GetBorderSide() != tc.sides {
				t.Errorf("GetBorderSide() = %v, want %v", s.GetBorderSide(), tc.sides)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Render — border glyphs appear in output
// ---------------------------------------------------------------------------

func TestStyle_Render_Border_AllSides(t *testing.T) {
	enableColorForBorder(t)

	s := New().WithBorder(BorderNormal())
	out := s.Render("hello")

	// All four corner and edge glyphs must be present.
	for _, glyph := range []string{"┌", "┐", "└", "┘", "─", "│"} {
		if !strings.Contains(out, glyph) {
			t.Errorf("Render with BorderNormal missing glyph %q\noutput: %q", glyph, out)
		}
	}
	// Content must still be present.
	if !strings.Contains(Strip(out), "hello") {
		t.Errorf("Render with border lost content, stripped output: %q", Strip(out))
	}
}

func TestStyle_Render_Border_Rounded(t *testing.T) {
	enableColorForBorder(t)

	s := New().WithBorder(BorderRounded())
	out := s.Render("hi")

	for _, glyph := range []string{"╭", "╮", "╰", "╯"} {
		if !strings.Contains(out, glyph) {
			t.Errorf("Render with BorderRounded missing glyph %q\noutput: %q", glyph, out)
		}
	}
}

func TestStyle_Render_Border_Thick(t *testing.T) {
	enableColorForBorder(t)

	s := New().WithBorder(BorderThick())
	out := s.Render("hi")

	for _, glyph := range []string{"┏", "┓", "┗", "┛", "━", "┃"} {
		if !strings.Contains(out, glyph) {
			t.Errorf("Render with BorderThick missing glyph %q\noutput: %q", glyph, out)
		}
	}
}

func TestStyle_Render_Border_Double(t *testing.T) {
	enableColorForBorder(t)

	s := New().WithBorder(BorderDouble())
	out := s.Render("hi")

	for _, glyph := range []string{"╔", "╗", "╚", "╝", "═", "║"} {
		if !strings.Contains(out, glyph) {
			t.Errorf("Render with BorderDouble missing glyph %q\noutput: %q", glyph, out)
		}
	}
}

// ---------------------------------------------------------------------------
// Render — BorderSide subsets
// ---------------------------------------------------------------------------

func TestStyle_Render_Border_TopOnly(t *testing.T) {
	enableColorForBorder(t)

	s := New().
		WithBorderStyle(BorderNormal()).
		WithBorderSide(BorderSideTop)
	out := s.Render("hello")

	if !strings.Contains(out, "─") {
		t.Errorf("top-only border missing top glyph, output: %q", out)
	}
	// No corner glyphs from bottom row.
	if strings.Contains(out, "└") || strings.Contains(out, "┘") {
		t.Errorf("top-only border contains bottom corners, output: %q", out)
	}
	// No left/right side glyphs.
	if strings.Contains(out, "│") {
		t.Errorf("top-only border contains side glyphs, output: %q", out)
	}
}

func TestStyle_Render_Border_BottomOnly(t *testing.T) {
	enableColorForBorder(t)

	s := New().
		WithBorderStyle(BorderNormal()).
		WithBorderSide(BorderSideBottom)
	out := s.Render("hello")

	if !strings.Contains(out, "─") {
		t.Errorf("bottom-only border missing bottom glyph, output: %q", out)
	}
	if strings.Contains(out, "┌") || strings.Contains(out, "┐") {
		t.Errorf("bottom-only border contains top corners, output: %q", out)
	}
}

func TestStyle_Render_Border_LeftOnly(t *testing.T) {
	enableColorForBorder(t)

	s := New().
		WithBorderStyle(BorderNormal()).
		WithBorderSide(BorderSideLeft)
	out := s.Render("hello")

	if !strings.Contains(out, "│") {
		t.Errorf("left-only border missing side glyph, output: %q", out)
	}
	// No top/bottom bar glyphs.
	if strings.Contains(out, "─") {
		t.Errorf("left-only border contains horizontal glyphs, output: %q", out)
	}
}

func TestStyle_Render_Border_RightOnly(t *testing.T) {
	enableColorForBorder(t)

	s := New().
		WithBorderStyle(BorderNormal()).
		WithBorderSide(BorderSideRight)
	out := s.Render("hello")

	if !strings.Contains(out, "│") {
		t.Errorf("right-only border missing side glyph, output: %q", out)
	}
}

func TestStyle_Render_Border_Vertical(t *testing.T) {
	enableColorForBorder(t)

	s := New().
		WithBorderStyle(BorderNormal()).
		WithBorderSide(BorderSideVertical)
	out := s.Render("hello")

	// Top and bottom bars present.
	if !strings.Contains(out, "─") {
		t.Errorf("vertical border missing horizontal glyphs, output: %q", out)
	}
	// No side glyphs.
	if strings.Contains(out, "│") {
		t.Errorf("vertical border contains side glyphs, output: %q", out)
	}
}

func TestStyle_Render_Border_Horizontal(t *testing.T) {
	enableColorForBorder(t)

	s := New().
		WithBorderStyle(BorderNormal()).
		WithBorderSide(BorderSideHorizontal)
	out := s.Render("hello")

	// Side glyphs present.
	if !strings.Contains(out, "│") {
		t.Errorf("horizontal border missing side glyphs, output: %q", out)
	}
	// No top/bottom bars.
	if strings.Contains(out, "─") {
		t.Errorf("horizontal border contains horizontal glyphs, output: %q", out)
	}
}

func TestStyle_Render_Border_None_NoGlyphs(t *testing.T) {
	enableColorForBorder(t)

	s := New().
		WithBorderStyle(BorderNormal()).
		WithBorderSide(BorderSideNone)
	out := s.Render("hello")

	for _, glyph := range []string{"┌", "┐", "└", "┘", "─", "│"} {
		if strings.Contains(out, glyph) {
			t.Errorf("BorderSideNone still renders glyph %q, output: %q", glyph, out)
		}
	}
}

// ---------------------------------------------------------------------------
// Render — border colour applies to glyphs only
// ---------------------------------------------------------------------------

func TestStyle_Render_BorderColor_AppliesToGlyphs(t *testing.T) {
	enableColorForBorder(t)

	borderCol := RGB(255, 0, 0) // red border
	s := New().
		WithBorder(BorderNormal()).
		WithBorderColor(borderCol)
	out := s.Render("hello")

	// Border glyph must be present.
	if !strings.Contains(out, "─") {
		t.Errorf("bordered output missing glyph, output: %q", out)
	}
	// Red RGB fg sequence (38;2;255;0;0) must appear — it colours the glyphs.
	if !strings.Contains(out, "38;2;255;0;0") {
		t.Errorf("border colour sequence missing from output: %q", out)
	}
	// Content must survive Strip intact.
	stripped := Strip(out)
	if !strings.Contains(stripped, "hello") {
		t.Errorf("strip of bordered+coloured output lost content: %q", stripped)
	}
}

func TestStyle_Render_BorderColor_ZeroColor_NoExtraSequence(t *testing.T) {
	enableColorForBorder(t)

	// No border colour set — glyphs should appear without colour wrapping.
	s := New().WithBorder(BorderNormal())
	out := s.Render("hello")

	// Content glyph present.
	if !strings.Contains(out, "─") {
		t.Errorf("border output missing glyph, output: %q", out)
	}
	// The outer style has no fg/bg and no border colour → zero SGR sequences
	// should appear in this specific test (zero Style apart from border).
	stripped := Strip(out)
	if !strings.Contains(stripped, "hello") {
		t.Errorf("stripped border output lost content: %q", stripped)
	}
}

// ---------------------------------------------------------------------------
// Render — NoBorder / zero BorderStyle skips border entirely
// ---------------------------------------------------------------------------

func TestStyle_Render_NoBorder_NoGlyphs(t *testing.T) {
	enableColorForBorder(t)

	s := New().WithBorder(NoBorder())
	out := s.Render("hello")

	for _, glyph := range []string{"┌", "┐", "└", "┘", "─", "│"} {
		if strings.Contains(out, glyph) {
			t.Errorf("NoBorder still rendered glyph %q, output: %q", glyph, out)
		}
	}
	if Strip(out) != "hello" {
		t.Errorf("NoBorder output != plain content, got %q", Strip(out))
	}
}

// ---------------------------------------------------------------------------
// Render — content still recoverable after Strip
// ---------------------------------------------------------------------------

func TestStyle_Render_Border_StripRoundTrip(t *testing.T) {
	enableColorForBorder(t)

	texts := []string{
		"hello",
		"longer content line",
		"line1\nline2\nline3",
	}
	presets := []BorderStyle{
		BorderNormal(),
		BorderRounded(),
		BorderThick(),
		BorderDouble(),
		BorderASCII(),
		BorderDashed(),
		BorderBlock(),
		BorderHidden(),
		BorderInnerHalfBlock(),
	}

	for _, text := range texts {
		for _, bs := range presets {
			s := New().WithBorder(bs)
			out := s.Render(text)
			stripped := Strip(out)
			// Every line of the original text must appear verbatim somewhere
			// in the stripped output (surrounded by border side glyphs).
			for _, line := range strings.Split(text, "\n") {
				if !strings.Contains(stripped, line) {
					t.Errorf("Strip(Render(%q, border)) missing line %q\nstripped: %q",
						text, line, stripped)
				}
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Render — remaining presets produce output with content intact
// ---------------------------------------------------------------------------

func TestStyle_Render_Border_Dashed(t *testing.T) {
	enableColorForBorder(t)

	s := New().WithBorder(BorderDashed())
	out := s.Render("hi")

	for _, glyph := range []string{"┌", "┐", "└", "┘", "┄", "┆"} {
		if !strings.Contains(out, glyph) {
			t.Errorf("Render with BorderDashed missing glyph %q\noutput: %q", glyph, out)
		}
	}
	if !strings.Contains(Strip(out), "hi") {
		t.Errorf("Render with BorderDashed lost content, stripped: %q", Strip(out))
	}
}

func TestStyle_Render_Border_Block(t *testing.T) {
	enableColorForBorder(t)

	s := New().WithBorder(BorderBlock())
	out := s.Render("hi")

	if !strings.Contains(out, "█") {
		t.Errorf("Render with BorderBlock missing glyph █\noutput: %q", out)
	}
	if !strings.Contains(Strip(out), "hi") {
		t.Errorf("Render with BorderBlock lost content, stripped: %q", Strip(out))
	}
}

func TestStyle_Render_Border_Hidden(t *testing.T) {
	enableColorForBorder(t)

	s := New().WithBorder(BorderHidden())
	out := s.Render("hi")

	// Hidden border uses spaces — content must still be present.
	if !strings.Contains(Strip(out), "hi") {
		t.Errorf("Render with BorderHidden lost content, stripped: %q", Strip(out))
	}
	// No box-drawing characters should appear.
	for _, glyph := range []string{"─", "│", "┌", "┐", "└", "┘"} {
		if strings.Contains(out, glyph) {
			t.Errorf("Render with BorderHidden contains unexpected glyph %q", glyph)
		}
	}
}

func TestStyle_Render_Border_InnerHalfBlock(t *testing.T) {
	enableColorForBorder(t)

	s := New().WithBorder(BorderInnerHalfBlock())
	out := s.Render("hi")

	for _, glyph := range []string{"▄", "▀", "▌", "▐"} {
		if !strings.Contains(out, glyph) {
			t.Errorf("Render with BorderInnerHalfBlock missing glyph %q\noutput: %q", glyph, out)
		}
	}
	if !strings.Contains(Strip(out), "hi") {
		t.Errorf("Render with BorderInnerHalfBlock lost content, stripped: %q", Strip(out))
	}
}

// ---------------------------------------------------------------------------
// Immutability
// ---------------------------------------------------------------------------

func TestStyle_Border_Immutability(t *testing.T) {
	base := New()
	_ = base.WithBorder(BorderNormal())
	if !base.GetBorderStyle().IsZero() {
		t.Error("WithBorder mutated the base style's border")
	}
	if base.GetBorderSide() != BorderSideNone {
		t.Error("WithBorder mutated the base style's border sides")
	}
}
