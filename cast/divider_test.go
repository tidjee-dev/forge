package cast_test

import (
	"strings"
	"testing"

	"github.com/tidjee-dev/forge/cast"
	"github.com/tidjee-dev/forge/ink"
)

func TestDivider_DefaultRender(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewDivider().Width(10).Render()
	if got != "──────────" {
		t.Errorf("expected 10 dashes, got %q", got)
	}
}

func TestDivider_DefaultWidth(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewDivider().Render()
	w := visibleWidthTest(got)
	if w != 80 {
		t.Errorf("expected default width 80, got %d", w)
	}
}

func TestDivider_CustomWidth(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewDivider().Width(20).Render()
	w := visibleWidthTest(got)
	if w != 20 {
		t.Errorf("expected width 20, got %d: %q", w, got)
	}
}

func TestDivider_CustomChar(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewDivider().Char("═").Width(5).Render()
	if got != "═════" {
		t.Errorf("expected 5x═, got %q", got)
	}
}

func TestDivider_CustomCharASCII(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewDivider().Char("-").Width(10).Render()
	if got != "----------" {
		t.Errorf("expected 10 hyphens, got %q", got)
	}
}

func TestDivider_Label(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewDivider().Label("Results").Width(30).Render()

	// The label should appear somewhere in the output.
	if !strings.Contains(got, "Results") {
		t.Errorf("expected label \"Results\" in output, got %q", got)
	}

	// Fill chars should appear on both sides.
	if !strings.Contains(got, "─") {
		t.Errorf("expected fill chars around label, got %q", got)
	}

	// The label should be surrounded by spaces.
	if !strings.Contains(got, " Results ") {
		t.Errorf("expected label padded with spaces, got %q", got)
	}
}

func TestDivider_LabelTotalWidth(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// The total visible width should equal the requested width.
	// With a label, the fill is split left + right of the label.
	got := cast.NewDivider().Label("Hi").Width(20).Render()
	w := visibleWidthTest(got)
	// The actual width may be slightly less than 20 because integer division
	// of the fill columns can leave up to 1 column off per side. We accept
	// anything in [18, 20].
	if w < 18 || w > 20 {
		t.Errorf("expected total width ~20 (18–20), got %d: %q", w, got)
	}
}

func TestDivider_LabelAlignCenter(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Default alignment is center — fill on both sides.
	got := cast.NewDivider().Label("X").Width(21).Render()

	leftPart, rightPart, ok := strings.Cut(got, "X")
	if !ok {
		t.Fatalf("label not found in %q", got)
	}

	if strings.TrimSpace(leftPart) == "" {
		t.Errorf("center: expected fill to the left of label, got %q", got)
	}
	if strings.TrimSpace(rightPart) == "" {
		t.Errorf("center: expected fill to the right of label, got %q", got)
	}
}

func TestDivider_LabelAlignCenterExplicit(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	implicit := cast.NewDivider().Label("Hi").Width(20).Render()
	explicit := cast.NewDivider().Label("Hi").LabelAlign(ink.JustifyCenter).Width(20).Render()

	if implicit != explicit {
		t.Errorf("explicit JustifyCenter should match default: implicit=%q explicit=%q", implicit, explicit)
	}
}

func TestDivider_LabelAlignStart(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewDivider().Label("Hi").LabelAlign(ink.JustifyStart).Width(20).Render()

	// Small fill on the left, large fill on the right.
	// e.g. "─── Hi ──────────────"
	leftPart, rightPart, ok := strings.Cut(got, "Hi")
	if !ok {
		t.Fatalf("start: label not found in %q", got)
	}

	if !strings.Contains(leftPart, "─") {
		t.Errorf("start: expected short fill to the left of label, got %q", got)
	}
	if !strings.Contains(rightPart, "─") {
		t.Errorf("start: expected fill to the right of label, got %q", got)
	}
	// Far side (right) should be longer than near side (left).
	if visibleWidthTest(strings.TrimSpace(rightPart)) <= visibleWidthTest(strings.TrimSpace(leftPart)) {
		t.Errorf("start: expected right fill to be longer than left fill, got %q", got)
	}
}

func TestDivider_LabelAlignEnd(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewDivider().Label("Hi").LabelAlign(ink.JustifyEnd).Width(20).Render()

	// Large fill on the left, small fill on the right.
	// e.g. "────────────────── Hi ───"
	leftPart, rightPart, ok := strings.Cut(got, "Hi")
	if !ok {
		t.Fatalf("end: label not found in %q", got)
	}

	if !strings.Contains(leftPart, "─") {
		t.Errorf("end: expected fill to the left of label, got %q", got)
	}
	if !strings.Contains(rightPart, "─") {
		t.Errorf("end: expected short fill to the right of label, got %q", got)
	}
	// Far side (left) should be longer than near side (right).
	if visibleWidthTest(strings.TrimSpace(leftPart)) <= visibleWidthTest(strings.TrimSpace(rightPart)) {
		t.Errorf("end: expected left fill to be longer than right fill, got %q", got)
	}
}

func TestDivider_LabelAlignStartTotalWidth(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewDivider().Label("Hi").LabelAlign(ink.JustifyStart).Width(20).Render()
	w := visibleWidthTest(got)
	if w < 18 || w > 20 {
		t.Errorf("start align: expected total width ~20 (18–20), got %d: %q", w, got)
	}
}

func TestDivider_LabelAlignEndTotalWidth(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewDivider().Label("Hi").LabelAlign(ink.JustifyEnd).Width(20).Render()
	w := visibleWidthTest(got)
	if w < 18 || w > 20 {
		t.Errorf("end align: expected total width ~20 (18–20), got %d: %q", w, got)
	}
}

func TestDivider_LabelAlignImmutability(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	base := cast.NewDivider().Label("Hi").Width(20)
	start := base.LabelAlign(ink.JustifyStart)
	end := base.LabelAlign(ink.JustifyEnd)

	// base should still be centered
	baseOut := base.Render()
	startOut := start.Render()
	endOut := end.Render()

	if baseOut == startOut {
		t.Errorf("expected base (center) to differ from start: %q", baseOut)
	}
	if baseOut == endOut {
		t.Errorf("expected base (center) to differ from end: %q", baseOut)
	}
	if startOut == endOut {
		t.Errorf("expected start and end to differ: %q", startOut)
	}
}

func TestDivider_NearFillCustom(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewDivider().Label("Hi").LabelAlign(ink.JustifyStart).NearFill(6).Width(20).Render()

	leftPart, rightPart, ok := strings.Cut(got, "Hi")
	if !ok {
		t.Fatalf("label not found in %q", got)
	}

	// Left side should be longer than default (3) — we asked for 6.
	if visibleWidthTest(strings.TrimSpace(leftPart)) <= 3 {
		t.Errorf("NearFill(6): expected left fill > 3 cols, got %q (left=%q)", got, leftPart)
	}
	if !strings.Contains(rightPart, "─") {
		t.Errorf("NearFill(6): expected fill on the right too, got %q", got)
	}
}

func TestDivider_NearFillZero(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewDivider().Label("Hi").LabelAlign(ink.JustifyStart).NearFill(0).Width(20).Render()

	// With NearFill(0) the label should be flush at the very start.
	if !strings.HasPrefix(got, " Hi ") {
		t.Errorf("NearFill(0): expected label flush at start, got %q", got)
	}
}

func TestDivider_NearFillZeroEnd(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewDivider().Label("Hi").LabelAlign(ink.JustifyEnd).NearFill(0).Width(20).Render()

	// With NearFill(0) the label should be flush at the very end.
	if !strings.HasSuffix(got, " Hi ") {
		t.Errorf("NearFill(0): expected label flush at end, got %q", got)
	}
}

func TestDivider_NearFillNegativeClamped(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Negative value should be clamped to 0 — same as NearFill(0).
	withNeg := cast.NewDivider().Label("Hi").LabelAlign(ink.JustifyStart).NearFill(-5).Width(20).Render()
	withZero := cast.NewDivider().Label("Hi").LabelAlign(ink.JustifyStart).NearFill(0).Width(20).Render()

	if withNeg != withZero {
		t.Errorf("NearFill(-5) should equal NearFill(0): neg=%q zero=%q", withNeg, withZero)
	}
}

func TestDivider_NearFillDefaultPreserved(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Default NearFill is 3 — explicit NearFill(3) should produce identical output.
	implicit := cast.NewDivider().Label("Hi").LabelAlign(ink.JustifyStart).Width(20).Render()
	explicit := cast.NewDivider().Label("Hi").LabelAlign(ink.JustifyStart).NearFill(3).Width(20).Render()

	if implicit != explicit {
		t.Errorf("default NearFill should equal NearFill(3): implicit=%q explicit=%q", implicit, explicit)
	}
}

func TestDivider_NearFillNoEffectOnCenter(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// NearFill has no effect when alignment is center.
	without := cast.NewDivider().Label("Hi").Width(20).Render()
	with := cast.NewDivider().Label("Hi").NearFill(10).Width(20).Render()

	if without != with {
		t.Errorf("NearFill should have no effect on JustifyCenter: without=%q with=%q", without, with)
	}
}

func TestDivider_NearFillImmutability(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	base := cast.NewDivider().Label("Hi").LabelAlign(ink.JustifyStart).Width(20)
	custom := base.NearFill(6)

	if base.Render() == custom.Render() {
		t.Errorf("NearFill setter mutated base divider: %q", base.Render())
	}
}

func TestDivider_Style(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	s := ink.New().WithForeground(ink.Muted)
	got := cast.NewDivider().Style(s).Width(10).Render()

	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI sequences in styled divider, got %q", got)
	}
	plain := ink.Strip(got)
	if plain != "──────────" {
		t.Errorf("expected 10 dashes (stripped), got %q", plain)
	}
}

func TestDivider_LabelStyle(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	ruleStyle := ink.New().WithForeground(ink.Muted)
	labelStyle := ink.New().WithForeground(ink.BrightWhite).WithBold(true)

	got := cast.NewDivider().
		Label("Section").
		Style(ruleStyle).
		LabelStyle(labelStyle).
		Width(30).
		Render()

	if !strings.Contains(got, "Section") {
		t.Errorf("expected label in output, got %q", got)
	}
	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI sequences, got %q", got)
	}
}

func TestDivider_NoColorMode(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	s := ink.New().WithForeground(ink.Danger)
	got := cast.NewDivider().Style(s).Width(10).Render()

	if strings.Contains(got, "\x1b[") {
		t.Errorf("expected no ANSI in NeverColorMode, got %q", got)
	}
}

func TestDivider_Immutability(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	base := cast.NewDivider().Width(10)
	withLabel := base.Label("Hi")

	// base should have no label
	if strings.Contains(base.Render(), "Hi") {
		t.Errorf("base divider was mutated to include label")
	}

	// withLabel should have the label
	if !strings.Contains(withLabel.Render(), "Hi") {
		t.Errorf("labelled divider missing label")
	}
}

func TestDivider_WidthNegative(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Negative width is clamped to 0, which falls back to default (80).
	got := cast.NewDivider().Width(-5).Render()
	w := visibleWidthTest(got)
	if w != 80 {
		t.Errorf("expected default width 80 for negative input, got %d", w)
	}
}

func TestDivider_LabelEqualsWidth(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// When the label (including padding) is wider than width, no fill is added
	// but the label is still rendered.
	got := cast.NewDivider().Label("Hello").Width(5).Render()
	if !strings.Contains(got, "Hello") {
		t.Errorf("expected label to appear even when wider than width, got %q", got)
	}
}

func TestDivider_EmptyChar(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// An empty Char should fall back to the default "─".
	got := cast.NewDivider().Char("").Width(5).Render()
	if got != "─────" {
		t.Errorf("expected default char fallback, got %q", got)
	}
}
