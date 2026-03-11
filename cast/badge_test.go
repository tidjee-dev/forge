package cast_test

import (
	"strings"
	"testing"

	"github.com/tidjee-dev/forge/cast"
	"github.com/tidjee-dev/forge/ink"
)

func TestBadge_RenderPlain(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewBadge("OK").Render()
	if got != " OK " {
		t.Errorf("expected \" OK \", got %q", got)
	}
}

func TestBadge_RenderEmpty(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewBadge("").Render()
	// Empty label → "  " (two spaces: one each side)
	if got != "  " {
		t.Errorf("expected \"  \", got %q", got)
	}
}

func TestBadge_Style(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	s := ink.New().WithForeground(ink.White).WithBackground(ink.Blue).WithBold(true)
	got := cast.NewBadge("INFO").Style(s).Render()

	// Should contain the label surrounded by spaces.
	if !strings.Contains(ink.Strip(got), " INFO ") {
		t.Errorf("rendered badge missing padded label, got %q", got)
	}
	// Should contain ANSI sequences when color is enabled.
	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI sequences in output, got %q", got)
	}
}

func TestBadge_Success(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewBadge("OK").Success().Render()
	plain := ink.Strip(got)
	if plain != " OK " {
		t.Errorf("expected plain \" OK \", got %q", plain)
	}
	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI sequences for Success badge, got %q", got)
	}
}

func TestBadge_Warning(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewBadge("WARN").Warning().Render()
	plain := ink.Strip(got)
	if plain != " WARN " {
		t.Errorf("expected plain \" WARN \", got %q", plain)
	}
}

func TestBadge_Danger(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewBadge("ERR").Danger().Render()
	plain := ink.Strip(got)
	if plain != " ERR " {
		t.Errorf("expected plain \" ERR \", got %q", plain)
	}
}

func TestBadge_Info(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewBadge("INFO").Info().Render()
	plain := ink.Strip(got)
	if plain != " INFO " {
		t.Errorf("expected plain \" INFO \", got %q", plain)
	}
}

func TestBadge_Neutral(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewBadge("v1.2.3").Neutral().Render()
	plain := ink.Strip(got)
	if plain != " v1.2.3 " {
		t.Errorf("expected plain \" v1.2.3 \", got %q", plain)
	}
}

func TestBadge_NoColorMode(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Even with a rich style, no ANSI should be emitted.
	got := cast.NewBadge("TEST").Success().Render()
	if strings.Contains(got, "\x1b[") {
		t.Errorf("expected no ANSI sequences in NeverColorMode, got %q", got)
	}
	if got != " TEST " {
		t.Errorf("expected \" TEST \" in NeverColorMode, got %q", got)
	}
}

func TestBadge_Immutability(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	original := cast.NewBadge("BASE")
	modified := original.Style(ink.New().WithBold(true))

	if original.Render() != " BASE " {
		t.Errorf("original badge mutated after Style() call")
	}
	if modified.Render() != " BASE " {
		t.Errorf("modified badge label changed unexpectedly, got %q", modified.Render())
	}
}

func TestBadge_SemanticChainingReplacesStyle(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Applying a second semantic shorthand should replace the first.
	badge := cast.NewBadge("X").Success().Danger()
	plain := ink.Strip(badge.Render())
	if plain != " X " {
		t.Errorf("expected \" X \", got %q", plain)
	}
	// Both Success and Danger should produce ANSI output.
	if !strings.Contains(badge.Render(), "\x1b[") {
		t.Errorf("expected ANSI in chained semantic badge, got %q", badge.Render())
	}
}
