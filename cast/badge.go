package cast

import "github.com/tidjee-dev/forge/ink"

// Badge is a short inline label with a styled (typically coloured) background,
// useful for status indicators such as "OK", "WARN", or "ERROR".
//
// Every setter returns a new Badge — the original is never mutated.
//
// Basic usage:
//
//	cast.NewBadge("OK").Success().Render()
//
// Custom style:
//
//	cast.NewBadge("v1.2.3").Style(ink.New().WithBackground(ink.Blue).WithForeground(ink.White)).Render()
type Badge struct {
	label string
	style ink.Style
}

// NewBadge returns a Badge with the given label and no styling applied.
// Calling Render on the zero-value result produces the plain label text.
func NewBadge(label string) Badge {
	return Badge{label: label}
}

// ---------------------------------------------------------------------------
// Setters
// ---------------------------------------------------------------------------

// Style replaces the entire style of the badge with s. It returns a new Badge.
func (b Badge) Style(s ink.Style) Badge {
	b.style = s
	return b
}

// Success applies the semantic success style: bright-green background, black
// foreground, bold text.
func (b Badge) Success() Badge {
	b.style = ink.New().
		WithBackground(ink.Success).
		WithForeground(ink.Black).
		WithBold(true)
	return b
}

// Warning applies the semantic warning style: amber background, black
// foreground, bold text.
func (b Badge) Warning() Badge {
	b.style = ink.New().
		WithBackground(ink.Warning).
		WithForeground(ink.Black).
		WithBold(true)
	return b
}

// Danger applies the semantic danger/error style: bright-red background, white
// foreground, bold text.
func (b Badge) Danger() Badge {
	b.style = ink.New().
		WithBackground(ink.Danger).
		WithForeground(ink.White).
		WithBold(true)
	return b
}

// Info applies the semantic info style: bright-blue background, white
// foreground, bold text.
func (b Badge) Info() Badge {
	b.style = ink.New().
		WithBackground(ink.Info).
		WithForeground(ink.White).
		WithBold(true)
	return b
}

// Neutral applies the semantic neutral style: mid-gray background, white
// foreground, bold text.
func (b Badge) Neutral() Badge {
	b.style = ink.New().
		WithBackground(ink.Muted).
		WithForeground(ink.White).
		WithBold(true)
	return b
}

// ---------------------------------------------------------------------------
// Render
// ---------------------------------------------------------------------------

// Render returns the badge as a styled string. A single space is added on
// each side of the label so the coloured background has visual breathing room.
// When colour is disabled the plain label is returned, still padded.
func (b Badge) Render() string {
	padded := " " + b.label + " "
	return b.style.Render(padded)
}
