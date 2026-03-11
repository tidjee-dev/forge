package cast

import (
	"strings"

	"github.com/tidjee-dev/forge/ink"
)

// Table is a bordered, aligned tabular data renderer. It supports headers,
// multiple data rows, optional zebra striping, per-column alignment and
// minimum widths, and several border styles.
//
// Every setter returns a new Table — the original is never mutated.
//
// Basic usage:
//
//	cast.NewTable().
//	    Headers("Name", "Status", "Port").
//	    Border(ink.BorderRounded()).
//	    HeaderStyle(ink.New().WithBold(true).WithForeground(ink.BrightWhite)).
//	    AddRow("api", "running", "8080").
//	    AddRow("db",  "stopped", "5432").
//	    Render()
type Table struct {
	headers      []string
	rows         [][]string
	borderStyle  ink.BorderStyle
	noBorder     bool
	headerStyle  ink.Style
	rowStyle     ink.Style
	altRowStyle  ink.Style
	colAligns    map[int]ink.JustifyContent
	colStyles    map[int]ink.Style
	colMinWidths map[int]int
}

// NewTable returns an empty Table with no headers, no rows, and no border
// applied. Calling Render on the zero-value result produces an empty string.
func NewTable() Table {
	return Table{
		colAligns:    make(map[int]ink.JustifyContent),
		colStyles:    make(map[int]ink.Style),
		colMinWidths: make(map[int]int),
	}
}

// ---------------------------------------------------------------------------
// Setters
// ---------------------------------------------------------------------------

// Headers sets the column headers. Calling Headers again replaces them.
func (t Table) Headers(cols ...string) Table {
	cp := make([]string, len(cols))
	copy(cp, cols)
	t.headers = cp
	return t
}

// AddRow appends a single data row. Extra cells beyond the header count are
// retained; missing cells are rendered as empty strings.
func (t Table) AddRow(cells ...string) Table {
	cp := make([]string, len(cells))
	copy(cp, cells)
	rows := make([][]string, len(t.rows)+1)
	copy(rows, t.rows)
	rows[len(t.rows)] = cp
	t.rows = rows
	return t
}

// AddRows appends multiple data rows at once.
func (t Table) AddRows(rows [][]string) Table {
	for _, r := range rows {
		t = t.AddRow(r...)
	}
	return t
}

// Border sets the border glyph set drawn around and inside the table. Pass
// ink.NoBorder() to explicitly remove a previously set border.
func (t Table) Border(bs ink.BorderStyle) Table {
	t.borderStyle = bs
	t.noBorder = false
	return t
}

// NoBorder removes any border from the table.
func (t Table) NoBorder() Table {
	t.borderStyle = ink.NoBorder()
	t.noBorder = true
	return t
}

// HeaderStyle sets the ink.Style applied to every cell in the header row.
func (t Table) HeaderStyle(s ink.Style) Table {
	t.headerStyle = s
	return t
}

// RowStyle sets the ink.Style applied to all data rows.
func (t Table) RowStyle(s ink.Style) Table {
	t.rowStyle = s
	return t
}

// AltRowStyle sets the ink.Style applied to alternating (even-indexed) data
// rows, creating a zebra-striping effect. Odd rows use RowStyle.
func (t Table) AltRowStyle(s ink.Style) Table {
	t.altRowStyle = s
	return t
}

// ColumnAlign sets the horizontal alignment for column col (0-indexed).
// Use ink.JustifyStart (default), ink.JustifyCenter, or ink.JustifyEnd.
func (t Table) ColumnAlign(col int, a ink.JustifyContent) Table {
	m := copyIntJustifyMap(t.colAligns)
	m[col] = a
	t.colAligns = m
	return t
}

// ColumnStyle sets the ink.Style applied to all cells in column col
// (0-indexed). This overrides the row style for that column.
func (t Table) ColumnStyle(col int, s ink.Style) Table {
	m := copyIntStyleMap(t.colStyles)
	m[col] = s
	t.colStyles = m
	return t
}

// ColumnMinWidth sets a minimum rendered width (in terminal columns) for
// column col (0-indexed). The column will be at least this wide even if all
// its cells are narrower.
func (t Table) ColumnMinWidth(col int, n int) Table {
	m := copyIntIntMap(t.colMinWidths)
	n = max(n, 0)
	m[col] = n
	t.colMinWidths = m
	return t
}

// ---------------------------------------------------------------------------
// Render
// ---------------------------------------------------------------------------

// Render returns the table as a multi-line string.
//
// Column widths are auto-sized to the widest cell (header or data) in each
// column, then clamped up to any ColumnMinWidth. ANSI sequences are stripped
// before measurement so styled input does not inflate column widths.
//
// When the table has no headers and no rows, an empty string is returned.
func (t Table) Render() string {
	colCount := t.columnCount()
	if colCount == 0 {
		return ""
	}

	// Compute per-column widths.
	widths := t.computeWidths(colCount)

	hasBorder := !t.noBorder && !t.borderStyle.IsZero()

	var sb strings.Builder

	if hasBorder {
		// Top border row: ╭──────┬─────────┬──────╮
		sb.WriteString(t.borderStyle.TopLeft)
		for ci := 0; ci < colCount; ci++ {
			sb.WriteString(strings.Repeat(t.borderStyle.Top, widths[ci]+2)) // +2 for padding
			if ci < colCount-1 {
				// Inner top junction — reuse Top char as fallback since
				// ink.BorderStyle doesn't expose a TopJunction glyph.
				// We approximate with the Top glyph; callers wanting a true
				// ┬ need a custom border.
				sb.WriteString(topJunction(t.borderStyle))
			}
		}
		sb.WriteString(t.borderStyle.TopRight)
		sb.WriteByte('\n')
	}

	// Header row ──────────────────────────────────────────────────────────
	if len(t.headers) > 0 {
		t.writeRow(&sb, t.headers, colCount, widths, t.headerStyle, hasBorder)

		if hasBorder {
			// Header separator: ├──────┼─────────┼──────┤
			sb.WriteString(t.borderStyle.Left)
			for ci := 0; ci < colCount; ci++ {
				sb.WriteString(strings.Repeat(t.borderStyle.Bottom, widths[ci]+2))
				if ci < colCount-1 {
					sb.WriteString(midJunction(t.borderStyle))
				}
			}
			sb.WriteString(t.borderStyle.Right)
			sb.WriteByte('\n')
		}
	}

	// Data rows ───────────────────────────────────────────────────────────
	hasAlt := isStyleSet(t.altRowStyle)

	for ri, row := range t.rows {
		var rowStyle ink.Style
		if hasAlt && ri%2 == 1 {
			rowStyle = t.altRowStyle
		} else {
			rowStyle = t.rowStyle
		}
		t.writeRow(&sb, row, colCount, widths, rowStyle, hasBorder)
	}

	if hasBorder {
		// Bottom border ───────────────────────────────────────────────────
		// ╰──────┴─────────┴──────╯
		sb.WriteString(t.borderStyle.BottomLeft)
		for ci := 0; ci < colCount; ci++ {
			sb.WriteString(strings.Repeat(t.borderStyle.Bottom, widths[ci]+2))
			if ci < colCount-1 {
				sb.WriteString(bottomJunction(t.borderStyle))
			}
		}
		sb.WriteString(t.borderStyle.BottomRight)
		sb.WriteByte('\n')
	}

	result := sb.String()
	// Trim trailing newline.
	result = strings.TrimRight(result, "\n")
	return result
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// writeRow writes one row (header or data) to sb.
func (t Table) writeRow(sb *strings.Builder, cells []string, colCount int, widths []int, rowSty ink.Style, hasBorder bool) {
	if hasBorder {
		sb.WriteString(t.borderStyle.Left)
	}

	for ci := 0; ci < colCount; ci++ {
		var cell string
		if ci < len(cells) {
			cell = cells[ci]
		}

		// Align the cell within its column width.
		align := t.colAligns[ci] // zero value = JustifyStart
		aligned := alignCell(cell, widths[ci], align)

		// Determine which style to use: column style > row style > none.
		var styled string
		if colSty, ok := t.colStyles[ci]; ok && isStyleSet(colSty) {
			styled = colSty.Render(aligned)
		} else if isStyleSet(rowSty) {
			styled = rowSty.Render(aligned)
		} else {
			styled = aligned
		}

		// Write " cell " (one space padding each side).
		sb.WriteByte(' ')
		sb.WriteString(styled)
		sb.WriteByte(' ')

		if hasBorder && ci < colCount-1 {
			sb.WriteString(t.borderStyle.Right)
		}
	}

	if hasBorder {
		sb.WriteString(t.borderStyle.Right)
	}
	sb.WriteByte('\n')
}

// alignCell pads or aligns cell text within a field of width w.
func alignCell(cell string, w int, align ink.JustifyContent) string {
	switch align {
	case ink.JustifyCenter:
		return centerPad(cell, w)
	case ink.JustifyEnd:
		return padLeft(cell, w)
	default:
		return padRight(cell, w)
	}
}

// columnCount returns the number of columns inferred from headers and rows.
func (t Table) columnCount() int {
	n := len(t.headers)
	for _, row := range t.rows {
		if len(row) > n {
			n = len(row)
		}
	}
	return n
}

// computeWidths calculates the rendered column widths. Each column is as wide
// as its widest cell (header or data), clamped up to any ColumnMinWidth. ANSI
// sequences are stripped before measurement.
func (t Table) computeWidths(colCount int) []int {
	widths := make([]int, colCount)

	// Seed from headers.
	for ci, h := range t.headers {
		if ci < colCount {
			w := visibleWidth(h)
			if w > widths[ci] {
				widths[ci] = w
			}
		}
	}

	// Seed from data rows.
	for _, row := range t.rows {
		for ci, cell := range row {
			if ci < colCount {
				w := visibleWidth(cell)
				if w > widths[ci] {
					widths[ci] = w
				}
			}
		}
	}

	// Apply minimum widths.
	for ci, minW := range t.colMinWidths {
		if ci < colCount && minW > widths[ci] {
			widths[ci] = minW
		}
	}

	return widths
}

// ---------------------------------------------------------------------------
// Border junction helpers
//
// ink.BorderStyle does not expose explicit junction glyphs (┬ ┼ ┴ ├ ┤).
// We derive sensible defaults by inspecting the border characters so that
// common border styles produce correct-looking output without requiring
// callers to specify junction glyphs explicitly.
// ---------------------------------------------------------------------------

// topJunction returns a best-guess top-inner junction glyph (┬) based on the
// border style's Top character. Falls back to the Top char itself.
func topJunction(bs ink.BorderStyle) string {
	switch bs.Top {
	case "─":
		return "┬"
	case "━":
		return "┳"
	case "═":
		return "╦"
	case "┄":
		return "┬"
	case "█":
		return "█"
	case " ":
		return " "
	default:
		return bs.Top
	}
}

// midJunction returns a best-guess middle-inner junction glyph (┼).
func midJunction(bs ink.BorderStyle) string {
	switch bs.Top {
	case "─":
		return "┼"
	case "━":
		return "╋"
	case "═":
		return "╬"
	case "┄":
		return "┼"
	case "█":
		return "█"
	case " ":
		return " "
	default:
		return bs.Left
	}
}

// bottomJunction returns a best-guess bottom-inner junction glyph (┴).
func bottomJunction(bs ink.BorderStyle) string {
	switch bs.Top {
	case "─":
		return "┴"
	case "━":
		return "┻"
	case "═":
		return "╩"
	case "┄":
		return "┴"
	case "█":
		return "█"
	case " ":
		return " "
	default:
		return bs.Bottom
	}
}

// ---------------------------------------------------------------------------
// Map copy helpers (preserve immutability)
// ---------------------------------------------------------------------------

func copyIntJustifyMap(m map[int]ink.JustifyContent) map[int]ink.JustifyContent {
	out := make(map[int]ink.JustifyContent, len(m)+1)
	for k, v := range m {
		out[k] = v
	}
	return out
}

func copyIntStyleMap(m map[int]ink.Style) map[int]ink.Style {
	out := make(map[int]ink.Style, len(m)+1)
	for k, v := range m {
		out[k] = v
	}
	return out
}

func copyIntIntMap(m map[int]int) map[int]int {
	out := make(map[int]int, len(m)+1)
	for k, v := range m {
		out[k] = v
	}
	return out
}
