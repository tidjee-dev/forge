package cast_test

import (
	"strings"
	"testing"

	"github.com/tidjee-dev/forge/cast"
	"github.com/tidjee-dev/forge/ink"
)

func TestList_RenderEmpty(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewList().Render()
	if got != "" {
		t.Errorf("expected empty string for empty list, got %q", got)
	}
}

func TestList_SingleItem(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewList("only").Render()
	if got != "• only" {
		t.Errorf("expected \"• only\", got %q", got)
	}
}

func TestList_MultipleItems(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewList("a", "b", "c").Render()
	lines := strings.Split(got, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %q", len(lines), got)
	}
	if lines[0] != "• a" {
		t.Errorf("line 0: expected \"• a\", got %q", lines[0])
	}
	if lines[1] != "• b" {
		t.Errorf("line 1: expected \"• b\", got %q", lines[1])
	}
	if lines[2] != "• c" {
		t.Errorf("line 2: expected \"• c\", got %q", lines[2])
	}
}

func TestList_CustomBullet(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewList("one", "two").Bullet("*").Render()
	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "* one" {
		t.Errorf("expected \"* one\", got %q", lines[0])
	}
	if lines[1] != "* two" {
		t.Errorf("expected \"* two\", got %q", lines[1])
	}
}

func TestList_EmptyBullet(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Empty bullet: no prefix, just the item text.
	got := cast.NewList("alpha", "beta").Bullet("").Render()
	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "alpha" {
		t.Errorf("expected \"alpha\", got %q", lines[0])
	}
	if lines[1] != "beta" {
		t.Errorf("expected \"beta\", got %q", lines[1])
	}
}

func TestList_Numbered(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewList("Buy carrots", "Buy celery", "Buy kohlrabi").Numbered().Render()
	lines := strings.Split(got, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %q", len(lines), got)
	}
	if lines[0] != "1. Buy carrots" {
		t.Errorf("line 0: expected \"1. Buy carrots\", got %q", lines[0])
	}
	if lines[1] != "2. Buy celery" {
		t.Errorf("line 1: expected \"2. Buy celery\", got %q", lines[1])
	}
	if lines[2] != "3. Buy kohlrabi" {
		t.Errorf("line 2: expected \"3. Buy kohlrabi\", got %q", lines[2])
	}
}

func TestList_NumberedSingleItem(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewList("only").Numbered().Render()
	if got != "1. only" {
		t.Errorf("expected \"1. only\", got %q", got)
	}
}

func TestList_AddItem(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewList("first").AddItem("second").AddItem("third").Render()
	lines := strings.Split(got, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %q", len(lines), got)
	}
	if lines[2] != "• third" {
		t.Errorf("expected \"• third\", got %q", lines[2])
	}
}

func TestList_Indent(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewList("item").Indent(2).Render()
	if got != "  • item" {
		t.Errorf("expected \"  • item\", got %q", got)
	}
}

func TestList_IndentMultiLine(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewList("a", "b").Indent(4).Render()
	lines := strings.Split(got, "\n")
	for i, line := range lines {
		if !strings.HasPrefix(line, "    ") {
			t.Errorf("line %d: expected 4-space indent, got %q", i, line)
		}
	}
}

func TestList_IndentNegativeClamped(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Negative indent is clamped to 0.
	got := cast.NewList("item").Indent(-3).Render()
	if got != "• item" {
		t.Errorf("expected no indent, got %q", got)
	}
}

func TestList_ItemStyle(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	s := ink.New().WithForeground(ink.White)
	got := cast.NewList("hello").ItemStyle(s).Render()

	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI sequences in styled list, got %q", got)
	}
	plain := ink.Strip(got)
	if !strings.Contains(plain, "hello") {
		t.Errorf("expected item text in output, got %q", plain)
	}
}

func TestList_BulletStyle(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	bulletStyle := ink.New().WithForeground(ink.Muted)
	itemStyle := ink.New().WithForeground(ink.White)

	got := cast.NewList("item").
		BulletStyle(bulletStyle).
		ItemStyle(itemStyle).
		Render()

	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI sequences, got %q", got)
	}
	plain := ink.Strip(got)
	if !strings.Contains(plain, "item") {
		t.Errorf("expected item text in stripped output, got %q", plain)
	}
}

func TestList_BulletStyleFallsBackToItemStyle(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// When no BulletStyle is set, the bullet should use the ItemStyle.
	itemStyle := ink.New().WithForeground(ink.Cyan)
	got := cast.NewList("x").ItemStyle(itemStyle).Render()

	// Both the bullet and the item should be styled (i.e. ANSI appears).
	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI sequences, got %q", got)
	}
}

func TestList_NumberedThenBullet(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Bullet() after Numbered() should disable numbering.
	got := cast.NewList("a", "b").Numbered().Bullet("-").Render()
	lines := strings.Split(got, "\n")
	if lines[0] != "- a" {
		t.Errorf("expected \"- a\", got %q", lines[0])
	}
}

func TestList_NoColorMode(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	s := ink.New().WithForeground(ink.Danger)
	got := cast.NewList("err").ItemStyle(s).Render()

	if strings.Contains(got, "\x1b[") {
		t.Errorf("expected no ANSI in NeverColorMode, got %q", got)
	}
	if got != "• err" {
		t.Errorf("expected \"• err\", got %q", got)
	}
}

func TestList_Immutability(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	base := cast.NewList("first")
	extended := base.AddItem("second")

	baseLines := strings.Split(base.Render(), "\n")
	if len(baseLines) != 1 {
		t.Errorf("base list was mutated, got %d lines", len(baseLines))
	}

	extLines := strings.Split(extended.Render(), "\n")
	if len(extLines) != 2 {
		t.Errorf("extended list should have 2 lines, got %d", len(extLines))
	}
}

func TestList_NumberedImmutability(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	base := cast.NewList("a", "b")
	numbered := base.Numbered()

	// base should still use bullet
	if strings.Contains(base.Render(), "1.") {
		t.Errorf("base list was mutated to numbered mode")
	}
	// numbered should use numbers
	if !strings.Contains(numbered.Render(), "1.") {
		t.Errorf("numbered list missing \"1.\"")
	}
}

func TestList_IndentWithNumbered(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewList("x", "y").Numbered().Indent(3).Render()
	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "   1. x" {
		t.Errorf("expected \"   1. x\", got %q", lines[0])
	}
	if lines[1] != "   2. y" {
		t.Errorf("expected \"   2. y\", got %q", lines[1])
	}
}
