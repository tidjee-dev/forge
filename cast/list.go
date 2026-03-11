package cast

import (
	"fmt"
	"strings"

	"github.com/tidjee-dev/forge/ink"
)

// List is a vertical list of items with an optional bullet character or
// sequential numbering.
//
// Every setter returns a new List — the original is never mutated.
//
// Basic usage:
//
//	cast.NewList("Buy carrots", "Buy celery", "Buy kohlrabi").
//	    Bullet("•").
//	    Render()
//
// Numbered variant:
//
//	cast.NewList("one", "two", "three").Numbered().Render()
type List struct {
	items       []string
	bullet      string
	numbered    bool
	itemStyle   ink.Style
	bulletStyle ink.Style
	indent      int
}

// defaultBullet is the bullet character used when none is specified and the
// list is not numbered.
const defaultBullet = "•"

// NewList returns a List pre-populated with the provided items and a default
// bullet character. Calling Render on the result immediately produces valid
// output.
func NewList(items ...string) List {
	// Make a defensive copy so the caller's slice cannot alias internal state.
	cp := make([]string, len(items))
	copy(cp, items)
	return List{
		items:  cp,
		bullet: defaultBullet,
	}
}

// ---------------------------------------------------------------------------
// Setters
// ---------------------------------------------------------------------------

// AddItem appends a single item to the list and returns a new List.
func (l List) AddItem(item string) List {
	items := make([]string, len(l.items)+1)
	copy(items, l.items)
	items[len(l.items)] = item
	l.items = items
	return l
}

// Bullet sets a custom bullet character. It implicitly disables Numbered mode.
// Pass an empty string to suppress the bullet entirely.
func (l List) Bullet(s string) List {
	l.bullet = s
	l.numbered = false
	return l
}

// Numbered switches the list to sequential numbering ("1.", "2.", …).
// It takes precedence over any Bullet setting while active.
func (l List) Numbered() List {
	l.numbered = true
	return l
}

// ItemStyle sets the ink.Style applied to each item's text.
func (l List) ItemStyle(s ink.Style) List {
	l.itemStyle = s
	return l
}

// BulletStyle sets the ink.Style applied to the bullet or number prefix.
// When not set the bullet inherits the ItemStyle.
func (l List) BulletStyle(s ink.Style) List {
	l.bulletStyle = s
	return l
}

// Indent sets the number of leading spaces prepended to every item line.
func (l List) Indent(n int) List {
	l.indent = max(n, 0)
	return l
}

// ---------------------------------------------------------------------------
// Render
// ---------------------------------------------------------------------------

// Render returns the list as a newline-separated string. Each item occupies
// exactly one line: [indent][bullet/number] [item text].
//
// When the list is empty an empty string is returned.
func (l List) Render() string {
	if len(l.items) == 0 {
		return ""
	}

	indentStr := strings.Repeat(" ", l.indent)

	// Decide whether we have a separate bullet style or fall back to itemStyle.
	hasBulletStyle := isStyleSet(l.bulletStyle)

	var sb strings.Builder
	for i, item := range l.items {
		// Build the prefix (bullet or number).
		var prefix string
		if l.numbered {
			prefix = fmt.Sprintf("%d.", i+1)
		} else {
			prefix = l.bullet
		}

		// Style the prefix.
		var styledPrefix string
		if prefix == "" {
			styledPrefix = ""
		} else if hasBulletStyle {
			styledPrefix = l.bulletStyle.Render(prefix)
		} else {
			styledPrefix = l.itemStyle.Render(prefix)
		}

		// Style the item text.
		styledItem := l.itemStyle.Render(item)

		// Assemble the line.
		if prefix == "" {
			sb.WriteString(indentStr + styledItem)
		} else {
			sb.WriteString(indentStr + styledPrefix + " " + styledItem)
		}

		if i < len(l.items)-1 {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}
