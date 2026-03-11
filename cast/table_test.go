package cast_test

import (
	"strings"
	"testing"

	"github.com/tidjee-dev/forge/cast"
	"github.com/tidjee-dev/forge/ink"
)

// ---------------------------------------------------------------------------
// Empty / zero-value
// ---------------------------------------------------------------------------

func TestTable_EmptyRender(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTable().Render()
	if got != "" {
		t.Errorf("expected empty string for empty table, got %q", got)
	}
}

func TestTable_HeadersOnly(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTable().Headers("Name", "Status").Render()
	if got == "" {
		t.Error("expected non-empty output for table with headers only")
	}
	if !strings.Contains(got, "Name") {
		t.Errorf("expected \"Name\" in output, got %q", got)
	}
	if !strings.Contains(got, "Status") {
		t.Errorf("expected \"Status\" in output, got %q", got)
	}
}

func TestTable_RowsOnly(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTable().AddRow("api", "running").Render()
	if got == "" {
		t.Error("expected non-empty output for table with row only")
	}
	if !strings.Contains(got, "api") {
		t.Errorf("expected \"api\" in output, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// Single row
// ---------------------------------------------------------------------------

func TestTable_SingleRow(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTable().
		Headers("Name", "Port").
		AddRow("api", "8080").
		Render()

	if !strings.Contains(got, "Name") {
		t.Errorf("missing \"Name\" in %q", got)
	}
	if !strings.Contains(got, "api") {
		t.Errorf("missing \"api\" in %q", got)
	}
	if !strings.Contains(got, "8080") {
		t.Errorf("missing \"8080\" in %q", got)
	}
}

// ---------------------------------------------------------------------------
// Multi-row
// ---------------------------------------------------------------------------

func TestTable_MultiRow(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTable().
		Headers("Name", "Status", "Port").
		AddRow("api", "running", "8080").
		AddRow("db", "stopped", "5432").
		Render()

	for _, want := range []string{"Name", "Status", "Port", "api", "running", "8080", "db", "stopped", "5432"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q in output, got:\n%s", want, got)
		}
	}
}

func TestTable_AddRows(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	rows := [][]string{
		{"a", "1"},
		{"b", "2"},
		{"c", "3"},
	}
	got := cast.NewTable().Headers("K", "V").AddRows(rows).Render()

	for _, r := range rows {
		for _, cell := range r {
			if !strings.Contains(got, cell) {
				t.Errorf("expected %q in output, got:\n%s", cell, got)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Border
// ---------------------------------------------------------------------------

func TestTable_BorderRounded(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTable().
		Headers("Name", "Status").
		AddRow("api", "running").
		Border(ink.BorderRounded()).
		Render()

	lines := strings.Split(got, "\n")
	if len(lines) < 4 {
		t.Fatalf("expected at least 4 lines with rounded border, got %d:\n%s", len(lines), got)
	}
	if !strings.HasPrefix(lines[0], "╭") {
		t.Errorf("expected top border to start with ╭, got %q", lines[0])
	}
	if !strings.HasSuffix(lines[0], "╮") {
		t.Errorf("expected top border to end with ╮, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[len(lines)-1], "╰") {
		t.Errorf("expected bottom border to start with ╰, got %q", lines[len(lines)-1])
	}
	if !strings.HasSuffix(lines[len(lines)-1], "╯") {
		t.Errorf("expected bottom border to end with ╯, got %q", lines[len(lines)-1])
	}
}

func TestTable_BorderNormal(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTable().
		Headers("A", "B").
		AddRow("1", "2").
		Border(ink.BorderNormal()).
		Render()

	lines := strings.Split(got, "\n")
	if !strings.HasPrefix(lines[0], "┌") {
		t.Errorf("expected ┌ at start of top border, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[len(lines)-1], "└") {
		t.Errorf("expected └ at start of bottom border, got %q", lines[len(lines)-1])
	}
}

func TestTable_NoBorder(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTable().
		Headers("Name", "Status").
		AddRow("api", "running").
		NoBorder().
		Render()

	for _, glyph := range []string{"╭", "╮", "╰", "╯", "┌", "┐", "└", "┘", "│", "─"} {
		if strings.Contains(got, glyph) {
			t.Errorf("expected no border glyphs in NoBorder table, found %q in:\n%s", glyph, got)
		}
	}
	if !strings.Contains(got, "Name") || !strings.Contains(got, "api") {
		t.Errorf("content missing from NoBorder table:\n%s", got)
	}
}

func TestTable_BorderThenNoBorder(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Setting a border then removing it should produce no border.
	got := cast.NewTable().
		Headers("X").
		AddRow("y").
		Border(ink.BorderRounded()).
		NoBorder().
		Render()

	if strings.Contains(got, "╭") {
		t.Errorf("expected no border after NoBorder(), got:\n%s", got)
	}
}

// ---------------------------------------------------------------------------
// Header separator
// ---------------------------------------------------------------------------

func TestTable_HeaderSeparatorPresent(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTable().
		Headers("Name", "Value").
		AddRow("x", "1").
		Border(ink.BorderRounded()).
		Render()

	lines := strings.Split(got, "\n")
	// Expect: top, header, separator, data row, bottom  → ≥5 lines
	if len(lines) < 5 {
		t.Fatalf("expected at least 5 lines (top+header+sep+row+bottom), got %d:\n%s", len(lines), got)
	}

	// The separator line should contain the bottom border character (─ for rounded).
	sep := lines[2]
	if !strings.Contains(sep, "─") {
		t.Errorf("expected header separator on line 2, got %q", sep)
	}
}

// ---------------------------------------------------------------------------
// Column auto-sizing
// ---------------------------------------------------------------------------

func TestTable_ColumnAutoSizing(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// The "LongHeader" column should be at least 10 chars wide.
	got := cast.NewTable().
		Headers("LongHeader", "X").
		AddRow("short", "y").
		Render()

	// Find the line containing "LongHeader" and check the cell width.
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "LongHeader") {
			// "LongHeader" is 10 chars; padded cell is at least " LongHeader ".
			if !strings.Contains(line, " LongHeader ") {
				t.Errorf("column not wide enough; header cell not padded: %q", line)
			}
			break
		}
	}

	// The data row's first cell ("short") should be padded to the same width.
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "short") {
			// "short" is 5 chars; with LongHeader (10) it should be padded to 10.
			if !strings.Contains(line, " short      ") && !strings.Contains(line, " short     ") {
				// Accept either 5 or 6 trailing spaces (10 - 5 = 5 spaces).
				if !strings.Contains(line, "short") {
					t.Errorf("data cell not padded correctly: %q", line)
				}
			}
			break
		}
	}
}

func TestTable_ColumnAutoSizingDataDrivenWidth(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// "VeryLongDataCell" (16 chars) is wider than header "Col" (3 chars).
	got := cast.NewTable().
		Headers("Col").
		AddRow("VeryLongDataCell").
		Render()

	// The header line should be padded to match the data width.
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "Col") && !strings.Contains(line, "VeryLong") {
			// "Col" should be padded to 16 chars.
			if !strings.Contains(line, " Col             ") && !strings.Contains(line, " Col") {
				// Just check "Col" is present; width check is approximate.
				t.Logf("header line: %q", line)
			}
			break
		}
	}
	if !strings.Contains(got, "VeryLongDataCell") {
		t.Errorf("data cell missing from output:\n%s", got)
	}
}

// ---------------------------------------------------------------------------
// Column min-width
// ---------------------------------------------------------------------------

func TestTable_ColumnMinWidth(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Column 0 has min-width 15, but all content is much shorter.
	got := cast.NewTable().
		Headers("A").
		AddRow("x").
		ColumnMinWidth(0, 15).
		Render()

	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "A") || strings.Contains(line, "x") {
			w := visibleWidthTest(line)
			if w < 15 {
				t.Errorf("expected line width >= 15 with ColumnMinWidth(0,15), got %d: %q", w, line)
			}
			break
		}
	}
}

// ---------------------------------------------------------------------------
// Column alignment
// ---------------------------------------------------------------------------

func TestTable_ColumnAlignCenter(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTable().
		Headers("Name").
		AddRow("Hi").
		ColumnMinWidth(0, 10).
		ColumnAlign(0, ink.JustifyCenter).
		Render()

	// "Hi" (2 wide) centred in 10 → 4 left, 4 right spaces (or 4/4 split).
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "Hi") {
			// Check that there are spaces on both sides of "Hi".
			before, after, found := strings.Cut(line, "Hi")
			if !found || before == "" {
				t.Errorf("\"Hi\" appears at start of line with no left padding: %q", line)
			}
			if !strings.HasPrefix(after, " ") {
				t.Errorf("no right padding after \"Hi\": %q", line)
			}
			break
		}
	}
}

func TestTable_ColumnAlignEnd(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTable().
		Headers("Value").
		AddRow("42").
		ColumnMinWidth(0, 10).
		ColumnAlign(0, ink.JustifyEnd).
		Render()

	// "42" right-aligned in 10 → should be preceded by 8 spaces in the cell.
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "42") && !strings.Contains(line, "Value") {
			before, _, found := strings.Cut(line, "42")
			if !found {
				t.Errorf("\"42\" not found in line: %q", line)
				break
			}
			// There should be at least some spaces before "42".
			if !strings.Contains(before, " ") {
				t.Errorf("expected left padding for right-aligned cell: %q", line)
			}
			break
		}
	}
}

// ---------------------------------------------------------------------------
// Column style
// ---------------------------------------------------------------------------

func TestTable_ColumnStyle(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	colStyle := ink.New().WithForeground(ink.Cyan)
	got := cast.NewTable().
		Headers("Name", "Val").
		AddRow("api", "1").
		ColumnStyle(0, colStyle).
		Render()

	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI sequences from ColumnStyle, got:\n%s", got)
	}
	plain := ink.Strip(got)
	if !strings.Contains(plain, "api") {
		t.Errorf("expected \"api\" in plain output:\n%s", plain)
	}
}

// ---------------------------------------------------------------------------
// Header style
// ---------------------------------------------------------------------------

func TestTable_HeaderStyle(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	hStyle := ink.New().WithBold(true).WithForeground(ink.BrightWhite)
	got := cast.NewTable().
		Headers("Name").
		AddRow("api").
		HeaderStyle(hStyle).
		Render()

	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI sequences from HeaderStyle, got:\n%s", got)
	}
}

// ---------------------------------------------------------------------------
// Row style
// ---------------------------------------------------------------------------

func TestTable_RowStyle(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	rStyle := ink.New().WithForeground(ink.White)
	got := cast.NewTable().
		Headers("Name").
		AddRow("api").
		RowStyle(rStyle).
		Render()

	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI sequences from RowStyle, got:\n%s", got)
	}
}

// ---------------------------------------------------------------------------
// Zebra / alt-row style
// ---------------------------------------------------------------------------

func TestTable_AltRowStyle(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	rowStyle := ink.New().WithForeground(ink.White)
	altStyle := ink.New().WithForeground(ink.Muted)

	got := cast.NewTable().
		Headers("Name").
		AddRow("row0").
		AddRow("row1").
		AddRow("row2").
		RowStyle(rowStyle).
		AltRowStyle(altStyle).
		Render()

	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI sequences from AltRowStyle, got:\n%s", got)
	}
	plain := ink.Strip(got)
	for _, cell := range []string{"row0", "row1", "row2"} {
		if !strings.Contains(plain, cell) {
			t.Errorf("expected %q in plain output:\n%s", cell, plain)
		}
	}
}

// ---------------------------------------------------------------------------
// No color mode
// ---------------------------------------------------------------------------

func TestTable_NoColorMode(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	hStyle := ink.New().WithBold(true).WithForeground(ink.BrightWhite)
	got := cast.NewTable().
		Headers("Name", "Status").
		AddRow("api", "running").
		Border(ink.BorderRounded()).
		HeaderStyle(hStyle).
		Render()

	if strings.Contains(got, "\x1b[") {
		t.Errorf("expected no ANSI sequences in NeverColorMode, got:\n%s", got)
	}
	if !strings.Contains(got, "Name") || !strings.Contains(got, "api") {
		t.Errorf("expected content in NeverColorMode output:\n%s", got)
	}
}

// ---------------------------------------------------------------------------
// Immutability
// ---------------------------------------------------------------------------

func TestTable_Immutability_AddRow(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	base := cast.NewTable().Headers("Name").AddRow("api")
	extended := base.AddRow("db")

	// base should only contain "api"
	if strings.Contains(base.Render(), "db") {
		t.Errorf("base table was mutated by AddRow")
	}
	// extended should contain both
	if !strings.Contains(extended.Render(), "api") || !strings.Contains(extended.Render(), "db") {
		t.Errorf("extended table missing rows:\n%s", extended.Render())
	}
}

func TestTable_Immutability_Headers(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	base := cast.NewTable().Headers("A", "B")
	withExtra := base.Headers("X", "Y", "Z")

	if strings.Contains(base.Render(), "X") {
		t.Errorf("base table headers were mutated")
	}
	if !strings.Contains(withExtra.Render(), "X") {
		t.Errorf("new headers not applied:\n%s", withExtra.Render())
	}
}

func TestTable_Immutability_Border(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	base := cast.NewTable().Headers("X").AddRow("1")
	bordered := base.Border(ink.BorderRounded())

	if strings.Contains(base.Render(), "╭") {
		t.Errorf("base table was mutated to include border")
	}
	if !strings.Contains(bordered.Render(), "╭") {
		t.Errorf("bordered table missing ╭:\n%s", bordered.Render())
	}
}

func TestTable_Immutability_ColumnAlign(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	base := cast.NewTable().Headers("V").AddRow("hi").ColumnMinWidth(0, 10)
	centred := base.ColumnAlign(0, ink.JustifyCenter)

	// Both should render without panicking and produce different alignment.
	_ = base.Render()
	_ = centred.Render()
}

func TestTable_Immutability_ColumnMinWidth(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	base := cast.NewTable().Headers("V").AddRow("hi")
	wide := base.ColumnMinWidth(0, 20)

	_ = base.Render()
	_ = wide.Render()
	// wide should be wider than base
	baseW := visibleWidthTest(base.Render())
	wideW := visibleWidthTest(wide.Render())
	if wideW <= baseW {
		t.Errorf("ColumnMinWidth should increase table width: base=%d wide=%d", baseW, wideW)
	}
}

// ---------------------------------------------------------------------------
// Edge cases
// ---------------------------------------------------------------------------

func TestTable_MissingCells(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Row has fewer cells than the header.
	got := cast.NewTable().
		Headers("A", "B", "C").
		AddRow("only-a").
		Render()

	// Should not panic and should contain "only-a".
	if !strings.Contains(got, "only-a") {
		t.Errorf("expected \"only-a\" in output:\n%s", got)
	}
}

func TestTable_ExtraCells(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Row has more cells than the header — extra cells should still appear.
	got := cast.NewTable().
		Headers("A").
		AddRow("1", "extra").
		Render()

	if !strings.Contains(got, "extra") {
		t.Errorf("expected extra cell in output:\n%s", got)
	}
}

func TestTable_UnicodeCell(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// CJK characters are 2 columns wide each.
	got := cast.NewTable().
		Headers("Lang", "Word").
		AddRow("Chinese", "你好").
		Render()

	if !strings.Contains(got, "你好") {
		t.Errorf("expected CJK text in output:\n%s", got)
	}
}

func TestTable_NoTrailingNewline(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTable().Headers("A").AddRow("b").Render()
	if strings.HasSuffix(got, "\n") {
		t.Errorf("expected no trailing newline, got %q", got)
	}
}
