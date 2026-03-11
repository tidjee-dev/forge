package cast

import "github.com/tidjee-dev/forge/ink"

// visibleWidth returns the number of terminal columns occupied by s after
// stripping any ANSI escape sequences. It delegates to ink.Strip for
// sequence removal, then counts rune column widths using the same logic ink
// uses internally.
func visibleWidth(s string) int {
	plain := ink.Strip(s)
	w := 0
	for _, r := range plain {
		w += runeColWidth(r)
	}
	return w
}

// runeColWidth returns the number of terminal columns occupied by a single
// rune. This mirrors ink's internal runeColumnWidth so that cast never
// diverges in its width calculations.
//
//   - 0 for C0/C1 controls, combining marks, and zero-width characters
//   - 2 for wide Unicode characters (CJK, fullwidth, common emoji)
//   - 1 for everything else
func runeColWidth(r rune) int {
	switch {
	case r < 0x20:
		return 0
	case r == 0x7F:
		return 0
	case r >= 0x80 && r < 0xA0:
		return 0
	case r >= 0x0300 && r <= 0x036F:
		return 0
	case r >= 0x200B && r <= 0x200F:
		return 0
	case r == 0xFEFF:
		return 0
	case r >= 0x1100 && r <= 0x115F:
		return 2
	case r >= 0x2E80 && r <= 0x303E:
		return 2
	case r >= 0x3040 && r <= 0x33FF:
		return 2
	case r >= 0x3400 && r <= 0x4DBF:
		return 2
	case r >= 0x4E00 && r <= 0x9FFF:
		return 2
	case r >= 0xA000 && r <= 0xA4CF:
		return 2
	case r >= 0xAC00 && r <= 0xD7AF:
		return 2
	case r >= 0xF900 && r <= 0xFAFF:
		return 2
	case r >= 0xFE10 && r <= 0xFE1F:
		return 2
	case r >= 0xFE30 && r <= 0xFE6F:
		return 2
	case r >= 0xFF01 && r <= 0xFF60:
		return 2
	case r >= 0xFFE0 && r <= 0xFFE6:
		return 2
	case r >= 0x1B000 && r <= 0x1B0FF:
		return 2
	case r >= 0x1F004 && r <= 0x1F0CF:
		return 2
	case r >= 0x1F300 && r <= 0x1F9FF:
		return 2
	case r >= 0x20000 && r <= 0x2FFFD:
		return 2
	case r >= 0x30000 && r <= 0x3FFFD:
		return 2
	default:
		return 1
	}
}

// padRight appends spaces to s until its visible width reaches width.
// If s is already at or beyond width, s is returned unchanged.
func padRight(s string, width int) string {
	w := visibleWidth(s)
	if w >= width {
		return s
	}
	for i := 0; i < width-w; i++ {
		s += " "
	}
	return s
}

// padLeft prepends spaces to s until its visible width reaches width.
// If s is already at or beyond width, s is returned unchanged.
func padLeft(s string, width int) string {
	w := visibleWidth(s)
	if w >= width {
		return s
	}
	pad := ""
	for i := 0; i < width-w; i++ {
		pad += " "
	}
	return pad + s
}

// centerPad centres s within a field of the given width by distributing
// padding spaces evenly on both sides. If s is already at or beyond width,
// s is returned unchanged.
func centerPad(s string, width int) string {
	w := visibleWidth(s)
	if w >= width {
		return s
	}
	total := width - w
	left := total / 2
	right := total - left
	lpad := ""
	rpad := ""
	for i := 0; i < left; i++ {
		lpad += " "
	}
	for i := 0; i < right; i++ {
		rpad += " "
	}
	return lpad + s + rpad
}

// repeatStr repeats the string s n times. Unlike strings.Repeat it accepts
// n <= 0 gracefully and returns "".
func repeatStr(s string, n int) string {
	if n <= 0 || s == "" {
		return ""
	}
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}

// isStyleSet reports whether s has any attribute configured — i.e. whether it
// differs from the zero-value Style. Because ink.Style contains a Layout which
// holds slices, direct == comparison is not allowed in Go. Instead we check the
// observable attributes that cast cares about.
//
// This is intentionally conservative: it returns true as soon as any common
// attribute is non-zero. The result is used only to decide whether to apply a
// style at all, not to compare two styles for equality.
func isStyleSet(s ink.Style) bool {
	return !s.GetForeground().IsZeroColor() ||
		!s.GetBackground().IsZeroColor() ||
		s.IsBold() ||
		s.IsItalic() ||
		s.IsUnderline() ||
		s.IsDim() ||
		s.IsStrikethrough() ||
		s.IsInverse() ||
		!s.GetBorderStyle().IsZero()
}
