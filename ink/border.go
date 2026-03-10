package ink

// ---------------------------------------------------------------------------
// BorderSide — bitmask for selective border edges
// ---------------------------------------------------------------------------

// BorderSide is a bitmask that selects which edges of a border are drawn.
//
//	ink.New().
//	    WithBorderStyle(ink.BorderRounded()).
//	    WithBorderSides(ink.BorderSideTop | ink.BorderSideBottom)
type BorderSide uint8

const (
	// BorderSideNone disables every edge (no border drawn).
	BorderSideNone BorderSide = 0

	BorderSideTop    BorderSide = 1 << iota // Draw the top edge of the border.
	BorderSideRight                         // Draw the right edge of the border.
	BorderSideBottom                        // Draw the bottom edge of the border.
	BorderSideLeft                          // Draw the left edge of the border.

	// BorderSideAll enables all four edges. This is the default when a border
	// style is set via [Style.WithBorder] or [Style.WithBorderStyle].
	BorderSideAll = BorderSideLeft | BorderSideTop | BorderSideRight | BorderSideBottom

	// BorderSideVertical draws only the top and bottom edges of the border are drawn.
	BorderSideVertical = BorderSideTop | BorderSideBottom

	// BorderSideHorizontal draws only the left and right edges of the border are drawn.
	BorderSideHorizontal = BorderSideLeft | BorderSideRight
)

// ---------------------------------------------------------------------------
// BorderStyle — glyph set
// ---------------------------------------------------------------------------

// BorderStyle holds the nine box-drawing characters used to render a border:
// the four corner glyphs, the four edge glyphs (top, right, bottom, left),
// and the junction character that appears between repeated horizontal runs.
//
// All fields are single-rune strings. Multi-rune values are allowed but each
// character should occupy exactly one terminal column so that alignment is
// preserved.
//
// Build one with a preset ([BorderNormal], [BorderRounded], …) or fill the
// struct directly for a custom design.
type BorderStyle struct {
	// Corner glyphs
	TopLeft     string // e.g. "┌"
	TopRight    string // e.g. "┐"
	BottomLeft  string // e.g. "└"
	BottomRight string // e.g. "┘"

	// Edge glyphs
	Top    string // e.g. "─" (repeated across the top)
	Bottom string // e.g. "─" (repeated across the bottom)
	Left   string // e.g. "│" (appears on the left of every content row)
	Right  string // e.g. "│" (appears on the right of every content row)
}

// IsZero reports whether every glyph in the BorderStyle is empty, which is
// the signal used by Style to skip border rendering entirely.
func (s BorderStyle) IsZero() bool {
	return s.TopLeft == "" &&
		s.TopRight == "" &&
		s.BottomLeft == "" &&
		s.BottomRight == "" &&
		s.Top == "" &&
		s.Bottom == "" &&
		s.Left == "" &&
		s.Right == ""
}

// ---------------------------------------------------------------------------
// Preset constructors
// ---------------------------------------------------------------------------

// NoBorder returns a zero BorderStyle (all glyphs empty). Assigning this to a
// Style clears any previously set border.
//
//	s = s.WithBorderStyle(ink.NoBorder())
func NoBorder() BorderStyle {
	return BorderStyle{}
}

// BorderNormal returns the classic thin Unicode box-drawing border:
//
//	┌─────┐
//	│  …  │
//	└─────┘
func BorderNormal() BorderStyle {
	return BorderStyle{
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
	}
}

// BorderRounded returns a border with rounded corners:
//
//	╭─────╮
//	│  …  │
//	╰─────╯
func BorderRounded() BorderStyle {
	return BorderStyle{
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
	}
}

// BorderThick returns a heavy/thick Unicode box-drawing border:
//
//	┏━━━━━┓
//	┃  …  ┃
//	┗━━━━━┛
func BorderThick() BorderStyle {
	return BorderStyle{
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
		Top:         "━",
		Bottom:      "━",
		Left:        "┃",
		Right:       "┃",
	}
}

// BorderDouble returns a double-line Unicode box-drawing border:
//
//	╔═════╗
//	║  …  ║
//	╚═════╝
func BorderDouble() BorderStyle {
	return BorderStyle{
		TopLeft:     "╔",
		TopRight:    "╗",
		BottomLeft:  "╚",
		BottomRight: "╝",
		Top:         "═",
		Bottom:      "═",
		Left:        "║",
		Right:       "║",
	}
}

// BorderASCII returns a plain ASCII border, safe for terminals that lack
// Unicode support:
//
//	+-----+
//	|  …  |
//	+-----+
func BorderASCII() BorderStyle {
	return BorderStyle{
		TopLeft:     "+",
		TopRight:    "+",
		BottomLeft:  "+",
		BottomRight: "+",
		Top:         "-",
		Bottom:      "-",
		Left:        "|",
		Right:       "|",
	}
}

// BorderDashed returns a dashed/dotted border using Unicode light-dash glyphs:
//
//	┌┄┄┄┄┄┐
//	┆  …  ┆
//	└┄┄┄┄┄┘
func BorderDashed() BorderStyle {
	return BorderStyle{
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
		Top:         "┄",
		Bottom:      "┄",
		Left:        "┆",
		Right:       "┆",
	}
}

// BorderBlock returns a border built entirely from full-block characters.
// This gives a solid, filled appearance and works well with background colors:
//
//	███████
//	█  …  █
//	███████
func BorderBlock() BorderStyle {
	return BorderStyle{
		TopLeft:     "█",
		TopRight:    "█",
		BottomLeft:  "█",
		BottomRight: "█",
		Top:         "█",
		Bottom:      "█",
		Left:        "█",
		Right:       "█",
	}
}

// BorderHidden returns a border where every glyph is a single space. The border
// occupies space in the layout without drawing any visible line. This is useful
// for aligning content that sits alongside bordered blocks.
func BorderHidden() BorderStyle {
	return BorderStyle{
		TopLeft:     " ",
		TopRight:    " ",
		BottomLeft:  " ",
		BottomRight: " ",
		Top:         " ",
		Bottom:      " ",
		Left:        " ",
		Right:       " ",
	}
}

// BorderInnerHalfBlock returns a border that uses half-block characters to create
// an inner-shadow / inset look. Works best against a contrasting background:
//
// ▄▄▄▄▄▄▄
// ▌  …  ▐
// ▀▀▀▀▀▀▀
func BorderInnerHalfBlock() BorderStyle {
	return BorderStyle{
		TopLeft:     "▄",
		TopRight:    "▄",
		BottomLeft:  "▀",
		BottomRight: "▀",
		Top:         "▄",
		Bottom:      "▀",
		Left:        "▌",
		Right:       "▐",
	}
}
