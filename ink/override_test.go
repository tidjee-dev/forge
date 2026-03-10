package ink

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Zero override is identity
// ---------------------------------------------------------------------------

func TestStyleOverride_Zero_IsIdentity(t *testing.T) {
	base := New().
		WithForeground(Red).
		WithBackground(Blue).
		WithBold(true).
		WithItalic(true).
		WithUnderline(true).
		WithDim(true).
		WithStrikethrough(true).
		WithInverse(true).
		WithBorder(BorderNormal()).
		WithBorderColor(Green).
		WithBorderSide(BorderSideTop | BorderSideBottom)

	result := base.Merge(Override())

	if result.GetForeground() != base.GetForeground() {
		t.Error("zero override changed fg")
	}
	if result.GetBackground() != base.GetBackground() {
		t.Error("zero override changed bg")
	}
	if result.IsBold() != base.IsBold() {
		t.Error("zero override changed bold")
	}
	if result.IsItalic() != base.IsItalic() {
		t.Error("zero override changed italic")
	}
	if result.IsUnderline() != base.IsUnderline() {
		t.Error("zero override changed underline")
	}
	if result.IsDim() != base.IsDim() {
		t.Error("zero override changed dim")
	}
	if result.IsStrikethrough() != base.IsStrikethrough() {
		t.Error("zero override changed strikethrough")
	}
	if result.IsInverse() != base.IsInverse() {
		t.Error("zero override changed inverse")
	}
	if result.GetBorderStyle() != base.GetBorderStyle() {
		t.Error("zero override changed border style")
	}
	if result.GetBorderColor() != base.GetBorderColor() {
		t.Error("zero override changed border color")
	}

	if result.GetBorderSide() != base.GetBorderSide() {
		t.Error("zero override changed border sides")
	}
}

// ---------------------------------------------------------------------------
// Foreground / Background
// ---------------------------------------------------------------------------

func TestStyleOverride_Foreground_ReplacesBase(t *testing.T) {
	base := New().WithForeground(Red)
	result := base.Merge(Override().Foreground(Blue))
	if result.GetForeground() != Blue {
		t.Errorf("fg after override = %v, want %v", result.GetForeground(), Blue)
	}
}

func TestStyleOverride_Foreground_AbsentPreservesBase(t *testing.T) {
	base := New().WithForeground(Red)
	result := base.Merge(Override().Bold())
	if result.GetForeground() != Red {
		t.Errorf("fg preserved = %v, want %v", result.GetForeground(), Red)
	}
}

func TestStyleOverride_Background_ReplacesBase(t *testing.T) {
	base := New().WithBackground(Red)
	result := base.Merge(Override().Background(Blue))
	if result.GetBackground() != Blue {
		t.Errorf("bg after override = %v, want %v", result.GetBackground(), Blue)
	}
}

func TestStyleOverride_Background_AbsentPreservesBase(t *testing.T) {
	base := New().WithBackground(Green)
	result := base.Merge(Override().Bold())
	if result.GetBackground() != Green {
		t.Errorf("bg preserved = %v, want %v", result.GetBackground(), Green)
	}
}

// ---------------------------------------------------------------------------
// Bold
// ---------------------------------------------------------------------------

func TestStyleOverride_Bold_EnablesOnBase(t *testing.T) {
	base := New()
	result := base.Merge(Override().Bold())
	if !result.IsBold() {
		t.Error("Bold() override did not enable bold on base")
	}
}

func TestStyleOverride_NoBold_DisablesOnBase(t *testing.T) {
	base := New().WithBold(true)
	result := base.Merge(Override().NoBold())
	if result.IsBold() {
		t.Error("NoBold() override did not disable bold on base")
	}
}

func TestStyleOverride_Bold_Absent_PreservesBase(t *testing.T) {
	base := New().WithBold(true)
	result := base.Merge(Override().Italic())
	if !result.IsBold() {
		t.Error("absent bold override changed base bold state")
	}
}

// ---------------------------------------------------------------------------
// Italic
// ---------------------------------------------------------------------------

func TestStyleOverride_Italic_EnablesOnBase(t *testing.T) {
	base := New()
	result := base.Merge(Override().Italic())
	if !result.IsItalic() {
		t.Error("Italic() override did not enable italic")
	}
}

func TestStyleOverride_NoItalic_DisablesOnBase(t *testing.T) {
	base := New().WithItalic(true)
	result := base.Merge(Override().NoItalic())
	if result.IsItalic() {
		t.Error("NoItalic() override did not disable italic")
	}
}

func TestStyleOverride_Italic_Absent_PreservesBase(t *testing.T) {
	base := New().WithItalic(true)
	result := base.Merge(Override().Bold())
	if !result.IsItalic() {
		t.Error("absent italic override changed base italic state")
	}
}

// ---------------------------------------------------------------------------
// Underline
// ---------------------------------------------------------------------------

func TestStyleOverride_Underline_EnablesOnBase(t *testing.T) {
	base := New()
	result := base.Merge(Override().Underline())
	if !result.IsUnderline() {
		t.Error("Underline() override did not enable underline")
	}
}

func TestStyleOverride_NoUnderline_DisablesOnBase(t *testing.T) {
	base := New().WithUnderline(true)
	result := base.Merge(Override().NoUnderline())
	if result.IsUnderline() {
		t.Error("NoUnderline() override did not disable underline")
	}
}

func TestStyleOverride_Underline_Absent_PreservesBase(t *testing.T) {
	base := New().WithUnderline(true)
	result := base.Merge(Override())
	if !result.IsUnderline() {
		t.Error("absent underline override changed base underline state")
	}
}

// ---------------------------------------------------------------------------
// Dim
// ---------------------------------------------------------------------------

func TestStyleOverride_Dim_EnablesOnBase(t *testing.T) {
	base := New()
	result := base.Merge(Override().Dim())
	if !result.IsDim() {
		t.Error("Dim() override did not enable dim")
	}
}

func TestStyleOverride_NoDim_DisablesOnBase(t *testing.T) {
	base := New().WithDim(true)
	result := base.Merge(Override().NoDim())
	if result.IsDim() {
		t.Error("NoDim() override did not disable dim")
	}
}

func TestStyleOverride_Dim_Absent_PreservesBase(t *testing.T) {
	base := New().WithDim(true)
	result := base.Merge(Override())
	if !result.IsDim() {
		t.Error("absent dim override changed base dim state")
	}
}

// ---------------------------------------------------------------------------
// Strikethrough
// ---------------------------------------------------------------------------

func TestStyleOverride_Strike_EnablesOnBase(t *testing.T) {
	base := New()
	result := base.Merge(Override().Strike())
	if !result.IsStrikethrough() {
		t.Error("Strike() override did not enable strikethrough")
	}
}

func TestStyleOverride_NoStrike_DisablesOnBase(t *testing.T) {
	base := New().WithStrikethrough(true)
	result := base.Merge(Override().NoStrike())
	if result.IsStrikethrough() {
		t.Error("NoStrike() override did not disable strikethrough")
	}
}

func TestStyleOverride_Strike_Absent_PreservesBase(t *testing.T) {
	base := New().WithStrikethrough(true)
	result := base.Merge(Override())
	if !result.IsStrikethrough() {
		t.Error("absent strike override changed base strikethrough state")
	}
}

// ---------------------------------------------------------------------------
// Inverse
// ---------------------------------------------------------------------------

func TestStyleOverride_Inverse_EnablesOnBase(t *testing.T) {
	base := New()
	result := base.Merge(Override().Inverse())
	if !result.IsInverse() {
		t.Error("Inverse() override did not enable inverse")
	}
}

func TestStyleOverride_NoInverse_DisablesOnBase(t *testing.T) {
	base := New().WithInverse(true)
	result := base.Merge(Override().NoInverse())
	if result.IsInverse() {
		t.Error("NoInverse() override did not disable inverse")
	}
}

func TestStyleOverride_Inverse_Absent_PreservesBase(t *testing.T) {
	base := New().WithInverse(true)
	result := base.Merge(Override())
	if !result.IsInverse() {
		t.Error("absent inverse override changed base inverse state")
	}
}

// ---------------------------------------------------------------------------
// Layout
// ---------------------------------------------------------------------------

func TestStyleOverride_WithLayout_ReplacesBase(t *testing.T) {
	baseLayout := NewLayout().WithUniformPadding(2)
	patchLayout := NewLayout().WithUniformPadding(5)

	base := New().WithLayout(baseLayout)
	result := base.Merge(Override().WithLayout(patchLayout))

	if result.GetLayout().PaddingTop() != 5 {
		t.Errorf("layout override: PaddingTop = %d, want 5", result.GetLayout().PaddingTop())
	}
}

func TestStyleOverride_WithLayout_Absent_PreservesBase(t *testing.T) {
	baseLayout := NewLayout().WithUniformPadding(3)
	base := New().WithLayout(baseLayout)
	result := base.Merge(Override().Bold())

	if result.GetLayout().PaddingTop() != 3 {
		t.Errorf("absent layout override changed base layout, PaddingTop = %d, want 3",
			result.GetLayout().PaddingTop())
	}
}

func TestStyleOverride_WithLayout_ZeroLayout_ClearsBase(t *testing.T) {
	base := New().WithLayout(NewLayout().WithUniformPadding(4))
	result := base.Merge(Override().WithLayout(NewLayout()))

	if !result.GetLayout().IsZero() {
		t.Error("WithLayout(NewLayout()) did not clear the base layout")
	}
}

// ---------------------------------------------------------------------------
// Border
// ---------------------------------------------------------------------------

func TestStyleOverride_WithBorder_ReplacesBase(t *testing.T) {
	base := New().WithBorder(BorderNormal())
	result := base.Merge(Override().WithBorder(BorderRounded()))

	if result.GetBorderStyle() != BorderRounded() {
		t.Error("WithBorder override did not replace base border style")
	}
}

func TestStyleOverride_WithBorder_Absent_PreservesBase(t *testing.T) {
	base := New().WithBorder(BorderThick())
	result := base.Merge(Override().Bold())

	if result.GetBorderStyle() != BorderThick() {
		t.Error("absent border override changed base border style")
	}
}

func TestStyleOverride_WithNoBorder_ClearsBase(t *testing.T) {
	base := New().WithBorder(BorderDouble())
	result := base.Merge(Override().WithNoBorder())

	if !result.GetBorderStyle().IsZero() {
		t.Error("WithNoBorder did not clear the border style")
	}
	if result.GetBorderSide() != BorderSideNone {
		t.Errorf("WithNoBorder left sides = %v, want BorderSideNone", result.GetBorderSide())
	}
}

// ---------------------------------------------------------------------------
// Border Color
// ---------------------------------------------------------------------------

func TestStyleOverride_WithNoBorder_AlsoClearsBorderColor(t *testing.T) {
	base := New().
		WithBorder(BorderDouble()).
		WithBorderColor(Red)
	result := base.Merge(Override().WithNoBorder())

	if !result.GetBorderStyle().IsZero() {
		t.Error("WithNoBorder did not clear the border style")
	}
	if result.GetBorderSide() != BorderSideNone {
		t.Errorf("WithNoBorder left sides = %v, want BorderSideNone", result.GetBorderSide())
	}
	if !result.GetBorderColor().IsZeroColor() {
		t.Errorf("WithNoBorder did not clear border color, got %v", result.GetBorderColor())
	}
}

func TestStyleOverride_WithBorderColor_ReplacesBase(t *testing.T) {
	base := New().WithBorderColor(Red)
	result := base.Merge(Override().WithBorderColor(Blue))

	if result.GetBorderColor() != Blue {
		t.Errorf("border color after override = %v, want %v", result.GetBorderColor(), Blue)
	}
}

func TestStyleOverride_WithBorderColor_Absent_PreservesBase(t *testing.T) {
	base := New().WithBorderColor(Green)
	result := base.Merge(Override().Bold())

	if result.GetBorderColor() != Green {
		t.Errorf("absent border color override changed base, got %v, want %v",
			result.GetBorderColor(), Green)
	}
}

// ---------------------------------------------------------------------------
// Border Sides
// ---------------------------------------------------------------------------

func TestStyleOverride_WithBorderSides_ReplacesBase(t *testing.T) {
	base := New().WithBorderSide(BorderSideAll)
	result := base.Merge(Override().WithBorderSides(BorderSideTop | BorderSideBottom))

	if result.GetBorderSide() != BorderSideTop|BorderSideBottom {
		t.Errorf("border sides after override = %v, want Top|Bottom", result.GetBorderSide())
	}
}

func TestStyleOverride_WithBorderSides_Absent_PreservesBase(t *testing.T) {
	base := New().WithBorderSide(BorderSideVertical)
	result := base.Merge(Override().Bold())

	if result.GetBorderSide() != BorderSideVertical {
		t.Errorf("absent border sides override changed base, got %v, want Vertical",
			result.GetBorderSide())
	}
}

// ---------------------------------------------------------------------------
// Compound: multiple overrides in one patch
// ---------------------------------------------------------------------------

func TestStyleOverride_Compound(t *testing.T) {
	base := New().
		WithBold(true).
		WithItalic(true).
		WithForeground(Red).
		WithBackground(Blue)

	patch := Override().
		NoBold().
		Foreground(Green).
		WithBorderSides(BorderSideAll)

	result := base.Merge(patch)

	// NoBold explicitly unsets bold.
	if result.IsBold() {
		t.Error("compound: NoBold did not disable bold")
	}
	// Italic was not in the patch — must be preserved.
	if !result.IsItalic() {
		t.Error("compound: absent italic override changed base italic")
	}
	// Fg replaced.
	if result.GetForeground() != Green {
		t.Errorf("compound: fg = %v, want Green", result.GetForeground())
	}
	// Bg not in patch — must be preserved.
	if result.GetBackground() != Blue {
		t.Errorf("compound: bg = %v, want Blue", result.GetBackground())
	}
	// Border sides set.
	if result.GetBorderSide() != BorderSideAll {
		t.Errorf("compound: border sides = %v, want All", result.GetBorderSide())
	}
}

// ---------------------------------------------------------------------------
// Receiver is never mutated
// ---------------------------------------------------------------------------

func TestStyleOverride_Merge_DoesNotMutateReceiver(t *testing.T) {
	base := New().WithBold(true).WithForeground(Red)
	_ = base.Merge(Override().NoBold().Foreground(Blue))

	if !base.IsBold() {
		t.Error("Merge mutated receiver: bold was unset on base")
	}
	if base.GetForeground() != Red {
		t.Errorf("Merge mutated receiver: fg changed to %v", base.GetForeground())
	}
}

// ---------------------------------------------------------------------------
// Override() constructor returns empty override
// ---------------------------------------------------------------------------

func TestOverride_Constructor_IsEmpty(t *testing.T) {
	o := Override()
	base := New().
		WithBold(true).
		WithForeground(Red).
		WithBackground(Blue)

	result := base.Merge(o)

	// Everything should be identical to base.
	if result.IsBold() != base.IsBold() {
		t.Error("empty Override changed bold")
	}
	if result.GetForeground() != base.GetForeground() {
		t.Error("empty Override changed fg")
	}
	if result.GetBackground() != base.GetBackground() {
		t.Error("empty Override changed bg")
	}
}
