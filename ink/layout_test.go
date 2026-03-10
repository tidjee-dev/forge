package ink

import (
	"testing"
)

// ---------------------------------------------------------------------------
// TrackSize constructors
// ---------------------------------------------------------------------------

func TestTrackSize_FixedTrack(t *testing.T) {
	tr := FixedTrack(20)
	if tr.Kind() != TrackFixed {
		t.Errorf("FixedTrack kind = %v, want TrackFixed", tr.Kind())
	}
	if tr.Value() != 20 {
		t.Errorf("FixedTrack value = %d, want 20", tr.Value())
	}
	if tr.IsZero() {
		t.Error("FixedTrack(20).IsZero() = true, want false")
	}
}

func TestTrackSize_FixedTrack_Zero(t *testing.T) {
	tr := FixedTrack(0)
	if !tr.IsZero() {
		t.Error("FixedTrack(0).IsZero() = false, want true")
	}
}

func TestTrackSize_FrTrack(t *testing.T) {
	tr := FrTrack(3)
	if tr.Kind() != TrackFr {
		t.Errorf("FrTrack kind = %v, want TrackFr", tr.Kind())
	}
	if tr.Value() != 3 {
		t.Errorf("FrTrack value = %d, want 3", tr.Value())
	}
}

func TestTrackSize_FrTrack_ClampsBelowOne(t *testing.T) {
	tr := FrTrack(0)
	if tr.Value() != 1 {
		t.Errorf("FrTrack(0) value = %d, want 1 (clamped)", tr.Value())
	}
	tr = FrTrack(-5)
	if tr.Value() != 1 {
		t.Errorf("FrTrack(-5) value = %d, want 1 (clamped)", tr.Value())
	}
}

func TestTrackSize_AutoTrack(t *testing.T) {
	tr := AutoTrack()
	if tr.Kind() != TrackAuto {
		t.Errorf("AutoTrack kind = %v, want TrackAuto", tr.Kind())
	}
	if tr.IsZero() {
		t.Error("AutoTrack().IsZero() = true, want false")
	}
}

func TestTrackSize_MinContentTrack(t *testing.T) {
	tr := MinContentTrack()
	if tr.Kind() != TrackMinContent {
		t.Errorf("MinContentTrack kind = %v, want TrackMinContent", tr.Kind())
	}
}

func TestTrackSize_MaxContentTrack(t *testing.T) {
	tr := MaxContentTrack()
	if tr.Kind() != TrackMaxContent {
		t.Errorf("MaxContentTrack kind = %v, want TrackMaxContent", tr.Kind())
	}
}

// ---------------------------------------------------------------------------
// NewLayout — zero value
// ---------------------------------------------------------------------------

func TestNewLayout_ZeroValue(t *testing.T) {
	l := NewLayout()
	if !l.IsZero() {
		t.Error("NewLayout().IsZero() = false, want true")
	}
}

func TestLayout_ZeroValue_Defaults(t *testing.T) {
	l := NewLayout()

	if l.Direction() != DirectionColumn {
		t.Errorf("default Direction = %v, want DirectionColumn", l.Direction())
	}
	if l.JustifyContent() != JustifyStart {
		t.Errorf("default JustifyContent = %v, want JustifyStart", l.JustifyContent())
	}
	if l.AlignItems() != AlignItemsStart {
		t.Errorf("default AlignItems = %v, want AlignItemsStart", l.AlignItems())
	}
	if l.AlignSelf() != AlignSelfAuto {
		t.Errorf("default AlignSelf = %v, want AlignSelfAuto", l.AlignSelf())
	}
	if l.Wrap() {
		t.Error("default Wrap = true, want false")
	}
	if l.Grow() {
		t.Error("default Grow = true, want false")
	}
	if l.Shrink() {
		t.Error("default Shrink = true, want false")
	}
	if l.Basis() != 0 {
		t.Errorf("default Basis = %d, want 0", l.Basis())
	}
	if len(l.Columns()) != 0 {
		t.Errorf("default Columns len = %d, want 0", len(l.Columns()))
	}
	if len(l.Rows()) != 0 {
		t.Errorf("default Rows len = %d, want 0", len(l.Rows()))
	}
	if l.AutoFlow() != AutoFlowRow {
		t.Errorf("default AutoFlow = %v, want AutoFlowRow", l.AutoFlow())
	}
	if l.ColSpan() != 1 {
		t.Errorf("default ColSpan = %d, want 1", l.ColSpan())
	}
	if l.RowSpan() != 1 {
		t.Errorf("default RowSpan = %d, want 1", l.RowSpan())
	}
	if l.ColumnGap() != 0 {
		t.Errorf("default ColumnGap = %d, want 0", l.ColumnGap())
	}
	if l.RowGap() != 0 {
		t.Errorf("default RowGap = %d, want 0", l.RowGap())
	}
	if l.Padding() != (Edges{}) {
		t.Errorf("default Padding = %v, want zero Edges", l.Padding())
	}
	if l.MinWidth() != 0 {
		t.Errorf("default MinWidth = %d, want 0", l.MinWidth())
	}
	if l.MaxWidth() != 0 {
		t.Errorf("default MaxWidth = %d, want 0", l.MaxWidth())
	}
	if l.MinHeight() != 0 {
		t.Errorf("default MinHeight = %d, want 0", l.MinHeight())
	}
	if l.MaxHeight() != 0 {
		t.Errorf("default MaxHeight = %d, want 0", l.MaxHeight())
	}
	if l.Overflow() != OverflowClip {
		t.Errorf("default Overflow = %v, want OverflowClip", l.Overflow())
	}
}

// ---------------------------------------------------------------------------
// Immutability — setters return a new copy
// ---------------------------------------------------------------------------

func TestLayout_Immutability(t *testing.T) {
	original := NewLayout()
	modified := original.WithDirection(DirectionRow)

	if original.Direction() != DirectionColumn {
		t.Error("WithDirection mutated the original Layout")
	}
	if modified.Direction() != DirectionRow {
		t.Error("WithDirection did not set the new value")
	}
}

// ---------------------------------------------------------------------------
// Flex setters
// ---------------------------------------------------------------------------

func TestLayout_WithDirection(t *testing.T) {
	tests := []struct {
		name  string
		input Direction
	}{
		{"Column", DirectionColumn},
		{"Row", DirectionRow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLayout().WithDirection(tt.input)
			if l.Direction() != tt.input {
				t.Errorf("Direction() = %v, want %v", l.Direction(), tt.input)
			}
		})
	}
}

func TestLayout_WithJustifyContent(t *testing.T) {
	values := []JustifyContent{
		JustifyStart, JustifyCenter, JustifyEnd,
		JustifySpaceBetween, JustifySpaceAround, JustifySpaceEvenly,
	}
	for _, v := range values {
		l := NewLayout().WithJustifyContent(v)
		if l.JustifyContent() != v {
			t.Errorf("WithJustifyContent(%v): got %v", v, l.JustifyContent())
		}
	}
}

func TestLayout_WithAlignItems(t *testing.T) {
	values := []AlignItems{AlignItemsStart, AlignItemsCenter, AlignItemsEnd, AlignItemsStretch}
	for _, v := range values {
		l := NewLayout().WithAlignItems(v)
		if l.AlignItems() != v {
			t.Errorf("WithAlignItems(%v): got %v", v, l.AlignItems())
		}
	}
}

func TestLayout_WithAlignSelf(t *testing.T) {
	values := []AlignSelf{AlignSelfAuto, AlignSelfStart, AlignSelfCenter, AlignSelfEnd, AlignSelfStretch}
	for _, v := range values {
		l := NewLayout().WithAlignSelf(v)
		if l.AlignSelf() != v {
			t.Errorf("WithAlignSelf(%v): got %v", v, l.AlignSelf())
		}
	}
}

func TestLayout_WithWrap(t *testing.T) {
	if !NewLayout().WithWrap(true).Wrap() {
		t.Error("WithWrap(true): Wrap() = false")
	}
	if NewLayout().WithWrap(false).Wrap() {
		t.Error("WithWrap(false): Wrap() = true")
	}
}

func TestLayout_WithGrow(t *testing.T) {
	if !NewLayout().WithGrow(true).Grow() {
		t.Error("WithGrow(true): Grow() = false")
	}
	if NewLayout().WithGrow(false).Grow() {
		t.Error("WithGrow(false): Grow() = true")
	}
}

func TestLayout_WithShrink(t *testing.T) {
	if !NewLayout().WithShrink(true).Shrink() {
		t.Error("WithShrink(true): Shrink() = false")
	}
	if NewLayout().WithShrink(false).Shrink() {
		t.Error("WithShrink(false): Shrink() = true")
	}
}

func TestLayout_WithBasis(t *testing.T) {
	l := NewLayout().WithBasis(40)
	if l.Basis() != 40 {
		t.Errorf("Basis() = %d, want 40", l.Basis())
	}
}

func TestLayout_WithBasis_ClampsNegative(t *testing.T) {
	l := NewLayout().WithBasis(-10)
	if l.Basis() != 0 {
		t.Errorf("WithBasis(-10): Basis() = %d, want 0", l.Basis())
	}
}

// ---------------------------------------------------------------------------
// Grid setters
// ---------------------------------------------------------------------------

func TestLayout_WithColumns(t *testing.T) {
	cols := []TrackSize{FrTrack(1), FrTrack(2), FixedTrack(10)}
	l := NewLayout().WithColumns(cols...)
	got := l.Columns()
	if len(got) != len(cols) {
		t.Fatalf("Columns() len = %d, want %d", len(got), len(cols))
	}
	for i, c := range cols {
		if got[i] != c {
			t.Errorf("Columns()[%d] = %v, want %v", i, got[i], c)
		}
	}
}

func TestLayout_WithColumns_Empty(t *testing.T) {
	l := NewLayout().WithColumns()
	if len(l.Columns()) != 0 {
		t.Errorf("WithColumns(): len = %d, want 0", len(l.Columns()))
	}
}

func TestLayout_WithRows(t *testing.T) {
	rows := []TrackSize{AutoTrack(), FixedTrack(5)}
	l := NewLayout().WithRows(rows...)
	got := l.Rows()
	if len(got) != len(rows) {
		t.Fatalf("Rows() len = %d, want %d", len(got), len(rows))
	}
	for i, r := range rows {
		if got[i] != r {
			t.Errorf("Rows()[%d] = %v, want %v", i, got[i], r)
		}
	}
}

func TestLayout_WithAutoFlow(t *testing.T) {
	l := NewLayout().WithAutoFlow(AutoFlowColumn)
	if l.AutoFlow() != AutoFlowColumn {
		t.Errorf("AutoFlow() = %v, want AutoFlowColumn", l.AutoFlow())
	}
	l = l.WithAutoFlow(AutoFlowRow)
	if l.AutoFlow() != AutoFlowRow {
		t.Errorf("AutoFlow() = %v, want AutoFlowRow", l.AutoFlow())
	}
}

func TestLayout_WithColSpan(t *testing.T) {
	l := NewLayout().WithColSpan(3)
	if l.ColSpan() != 3 {
		t.Errorf("ColSpan() = %d, want 3", l.ColSpan())
	}
}

func TestLayout_WithColSpan_ClampsToOne(t *testing.T) {
	tests := []struct{ input, want int }{{0, 1}, {-5, 1}}
	for _, tt := range tests {
		l := NewLayout().WithColSpan(tt.input)
		if l.ColSpan() != tt.want {
			t.Errorf("WithColSpan(%d): ColSpan() = %d, want %d", tt.input, l.ColSpan(), tt.want)
		}
	}
}

func TestLayout_WithRowSpan(t *testing.T) {
	l := NewLayout().WithRowSpan(2)
	if l.RowSpan() != 2 {
		t.Errorf("RowSpan() = %d, want 2", l.RowSpan())
	}
}

func TestLayout_WithRowSpan_ClampsToOne(t *testing.T) {
	tests := []struct{ input, want int }{{0, 1}, {-3, 1}}
	for _, tt := range tests {
		l := NewLayout().WithRowSpan(tt.input)
		if l.RowSpan() != tt.want {
			t.Errorf("WithRowSpan(%d): RowSpan() = %d, want %d", tt.input, l.RowSpan(), tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Gap setters
// ---------------------------------------------------------------------------

func TestLayout_WithGap(t *testing.T) {
	l := NewLayout().WithGap(2)
	if l.ColumnGap() != 2 {
		t.Errorf("ColumnGap() = %d, want 2", l.ColumnGap())
	}
	if l.RowGap() != 2 {
		t.Errorf("RowGap() = %d, want 2", l.RowGap())
	}
}

func TestLayout_WithGap_ClampsNegative(t *testing.T) {
	l := NewLayout().WithGap(-1)
	if l.ColumnGap() != 0 || l.RowGap() != 0 {
		t.Errorf("WithGap(-1): ColumnGap=%d RowGap=%d, want both 0", l.ColumnGap(), l.RowGap())
	}
}

func TestLayout_WithColumnGap(t *testing.T) {
	l := NewLayout().WithColumnGap(4)
	if l.ColumnGap() != 4 {
		t.Errorf("ColumnGap() = %d, want 4", l.ColumnGap())
	}
	if l.RowGap() != 0 {
		t.Errorf("RowGap() = %d, want 0 (untouched)", l.RowGap())
	}
}

func TestLayout_WithRowGap(t *testing.T) {
	l := NewLayout().WithRowGap(3)
	if l.RowGap() != 3 {
		t.Errorf("RowGap() = %d, want 3", l.RowGap())
	}
	if l.ColumnGap() != 0 {
		t.Errorf("ColumnGap() = %d, want 0 (untouched)", l.ColumnGap())
	}
}

func TestLayout_WithColumnGap_ClampsNegative(t *testing.T) {
	l := NewLayout().WithColumnGap(-2)
	if l.ColumnGap() != 0 {
		t.Errorf("WithColumnGap(-2): ColumnGap() = %d, want 0", l.ColumnGap())
	}
}

func TestLayout_WithRowGap_ClampsNegative(t *testing.T) {
	l := NewLayout().WithRowGap(-2)
	if l.RowGap() != 0 {
		t.Errorf("WithRowGap(-2): RowGap() = %d, want 0", l.RowGap())
	}
}

// ---------------------------------------------------------------------------
// Padding setters
// ---------------------------------------------------------------------------

func TestLayout_WithPadding(t *testing.T) {
	l := NewLayout().WithPadding(1, 2, 3, 4)
	p := l.Padding()
	if p.Top != 1 || p.Right != 2 || p.Bottom != 3 || p.Left != 4 {
		t.Errorf("Padding() = %v, want {1 2 3 4}", p)
	}
	if l.PaddingTop() != 1 {
		t.Errorf("PaddingTop() = %d, want 1", l.PaddingTop())
	}
	if l.PaddingRight() != 2 {
		t.Errorf("PaddingRight() = %d, want 2", l.PaddingRight())
	}
	if l.PaddingBottom() != 3 {
		t.Errorf("PaddingBottom() = %d, want 3", l.PaddingBottom())
	}
	if l.PaddingLeft() != 4 {
		t.Errorf("PaddingLeft() = %d, want 4", l.PaddingLeft())
	}
}

func TestLayout_WithPadding_ClampsNegative(t *testing.T) {
	l := NewLayout().WithPadding(-1, -2, -3, -4)
	p := l.Padding()
	if p.Top != 0 || p.Right != 0 || p.Bottom != 0 || p.Left != 0 {
		t.Errorf("WithPadding with negatives: Padding() = %v, want all 0", p)
	}
}

func TestLayout_WithPaddingX(t *testing.T) {
	l := NewLayout().WithPaddingX(3)
	if l.PaddingLeft() != 3 || l.PaddingRight() != 3 {
		t.Errorf("WithPaddingX(3): Left=%d Right=%d, want both 3", l.PaddingLeft(), l.PaddingRight())
	}
	if l.PaddingTop() != 0 || l.PaddingBottom() != 0 {
		t.Errorf("WithPaddingX(3): Top=%d Bottom=%d, want both 0", l.PaddingTop(), l.PaddingBottom())
	}
}

func TestLayout_WithPaddingY(t *testing.T) {
	l := NewLayout().WithPaddingY(2)
	if l.PaddingTop() != 2 || l.PaddingBottom() != 2 {
		t.Errorf("WithPaddingY(2): Top=%d Bottom=%d, want both 2", l.PaddingTop(), l.PaddingBottom())
	}
	if l.PaddingLeft() != 0 || l.PaddingRight() != 0 {
		t.Errorf("WithPaddingY(2): Left=%d Right=%d, want both 0", l.PaddingLeft(), l.PaddingRight())
	}
}

func TestLayout_WithUniformPadding(t *testing.T) {
	l := NewLayout().WithUniformPadding(4)
	p := l.Padding()
	if p.Top != 4 || p.Right != 4 || p.Bottom != 4 || p.Left != 4 {
		t.Errorf("WithUniformPadding(4): Padding() = %v, want all 4", p)
	}
}

func TestLayout_WithPaddingTop(t *testing.T) {
	l := NewLayout().WithPaddingTop(5)
	if l.PaddingTop() != 5 {
		t.Errorf("PaddingTop() = %d, want 5", l.PaddingTop())
	}
	if l.PaddingRight() != 0 || l.PaddingBottom() != 0 || l.PaddingLeft() != 0 {
		t.Error("WithPaddingTop affected other padding sides")
	}
}

func TestLayout_WithPaddingRight(t *testing.T) {
	l := NewLayout().WithPaddingRight(6)
	if l.PaddingRight() != 6 {
		t.Errorf("PaddingRight() = %d, want 6", l.PaddingRight())
	}
	if l.PaddingTop() != 0 || l.PaddingBottom() != 0 || l.PaddingLeft() != 0 {
		t.Error("WithPaddingRight affected other padding sides")
	}
}

func TestLayout_WithPaddingBottom(t *testing.T) {
	l := NewLayout().WithPaddingBottom(7)
	if l.PaddingBottom() != 7 {
		t.Errorf("PaddingBottom() = %d, want 7", l.PaddingBottom())
	}
	if l.PaddingTop() != 0 || l.PaddingRight() != 0 || l.PaddingLeft() != 0 {
		t.Error("WithPaddingBottom affected other padding sides")
	}
}

func TestLayout_WithPaddingLeft(t *testing.T) {
	l := NewLayout().WithPaddingLeft(8)
	if l.PaddingLeft() != 8 {
		t.Errorf("PaddingLeft() = %d, want 8", l.PaddingLeft())
	}
	if l.PaddingTop() != 0 || l.PaddingRight() != 0 || l.PaddingBottom() != 0 {
		t.Error("WithPaddingLeft affected other padding sides")
	}
}

func TestLayout_IndividualPadding_ClampsNegative(t *testing.T) {
	tests := []struct {
		name string
		l    Layout
		want int
		got  func(Layout) int
	}{
		{"Top", NewLayout().WithPaddingTop(-1), 0, Layout.PaddingTop},
		{"Right", NewLayout().WithPaddingRight(-1), 0, Layout.PaddingRight},
		{"Bottom", NewLayout().WithPaddingBottom(-1), 0, Layout.PaddingBottom},
		{"Left", NewLayout().WithPaddingLeft(-1), 0, Layout.PaddingLeft},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got(tt.l) != tt.want {
				t.Errorf("WithPadding%s(-1) = %d, want 0", tt.name, tt.got(tt.l))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Size constraint setters
// ---------------------------------------------------------------------------

func TestLayout_WithMinWidth(t *testing.T) {
	l := NewLayout().WithMinWidth(20)
	if l.MinWidth() != 20 {
		t.Errorf("MinWidth() = %d, want 20", l.MinWidth())
	}
}

func TestLayout_WithMaxWidth(t *testing.T) {
	l := NewLayout().WithMaxWidth(80)
	if l.MaxWidth() != 80 {
		t.Errorf("MaxWidth() = %d, want 80", l.MaxWidth())
	}
}

func TestLayout_WithMinHeight(t *testing.T) {
	l := NewLayout().WithMinHeight(5)
	if l.MinHeight() != 5 {
		t.Errorf("MinHeight() = %d, want 5", l.MinHeight())
	}
}

func TestLayout_WithMaxHeight(t *testing.T) {
	l := NewLayout().WithMaxHeight(24)
	if l.MaxHeight() != 24 {
		t.Errorf("MaxHeight() = %d, want 24", l.MaxHeight())
	}
}

func TestLayout_SizeConstraints_ClampsNegative(t *testing.T) {
	l := NewLayout().
		WithMinWidth(-1).
		WithMaxWidth(-1).
		WithMinHeight(-1).
		WithMaxHeight(-1)

	if l.MinWidth() != 0 {
		t.Errorf("WithMinWidth(-1) = %d, want 0", l.MinWidth())
	}
	if l.MaxWidth() != 0 {
		t.Errorf("WithMaxWidth(-1) = %d, want 0", l.MaxWidth())
	}
	if l.MinHeight() != 0 {
		t.Errorf("WithMinHeight(-1) = %d, want 0", l.MinHeight())
	}
	if l.MaxHeight() != 0 {
		t.Errorf("WithMaxHeight(-1) = %d, want 0", l.MaxHeight())
	}
}

func TestLayout_WithWidth(t *testing.T) {
	l := NewLayout().WithWidth(40)
	if l.MinWidth() != 40 {
		t.Errorf("WithWidth(40): MinWidth() = %d, want 40", l.MinWidth())
	}
	if l.MaxWidth() != 40 {
		t.Errorf("WithWidth(40): MaxWidth() = %d, want 40", l.MaxWidth())
	}
}

func TestLayout_WithHeight(t *testing.T) {
	l := NewLayout().WithHeight(20)
	if l.MinHeight() != 20 {
		t.Errorf("WithHeight(20): MinHeight() = %d, want 20", l.MinHeight())
	}
	if l.MaxHeight() != 20 {
		t.Errorf("WithHeight(20): MaxHeight() = %d, want 20", l.MaxHeight())
	}
}

func TestLayout_WithSize(t *testing.T) {
	l := NewLayout().WithSize(80, 24)
	if l.MinWidth() != 80 || l.MaxWidth() != 80 {
		t.Errorf("WithSize(80,24): width min=%d max=%d, want both 80", l.MinWidth(), l.MaxWidth())
	}
	if l.MinHeight() != 24 || l.MaxHeight() != 24 {
		t.Errorf("WithSize(80,24): height min=%d max=%d, want both 24", l.MinHeight(), l.MaxHeight())
	}
}

func TestLayout_WithSize_ClampsNegative(t *testing.T) {
	l := NewLayout().WithSize(-10, -5)
	if l.MinWidth() != 0 || l.MaxWidth() != 0 {
		t.Errorf("WithSize(-10,-5): width min=%d max=%d, want both 0", l.MinWidth(), l.MaxWidth())
	}
	if l.MinHeight() != 0 || l.MaxHeight() != 0 {
		t.Errorf("WithSize(-10,-5): height min=%d max=%d, want both 0", l.MinHeight(), l.MaxHeight())
	}
}

// ---------------------------------------------------------------------------
// Overflow setter
// ---------------------------------------------------------------------------

func TestLayout_WithOverflow(t *testing.T) {
	values := []Overflow{OverflowClip, OverflowHidden, OverflowWrap, OverflowScroll}
	for _, v := range values {
		l := NewLayout().WithOverflow(v)
		if l.Overflow() != v {
			t.Errorf("WithOverflow(%v): Overflow() = %v", v, l.Overflow())
		}
	}
}

// ---------------------------------------------------------------------------
// IsZero — detects non-zero after each category of setter
// ---------------------------------------------------------------------------

func TestLayout_IsZero_FalseAfterFlexSetter(t *testing.T) {
	tests := []struct {
		name string
		l    Layout
	}{
		{"WithDirection(Row)", NewLayout().WithDirection(DirectionRow)},
		{"WithJustifyContent(Center)", NewLayout().WithJustifyContent(JustifyCenter)},
		{"WithAlignItems(Center)", NewLayout().WithAlignItems(AlignItemsCenter)},
		{"WithAlignSelf(Center)", NewLayout().WithAlignSelf(AlignSelfCenter)},
		{"WithWrap(true)", NewLayout().WithWrap(true)},
		{"WithGrow(true)", NewLayout().WithGrow(true)},
		{"WithShrink(true)", NewLayout().WithShrink(true)},
		{"WithBasis(10)", NewLayout().WithBasis(10)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.l.IsZero() {
				t.Errorf("%s: IsZero() = true, want false", tt.name)
			}
		})
	}
}

func TestLayout_IsZero_FalseAfterGridSetter(t *testing.T) {
	tests := []struct {
		name string
		l    Layout
	}{
		{"WithColumns", NewLayout().WithColumns(FrTrack(1))},
		{"WithRows", NewLayout().WithRows(AutoTrack())},
		{"WithAutoFlow(Column)", NewLayout().WithAutoFlow(AutoFlowColumn)},
		{"WithColSpan(2)", NewLayout().WithColSpan(2)},
		{"WithRowSpan(2)", NewLayout().WithRowSpan(2)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.l.IsZero() {
				t.Errorf("%s: IsZero() = true, want false", tt.name)
			}
		})
	}
}

func TestLayout_IsZero_FalseAfterGapSetter(t *testing.T) {
	tests := []struct {
		name string
		l    Layout
	}{
		{"WithGap(1)", NewLayout().WithGap(1)},
		{"WithColumnGap(1)", NewLayout().WithColumnGap(1)},
		{"WithRowGap(1)", NewLayout().WithRowGap(1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.l.IsZero() {
				t.Errorf("%s: IsZero() = true, want false", tt.name)
			}
		})
	}
}

func TestLayout_IsZero_FalseAfterPaddingSetter(t *testing.T) {
	tests := []struct {
		name string
		l    Layout
	}{
		{"WithPadding", NewLayout().WithPadding(1, 0, 0, 0)},
		{"WithPaddingX", NewLayout().WithPaddingX(1)},
		{"WithPaddingY", NewLayout().WithPaddingY(1)},
		{"WithUniformPadding", NewLayout().WithUniformPadding(1)},
		{"WithPaddingTop", NewLayout().WithPaddingTop(1)},
		{"WithPaddingRight", NewLayout().WithPaddingRight(1)},
		{"WithPaddingBottom", NewLayout().WithPaddingBottom(1)},
		{"WithPaddingLeft", NewLayout().WithPaddingLeft(1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.l.IsZero() {
				t.Errorf("%s: IsZero() = true, want false", tt.name)
			}
		})
	}
}

func TestLayout_IsZero_FalseAfterSizeConstraint(t *testing.T) {
	tests := []struct {
		name string
		l    Layout
	}{
		{"WithMinWidth", NewLayout().WithMinWidth(10)},
		{"WithMaxWidth", NewLayout().WithMaxWidth(80)},
		{"WithMinHeight", NewLayout().WithMinHeight(5)},
		{"WithMaxHeight", NewLayout().WithMaxHeight(24)},
		{"WithWidth", NewLayout().WithWidth(40)},
		{"WithHeight", NewLayout().WithHeight(10)},
		{"WithSize", NewLayout().WithSize(80, 24)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.l.IsZero() {
				t.Errorf("%s: IsZero() = true, want false", tt.name)
			}
		})
	}
}

func TestLayout_IsZero_FalseAfterOverflow(t *testing.T) {
	l := NewLayout().WithOverflow(OverflowHidden)
	if l.IsZero() {
		t.Error("WithOverflow(OverflowHidden): IsZero() = true, want false")
	}
}

// ---------------------------------------------------------------------------
// Chaining — multiple setters in sequence
// ---------------------------------------------------------------------------

func TestLayout_Chaining(t *testing.T) {
	l := NewLayout().
		WithDirection(DirectionRow).
		WithJustifyContent(JustifySpaceBetween).
		WithAlignItems(AlignItemsStretch).
		WithGap(2).
		WithPadding(1, 2, 1, 2).
		WithMinWidth(20).
		WithMaxWidth(120).
		WithOverflow(OverflowWrap)

	if l.Direction() != DirectionRow {
		t.Errorf("Direction = %v, want DirectionRow", l.Direction())
	}
	if l.JustifyContent() != JustifySpaceBetween {
		t.Errorf("JustifyContent = %v, want JustifySpaceBetween", l.JustifyContent())
	}
	if l.AlignItems() != AlignItemsStretch {
		t.Errorf("AlignItems = %v, want AlignItemsStretch", l.AlignItems())
	}
	if l.ColumnGap() != 2 || l.RowGap() != 2 {
		t.Errorf("Gap: col=%d row=%d, want both 2", l.ColumnGap(), l.RowGap())
	}
	if p := l.Padding(); p.Top != 1 || p.Right != 2 || p.Bottom != 1 || p.Left != 2 {
		t.Errorf("Padding = %v, want {1 2 1 2}", p)
	}
	if l.MinWidth() != 20 || l.MaxWidth() != 120 {
		t.Errorf("Width: min=%d max=%d, want 20/120", l.MinWidth(), l.MaxWidth())
	}
	if l.Overflow() != OverflowWrap {
		t.Errorf("Overflow = %v, want OverflowWrap", l.Overflow())
	}
	if l.IsZero() {
		t.Error("IsZero() = true after chaining, want false")
	}
}

// ---------------------------------------------------------------------------
// Style.WithLayout / Style.GetLayout integration
// ---------------------------------------------------------------------------

func TestStyle_WithLayout_GetLayout(t *testing.T) {
	l := NewLayout().
		WithUniformPadding(2).
		WithJustifyContent(JustifyCenter).
		WithWidth(60)

	s := New().WithLayout(l)
	got := s.GetLayout()

	if got.Padding() != l.Padding() {
		t.Errorf("GetLayout().Padding() = %v, want %v", got.Padding(), l.Padding())
	}
	if got.JustifyContent() != JustifyCenter {
		t.Errorf("GetLayout().JustifyContent() = %v, want JustifyCenter", got.JustifyContent())
	}
	if got.MinWidth() != 60 || got.MaxWidth() != 60 {
		t.Errorf("GetLayout().Width: min=%d max=%d, want both 60", got.MinWidth(), got.MaxWidth())
	}
}

func TestStyle_WithLayout_IsImmutable(t *testing.T) {
	l := NewLayout().WithUniformPadding(1)
	s1 := New().WithLayout(l)
	s2 := s1.WithLayout(NewLayout().WithUniformPadding(5))

	if s1.GetLayout().PaddingTop() != 1 {
		t.Error("WithLayout mutated the original Style")
	}
	if s2.GetLayout().PaddingTop() != 5 {
		t.Error("WithLayout did not update the new Style")
	}
}

func TestStyle_GetLayout_ZeroWhenUnset(t *testing.T) {
	s := New()
	if !s.GetLayout().IsZero() {
		t.Error("GetLayout() on a style with no layout set: IsZero() = false, want true")
	}
}

// ---------------------------------------------------------------------------
// StyleOverride.WithLayout / Merge integration
// ---------------------------------------------------------------------------

func TestStyleOverride_WithLayout_Merge(t *testing.T) {
	base := New().WithLayout(NewLayout().WithUniformPadding(1).WithWidth(40))
	patch := Override().WithLayout(NewLayout().WithUniformPadding(3).WithWidth(80))

	result := base.Merge(patch)
	l := result.GetLayout()

	if l.PaddingTop() != 3 {
		t.Errorf("Merged layout PaddingTop = %d, want 3", l.PaddingTop())
	}
	if l.MinWidth() != 80 {
		t.Errorf("Merged layout MinWidth = %d, want 80", l.MinWidth())
	}
}

func TestStyleOverride_UnsetLayout_WithLayout_Merge(t *testing.T) {
	base := New().WithLayout(NewLayout().WithUniformPadding(2))
	patch := Override().WithLayout(NewLayout()) // explicitly clear layout

	result := base.Merge(patch)
	if !result.GetLayout().IsZero() {
		t.Error("Merge with empty Layout: result layout should be zero")
	}
}

func TestStyleOverride_NoLayout_Merge_PreservesBase(t *testing.T) {
	l := NewLayout().WithUniformPadding(4).WithWidth(60)
	base := New().WithLayout(l)
	patch := Override().Bold() // no layout in the override

	result := base.Merge(patch)
	if result.GetLayout().PaddingTop() != 4 {
		t.Errorf("Merge without layout override: PaddingTop = %d, want 4 (preserved from base)", result.GetLayout().PaddingTop())
	}
}
