package cast

import (
	"strings"

	"github.com/tidjee-dev/forge/ink"
)

// Banner is a full-width styled block, useful for section headers or
// prominent messages. It supports borders, alignment, and a fixed width.
//
// Every setter returns a new Banner — the original is never mutated.
//
// Basic usage:
//
//	cast.NewBanner("Server started on :8080").
//	    Border(ink.BorderRounded()).
//	    Width(40).
//	    Render()
type Banner struct {
	text        string
	style       ink.Style
	borderStyle ink.BorderStyle
	width       int
	align       ink.JustifyContent
}

// NewBanner returns a Banner with the given text and no styling applied.
// Calling Render on the zero-value result produces the plain text.
func NewBanner(text string) Banner {
	return Banner{text: text}
}

// ---------------------------------------------------------------------------
// Setters
// ---------------------------------------------------------------------------

// Style replaces the entire style of the banner with s. It returns a new Banner.
// Note: border and alignment are managed separately; Style controls text
// attributes (foreground, background, bold, etc.).
func (b Banner) Style(s ink.Style) Banner {
	b.style = s
	return b
}

// Border sets the border glyph set drawn around the banner. Pass
// ink.NoBorder() to explicitly remove a previously set border.
func (b Banner) Border(bs ink.BorderStyle) Banner {
	b.borderStyle = bs
	return b
}

// Width sets a fixed render width for the banner in terminal columns.
// When 0 (the default) the banner is as wide as its content plus padding.
func (b Banner) Width(n int) Banner {
	b.width = max(n, 0)
	return b
}

// Align controls the horizontal alignment of the text within the banner's
// content area. Use ink.JustifyStart (default), ink.JustifyCenter, or
// ink.JustifyEnd.
func (b Banner) Align(a ink.JustifyContent) Banner {
	b.align = a
	return b
}

// ---------------------------------------------------------------------------
// Render
// ---------------------------------------------------------------------------

// Render returns the banner as a string, applying style, optional border, and
// alignment. When a fixed Width is set the content area is padded or truncated
// to that width.
func (b Banner) Render() string {
	hasBorder := !b.borderStyle.IsZero()

	// Determine the inner content width.
	// If a fixed width is requested, the inner width accounts for the border
	// columns that will be added (1 left + 1 right when both sides are drawn).
	innerW := b.width
	if hasBorder && innerW > 0 {
		// Subtract left and right border columns so the overall rendered width
		// equals b.width.
		innerW = max(innerW-2, 0)
	}

	// Build the content line with alignment applied.
	content := b.renderContent(innerW)

	// Apply text style to the content.
	styled := b.style.Render(content)

	if !hasBorder {
		return styled
	}

	// Draw the border manually so that only the content portion is coloured
	// by b.style, while the border glyphs themselves are unstyled (the caller
	// may colour them via ink.Style.WithBorderColor if they wish — but since
	// Banner manages its own border here we keep it simple: no separate border
	// colour).
	return renderBorder(styled, b.borderStyle, visibleWidth(content))
}

// renderContent builds the (unstyled) text line, padded to innerW columns
// according to the chosen alignment. When innerW == 0 no padding is added.
func (b Banner) renderContent(innerW int) string {
	textW := visibleWidth(b.text)

	if innerW <= 0 {
		// No fixed width — just return the text.
		return b.text
	}

	if textW >= innerW {
		return b.text
	}

	switch b.align {
	case ink.JustifyCenter:
		return centerPad(b.text, innerW)
	case ink.JustifyEnd:
		return padLeft(b.text, innerW)
	default: // JustifyStart
		return padRight(b.text, innerW)
	}
}

// renderBorder wraps styledContent in a box drawn with bs. styledContent may
// contain ANSI sequences; innerW is its ANSI-stripped visible column count so
// that border lines are the correct length.
func renderBorder(styledContent string, bs ink.BorderStyle, innerW int) string {
	var sb strings.Builder

	top := bs.TopLeft + repeatStr(bs.Top, innerW) + bs.TopRight
	sb.WriteString(top)
	sb.WriteByte('\n')

	sb.WriteString(bs.Left)
	sb.WriteString(styledContent)
	sb.WriteString(bs.Right)
	sb.WriteByte('\n')

	bottom := bs.BottomLeft + repeatStr(bs.Bottom, innerW) + bs.BottomRight
	sb.WriteString(bottom)

	return sb.String()
}
