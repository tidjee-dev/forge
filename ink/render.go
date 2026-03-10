package ink

import (
	"os"
	"strings"
)

// ColorModeEnabled reports whether ANSI colour output should be emitted to
// stdout. It respects the global colour mode set via [SetGlobalColorMode]:
//   - [AlwaysColorMode] → always true
//   - [NeverColorMode]  → always false
//   - [AutoColorMode]   → auto-detected from environment and TTY state
func ColorModeEnabled() bool {
	return isColorEnabled(os.Stdout)
}

// ---------------------------------------------------------------------------
// SGR helpers
// ---------------------------------------------------------------------------

const (
	esc   = "\x1b["
	reset = esc + "0m"
)

// sgr builds a single SGR escape sequence from the given numeric parameters,
// e.g. sgr("1", "38", "2", "255", "0", "0") → "\x1b[1;38;2;255;0;0m".
func sgr(params ...string) string {
	return esc + strings.Join(params, ";") + "m"
}

// ---------------------------------------------------------------------------
// Style.Render
// ---------------------------------------------------------------------------

// Render applies the style's SGR attributes to s and returns the decorated
// string. When colour output is disabled (see [ColorModeEnabled]) or the
// style has no attributes set, s is returned unchanged.
//
// Rules:
//   - An empty s never produces a dangling SGR reset.
//   - A single ESC[0m reset is appended after the content when at least one
//     attribute is active.
//   - All active parameters are combined into one leading SGR sequence.
func (s Style) Render(text string) string {
	if text == "" {
		return text
	}

	colorEnabled := isColorEnabled(os.Stdout)

	var params []string

	if colorEnabled {
		// Text attributes
		if s.bold {
			params = append(params, "1")
		}
		if s.dim {
			params = append(params, "2")
		}
		if s.italic {
			params = append(params, "3")
		}
		if s.underline {
			params = append(params, "4")
		}
		if s.inverse {
			params = append(params, "7")
		}
		if s.strikethrough {
			params = append(params, "9")
		}

		// Foreground colour
		if !s.fgColor.IsZeroColor() {
			params = append(params, s.fgColor.fgParams()...)
		}

		// Background colour
		if !s.bgColor.IsZeroColor() {
			params = append(params, s.bgColor.bgParams()...)
		}
	}

	// Apply layout (padding, alignment, width constraints) to the raw text
	// before wrapping in SGR so that the padding spaces are also styled.
	styled := applyLayout(text, s.layout)

	// Apply border around the (already laid-out) content block.
	if !s.borderStyle.IsZero() && s.borderSides != BorderSideNone {
		styled = applyBorder(styled, s.borderStyle, s.borderSides, s.borderColor, colorEnabled)
	}

	if len(params) == 0 {
		return styled
	}

	return sgr(params...) + styled + reset
}

// String returns a human-readable description of the style, intended for
// debugging and test output.
func (s Style) String() string {
	var parts []string

	if !s.fgColor.IsZeroColor() {
		parts = append(parts, "fg="+s.fgColor.String())
	}
	if !s.bgColor.IsZeroColor() {
		parts = append(parts, "bg="+s.bgColor.String())
	}
	if s.bold {
		parts = append(parts, "bold")
	}
	if s.dim {
		parts = append(parts, "dim")
	}
	if s.italic {
		parts = append(parts, "italic")
	}
	if s.underline {
		parts = append(parts, "underline")
	}
	if s.strikethrough {
		parts = append(parts, "strikethrough")
	}
	if s.inverse {
		parts = append(parts, "inverse")
	}
	if !s.layout.IsZero() {
		parts = append(parts, "layout=<set>")
	}
	if !s.borderStyle.IsZero() {
		parts = append(parts, "border=<set>")
	}

	if len(parts) == 0 {
		return "Style{}"
	}
	return "Style{" + strings.Join(parts, " ") + "}"
}

// Unset returns the zero Style, clearing all attributes.
func (s Style) Unset() Style {
	return Style{}
}

// ---------------------------------------------------------------------------
// Layout application
// ---------------------------------------------------------------------------

// applyLayout applies padding, alignment, and width constraints from l to
// the multi-line text block. Each step is applied in order:
//  1. MaxWidth truncation (per line, with ellipsis)
//  2. MinWidth / alignment padding (per line)
//  3. Vertical (top/bottom) padding lines
//  4. Left/right padding columns
func applyLayout(text string, l Layout) string {
	if l.IsZero() {
		return text
	}

	lines := strings.Split(text, "\n")

	// 1. MaxWidth — truncate lines that exceed the limit.
	if l.maxWidth > 0 {
		for i, line := range lines {
			lines[i] = truncateToWidth(line, l.maxWidth)
		}
	}

	// 2. MinWidth / horizontal alignment.
	if l.minWidth > 0 {
		for i, line := range lines {
			w := runeWidth(line)
			if w < l.minWidth {
				pad := l.minWidth - w
				lines[i] = alignLine(line, pad, l.justifyContent)
			}
		}
	}

	// 3. Left/right padding — prepend/append spaces to every content line.
	if l.padding.Left > 0 || l.padding.Right > 0 {
		leftPad := strings.Repeat(" ", l.padding.Left)
		rightPad := strings.Repeat(" ", l.padding.Right)
		for i, line := range lines {
			lines[i] = leftPad + line + rightPad
		}
	}

	// 4. Top/bottom padding — insert blank lines.
	var out []string
	blankLine := strings.Repeat(" ", l.padding.Left+runeWidth(lines[0])+l.padding.Right)

	for i := 0; i < l.padding.Top; i++ {
		out = append(out, blankLine)
	}
	out = append(out, lines...)
	for i := 0; i < l.padding.Bottom; i++ {
		out = append(out, blankLine)
	}

	return strings.Join(out, "\n")
}

// alignLine distributes pad columns according to justification.
// JustifyCenter → split evenly; JustifyEnd → prepend all; default → append all.
func alignLine(line string, pad int, justify JustifyContent) string {
	switch justify {
	case JustifyCenter:
		left := pad / 2
		right := pad - left
		return strings.Repeat(" ", left) + line + strings.Repeat(" ", right)
	case JustifyEnd:
		return strings.Repeat(" ", pad) + line
	default: // JustifyStart
		return line + strings.Repeat(" ", pad)
	}
}

// ---------------------------------------------------------------------------
// Border application
// ---------------------------------------------------------------------------

// applyBorder draws border glyphs around the content block. Only the sides
// selected in sides are drawn; absent sides leave the content edges bare.
// borderColor sets the foreground color of glyphs; pass a zero Color for no
// color.
func applyBorder(text string, bs BorderStyle, sides BorderSide, borderColor Color, colorEnabled bool) string {
	lines := strings.Split(text, "\n")

	// Determine inner width: the widest content line.
	innerW := 0
	for _, l := range lines {
		if w := runeWidth(l); w > innerW {
			innerW = w
		}
	}

	hasTop := sides&BorderSideTop != 0
	hasBottom := sides&BorderSideBottom != 0
	hasLeft := sides&BorderSideLeft != 0
	hasRight := sides&BorderSideRight != 0

	// Colour wrapper for border glyphs.
	colorGlyph := func(g string) string {
		if !colorEnabled || borderColor.IsZeroColor() {
			return g
		}
		return sgr(borderColor.fgParams()...) + g + reset
	}

	totalW := innerW
	if hasLeft {
		totalW++
	}
	if hasRight {
		totalW++
	}

	var sb strings.Builder

	// Top border row.
	if hasTop {
		if hasLeft {
			sb.WriteString(colorGlyph(bs.TopLeft))
		}
		topBar := strings.Repeat(bs.Top, innerW)
		sb.WriteString(colorGlyph(topBar))
		if hasRight {
			sb.WriteString(colorGlyph(bs.TopRight))
		}
		sb.WriteByte('\n')
	}
	_ = totalW

	// Content rows.
	for _, line := range lines {
		if hasLeft {
			sb.WriteString(colorGlyph(bs.Left))
		}
		// Pad line to innerW so right border aligns.
		w := runeWidth(line)
		sb.WriteString(line)
		if w < innerW {
			sb.WriteString(strings.Repeat(" ", innerW-w))
		}
		if hasRight {
			sb.WriteString(colorGlyph(bs.Right))
		}
		sb.WriteByte('\n')
	}

	// Bottom border row.
	if hasBottom {
		if hasLeft {
			sb.WriteString(colorGlyph(bs.BottomLeft))
		}
		bottomBar := strings.Repeat(bs.Bottom, innerW)
		sb.WriteString(colorGlyph(bottomBar))
		if hasRight {
			sb.WriteString(colorGlyph(bs.BottomRight))
		}
		sb.WriteByte('\n')
	}

	// Trim the trailing newline that the last WriteByte('\n') added.
	result := sb.String()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result
}
