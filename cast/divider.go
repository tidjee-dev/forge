package cast

import (
	"strings"

	"github.com/tidjee-dev/forge/ink"
)

// Divider is a horizontal rule that spans a fixed or terminal width, with an
// optional inline label.
//
// Every setter returns a new Divider — the original is never mutated.
//
// Basic usage:
//
//	cast.NewDivider().Render()
//	cast.NewDivider().Label("Results").Render()
//	cast.NewDivider().Label("Results").LabelAlign(ink.JustifyStart).Render()
//	cast.NewDivider().Char("═").Style(ink.New().WithForeground(ink.Muted)).Render()
type Divider struct {
	label      string
	labelAlign ink.JustifyContent
	nearFill   int
	char       string
	width      int
	style      ink.Style
	labelStyle ink.Style
}

// defaultDividerChar is the fill character used when none is specified.
const defaultDividerChar = "─"

// defaultDividerWidth is the fallback width used when no width is set and
// terminal detection is unavailable.
const defaultDividerWidth = 80

// defaultNearFill is the number of fill characters placed on the near side of
// the label for JustifyStart and JustifyEnd alignments.
const defaultNearFill = 3

// NewDivider returns a Divider with default settings: "─" fill character,
// centred label alignment, 3-character near-side fill, and width determined
// at render time.
func NewDivider() Divider {
	return Divider{
		char:       defaultDividerChar,
		labelAlign: ink.JustifyCenter,
		nearFill:   defaultNearFill,
	}
}

// ---------------------------------------------------------------------------
// Setters
// ---------------------------------------------------------------------------

// Label sets an inline label placed within the rule according to LabelAlign
// (default: centred). Pass an empty string to remove a previously set label.
func (d Divider) Label(text string) Divider {
	d.label = text
	return d
}

// LabelAlign controls where the label sits within the rule:
//   - [ink.JustifyStart]  — label near the left, long fill on the right
//   - [ink.JustifyCenter] — label centred, fill split evenly (default)
//   - [ink.JustifyEnd]    — label near the right, long fill on the left
func (d Divider) LabelAlign(a ink.JustifyContent) Divider {
	d.labelAlign = a
	return d
}

// NearFill sets the number of fill characters placed on the short side of the
// label when using [ink.JustifyStart] or [ink.JustifyEnd] alignment. The
// default is 3, producing output like:
//
//	JustifyStart: "─── Label ──────────────────────"
//	JustifyEnd:   "────────────────────── Label ───"
//
// Values ≤ 0 are clamped to 0 (no near-side fill at all). Has no effect when
// alignment is [ink.JustifyCenter].
func (d Divider) NearFill(n int) Divider {
	d.nearFill = max(n, 0)
	return d
}

// Char sets the fill character(s) used to draw the rule. Defaults to "─".
// Multi-rune values are accepted but each rune should occupy exactly one
// terminal column for correct alignment.
func (d Divider) Char(s string) Divider {
	d.char = s
	return d
}

// Width sets a fixed render width for the divider in terminal columns.
// When 0 (the default) a fallback of 80 columns is used.
func (d Divider) Width(n int) Divider {
	n = max(n, 0)
	d.width = n
	return d
}

// Style sets the ink.Style applied to the fill characters of the rule.
func (d Divider) Style(s ink.Style) Divider {
	d.style = s
	return d
}

// LabelStyle sets the ink.Style applied to the label text only. When not set
// the label inherits the rule's Style.
func (d Divider) LabelStyle(s ink.Style) Divider {
	d.labelStyle = s
	return d
}

// ---------------------------------------------------------------------------
// Render
// ---------------------------------------------------------------------------

// Render returns the divider as a styled string. The rule occupies exactly
// Width() terminal columns. When a Label is set it is placed according to
// LabelAlign with a single space of padding on each side; the remaining
// columns are filled with the Char character.
func (d Divider) Render() string {
	width := d.width
	if width <= 0 {
		width = defaultDividerWidth
	}

	fillChar := d.char
	if fillChar == "" {
		fillChar = defaultDividerChar
	}
	charW := visibleWidth(fillChar)
	if charW == 0 {
		charW = 1
	}

	if d.label == "" {
		return d.renderPlain(width, fillChar, charW)
	}
	return d.renderLabelled(width, fillChar, charW)
}

// renderPlain builds a divider with no label.
func (d Divider) renderPlain(width int, fillChar string, charW int) string {
	count := width / charW
	count = max(count, 0)
	line := strings.Repeat(fillChar, count)
	return d.style.Render(line)
}

// renderLabelled builds a divider with a label positioned according to
// d.labelAlign. The label is padded with one space on each side.
//
// Layout examples (width = 80, label = "Results"):
//
//	JustifyStart:  "─── Results ────────────────────────────────────────────────────────────────────"
//	JustifyCenter: "─────────────────────────────────── Results ────────────────────────────────────"
//	JustifyEnd:    "──────────────────────────────────────────────────────────────────── Results ───"
func (d Divider) renderLabelled(width int, fillChar string, charW int) string {
	labelWithPad := " " + d.label + " "
	labelW := visibleWidth(labelWithPad)

	remaining := width - labelW
	remaining = max(remaining, 0)

	var leftCols, rightCols int

	switch d.labelAlign {
	case ink.JustifyStart:
		// Short near-side fill on the left, remainder on the right.
		nearCols := min(d.nearFill*charW, remaining)
		leftCols = nearCols
		rightCols = remaining - nearCols
	case ink.JustifyEnd:
		// Remainder on the left, short near-side fill on the right.
		nearCols := min(d.nearFill*charW, remaining)
		rightCols = nearCols
		leftCols = remaining - nearCols
	default: // JustifyCenter
		leftCols = remaining / 2
		rightCols = remaining - leftCols
	}

	leftCount := leftCols / charW
	rightCount := rightCols / charW

	leftFill := strings.Repeat(fillChar, leftCount)
	rightFill := strings.Repeat(fillChar, rightCount)

	var styledLabel string
	if !isStyleSet(d.labelStyle) {
		styledLabel = d.style.Render(labelWithPad)
	} else {
		styledLabel = d.labelStyle.Render(labelWithPad)
	}

	styledLeft := d.style.Render(leftFill)
	styledRight := d.style.Render(rightFill)

	return styledLeft + styledLabel + styledRight
}
