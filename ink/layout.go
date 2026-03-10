package ink

// ---------------------------------------------------------------------------
// Layout — CSS flex/grid-inspired terminal layout primitives
//
// Layout describes how a block of text is sized, padded, and positioned
// within the terminal. It is intentionally separate from Style (which owns
// SGR attributes and borders) so that the two concerns can be composed
// independently:
//
//   l := ink.NewLayout().
//       Padding(1, 2).
//       JustifyContent(ink.JustifyCenter).
//       AlignItems(ink.AlignItemsCenter).
//       Width(40).
//       Height(10)
//
//   s := ink.New().WithBold(true).WithLayout(l)
//
// Flex concepts that map cleanly onto terminal blocks
// ───────────────────────────────────────────────────
//
//   CSS                     │ ink equivalent
//   ────────────────────────┼─────────────────────────────────────────────────
//   flex-direction          │ Direction  (Row | Column)
//   justify-content         │ JustifyContent  (how content fills the main axis)
//   align-items             │ AlignItems      (cross-axis placement)
//   gap                     │ Gap / ColumnGap / RowGap
//   padding                 │ Padding / PaddingX / PaddingY
//   min/max-width/height    │ MinWidth / MaxWidth / MinHeight / MaxHeight
//   overflow                │ Overflow  (Hidden / Wrap / Scroll)
//   flex-grow               │ Grow  (bool — block expands to fill available space)
//   flex-shrink             │ Shrink  (bool — block may compress below natural size)
//   flex-basis              │ Basis  (explicit starting size before grow/shrink)
//   align-self              │ AlignSelf (per-block cross-axis override)
//
// Grid concepts that map onto terminal blocks
// ───────────────────────────────────────────
//
//   CSS                     │ ink equivalent
//   ────────────────────────┼─────────────────────────────────────────────────
//   grid-template-columns   │ Columns  ([]TrackSize — explicit column tracks)
//   grid-template-rows      │ Rows     ([]TrackSize — explicit row tracks)
//   column-gap / row-gap    │ ColumnGap / RowGap
//   grid-column / grid-row  │ ColSpan / RowSpan (on the child layout)
//   grid-auto-flow          │ AutoFlow (Row | Column)
//
// Overflow note
// ─────────────
// Terminals are character grids, not pixel surfaces. "Scroll" overflow means
// the text is clipped at MaxHeight lines and the consumer (e.g. anvil/viewport)
// is responsible for scrolling. "Wrap" means lines that exceed MaxWidth are
// soft-wrapped. "Hidden" means both axes are hard-clipped with no ellipsis.
// "Clip" (the default) clips width with an ellipsis (…) and clips height hard.
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Direction
// ---------------------------------------------------------------------------

// Direction controls the main axis along which flex children are laid out.
// It is the terminal equivalent of CSS flex-direction.
type Direction uint8

const (
	// DirectionColumn stacks children vertically (top to bottom).
	// This is the default and matches how terminal output naturally flows.
	DirectionColumn Direction = iota

	// DirectionRow places children horizontally (left to right) on one line.
	// Each child occupies a contiguous column region.
	DirectionRow
)

// ---------------------------------------------------------------------------
// JustifyContent
// ---------------------------------------------------------------------------

// JustifyContent controls how content is distributed along the main axis
// (horizontal for Row, vertical for Column).
// It is the terminal equivalent of CSS justify-content.
type JustifyContent uint8

const (
	// JustifyStart packs content toward the start of the main axis (default).
	JustifyStart JustifyContent = iota

	// JustifyCenter centers content on the main axis.
	JustifyCenter

	// JustifyEnd packs content toward the end of the main axis.
	JustifyEnd

	// JustifySpaceBetween distributes children with equal space between them;
	// first and last children are flush with the edges.
	JustifySpaceBetween

	// JustifySpaceAround distributes children with equal space on each side;
	// half-size space appears at the edges.
	JustifySpaceAround

	// JustifySpaceEvenly distributes children so that all gaps — including
	// the two edge gaps — are exactly equal.
	JustifySpaceEvenly
)

// ---------------------------------------------------------------------------
// AlignItems
// ---------------------------------------------------------------------------

// AlignItems controls placement along the cross axis
// (vertical for Row, horizontal for Column).
// It is the terminal equivalent of CSS align-items.
type AlignItems uint8

const (
	// AlignItemsStart aligns children to the start of the cross axis (default).
	AlignItemsStart AlignItems = iota

	// AlignItemsCenter centers children on the cross axis.
	AlignItemsCenter

	// AlignItemsEnd aligns children to the end of the cross axis.
	AlignItemsEnd

	// AlignItemsStretch expands children to fill the full cross-axis extent.
	// For Column direction this means all children share the same width.
	// For Row direction this means all children share the same height.
	AlignItemsStretch
)

// ---------------------------------------------------------------------------
// AlignSelf
// ---------------------------------------------------------------------------

// AlignSelf is a per-child override for AlignItems.
// It is the terminal equivalent of CSS align-self.
type AlignSelf uint8

const (
	// AlignSelfAuto defers to the parent's AlignItems value (default).
	AlignSelfAuto AlignSelf = iota

	// AlignSelfStart aligns this child to the start of the cross axis.
	AlignSelfStart

	// AlignSelfCenter centers this child on the cross axis.
	AlignSelfCenter

	// AlignSelfEnd aligns this child to the end of the cross axis.
	AlignSelfEnd

	// AlignSelfStretch stretches this child to fill the cross axis.
	AlignSelfStretch
)

// ---------------------------------------------------------------------------
// Overflow
// ---------------------------------------------------------------------------

// Overflow controls what happens when content exceeds the available space.
// It is the terminal equivalent of CSS overflow.
type Overflow uint8

const (
	// OverflowClip truncates width with an ellipsis (…) and clips height
	// without any indicator. This is the default.
	OverflowClip Overflow = iota

	// OverflowHidden clips both axes hard with no ellipsis.
	OverflowHidden

	// OverflowWrap soft-wraps lines that exceed MaxWidth at word boundaries
	// (falling back to character boundaries when no word boundary is found).
	OverflowWrap

	// OverflowScroll clips at MaxHeight lines; the consumer is responsible
	// for scrolling (used by anvil/viewport).
	OverflowScroll
)

// ---------------------------------------------------------------------------
// AutoFlow
// ---------------------------------------------------------------------------

// AutoFlow controls the direction in which auto-placed grid items flow.
// It is the terminal equivalent of CSS grid-auto-flow.
type AutoFlow uint8

const (
	// AutoFlowRow fills each row before adding a new one (default).
	AutoFlowRow AutoFlow = iota

	// AutoFlowColumn fills each column before adding a new one.
	AutoFlowColumn
)

// ---------------------------------------------------------------------------
// TrackSize — grid column / row sizing
// ---------------------------------------------------------------------------

// TrackKind identifies how a grid track is sized.
type TrackKind uint8

const (
	// TrackFixed is a fixed-width track measured in terminal columns / rows.
	TrackFixed TrackKind = iota

	// TrackFr is a fractional track that shares remaining space proportionally
	// with other fr tracks. It is the terminal equivalent of CSS fr units.
	TrackFr

	// TrackAuto sizes the track to fit its content (widest cell in the column,
	// or tallest cell in the row).
	TrackAuto

	// TrackMinContent sizes the track to the minimum content size.
	TrackMinContent

	// TrackMaxContent sizes the track to the maximum content size.
	TrackMaxContent
)

// TrackSize describes a single column or row track in a grid layout.
// Build one with the convenience constructors below.
type TrackSize struct {
	kind  TrackKind
	value int // columns for TrackFixed; numerator for TrackFr; ignored otherwise
}

// FixedTrack returns a grid track with a fixed size of n terminal columns/rows.
func FixedTrack(n int) TrackSize { return TrackSize{kind: TrackFixed, value: n} }

// FrTrack returns a fractional grid track with the given weight.
// FrTrack(1) alongside FrTrack(2) produces a 1:2 split of remaining space.
func FrTrack(weight int) TrackSize {
	if weight < 1 {
		weight = 1
	}
	return TrackSize{kind: TrackFr, value: weight}
}

// AutoTrack returns a grid track that sizes itself to fit its content.
func AutoTrack() TrackSize { return TrackSize{kind: TrackAuto} }

// MinContentTrack returns a grid track sized to the minimum content width/height.
func MinContentTrack() TrackSize { return TrackSize{kind: TrackMinContent} }

// MaxContentTrack returns a grid track sized to the maximum content width/height.
func MaxContentTrack() TrackSize { return TrackSize{kind: TrackMaxContent} }

// Kind returns the TrackKind of this track.
func (t TrackSize) Kind() TrackKind { return t.kind }

// Value returns the numeric value associated with the track
// (column count for TrackFixed, weight for TrackFr, 0 otherwise).
func (t TrackSize) Value() int { return t.value }

// IsZero reports whether the TrackSize is the zero value (unset).
func (t TrackSize) IsZero() bool { return t.kind == TrackFixed && t.value == 0 }

// ---------------------------------------------------------------------------
// Padding
// ---------------------------------------------------------------------------

// Edges holds the four sides of a padding or gap specification, in CSS order:
// Top, Right, Bottom, Left.
type Edges struct {
	Top, Right, Bottom, Left int
}

// uniform returns an Edges with all four sides set to v.
func uniform(v int) Edges { return Edges{v, v, v, v} }

// xy returns an Edges with horizontal (left/right) set to x and
// vertical (top/bottom) set to y.
func xy(x, y int) Edges { return Edges{Top: y, Right: x, Bottom: y, Left: x} }

// ---------------------------------------------------------------------------
// Layout
// ---------------------------------------------------------------------------

// Layout describes the spatial properties of a rendered block. It is
// intentionally separated from Style (which owns SGR text attributes and
// borders) so the two can be composed independently.
//
// All fields are unexported. Use the fluent builder methods to construct a
// Layout, and the accessor methods to read values back.
//
// The zero value of Layout is valid and represents a no-op layout:
// no padding, no size constraints, no alignment overrides.
type Layout struct {
	// ── Flex ────────────────────────────────────────────────────────────────
	direction      Direction
	justifyContent JustifyContent
	alignItems     AlignItems
	alignSelf      AlignSelf
	wrap           bool // equivalent to flex-wrap: wrap
	grow           bool // equivalent to flex-grow: 1  (false = 0)
	shrink         bool // equivalent to flex-shrink: 1 (false = 0)
	basis          int  // explicit flex-basis in terminal columns (0 = auto)

	// ── Grid ─────────────────────────────────────────────────────────────────
	columns  []TrackSize // grid-template-columns
	rows     []TrackSize // grid-template-rows
	autoFlow AutoFlow
	colSpan  int // grid-column span (0 = 1, default)
	rowSpan  int // grid-row span    (0 = 1, default)

	// ── Gap ──────────────────────────────────────────────────────────────────
	columnGap int
	rowGap    int

	// ── Padding ───────────────────────────────────────────────────────────────
	padding Edges

	// ── Size constraints ─────────────────────────────────────────────────────
	minWidth  int // 0 = no minimum
	maxWidth  int // 0 = no maximum
	minHeight int // 0 = no minimum
	maxHeight int // 0 = no maximum

	// ── Overflow ─────────────────────────────────────────────────────────────
	overflow Overflow
}

// NewLayout returns a zero-value Layout ready for configuration.
// The zero value is already valid (no-op layout), but this constructor
// makes intent explicit.
func NewLayout() Layout {
	return Layout{}
}

// ---------------------------------------------------------------------------
// Flex setters
// ---------------------------------------------------------------------------

// WithDirection sets the main axis direction (Row or Column).
func (l Layout) WithDirection(d Direction) Layout {
	l.direction = d
	return l
}

// WithJustifyContent sets how content is distributed along the main axis.
func (l Layout) WithJustifyContent(j JustifyContent) Layout {
	l.justifyContent = j
	return l
}

// WithAlignItems sets the default cross-axis placement for children.
func (l Layout) WithAlignItems(a AlignItems) Layout {
	l.alignItems = a
	return l
}

// WithAlignSelf overrides the parent's AlignItems for this specific block.
func (l Layout) WithAlignSelf(a AlignSelf) Layout {
	l.alignSelf = a
	return l
}

// WithWrap enables soft-wrapping of children that overflow the main axis.
// Equivalent to CSS flex-wrap: wrap.
func (l Layout) WithWrap(w bool) Layout {
	l.wrap = w
	return l
}

// WithGrow sets whether this block may expand to fill available main-axis
// space. Equivalent to CSS flex-grow: 1 (true) / 0 (false).
func (l Layout) WithGrow(g bool) Layout {
	l.grow = g
	return l
}

// WithShrink sets whether this block may compress below its natural size
// when space is scarce. Equivalent to CSS flex-shrink: 1 (true) / 0 (false).
func (l Layout) WithShrink(s bool) Layout {
	l.shrink = s
	return l
}

// WithBasis sets the explicit starting width (in terminal columns) before
// grow/shrink is applied. 0 means "auto". Equivalent to CSS flex-basis.
func (l Layout) WithBasis(columns int) Layout {
	if columns < 0 {
		columns = 0
	}
	l.basis = columns
	return l
}

// ---------------------------------------------------------------------------
// Grid setters
// ---------------------------------------------------------------------------

// WithColumns sets the grid column track definitions.
// Example: two equal fractional columns and one fixed 10-column sidebar.
//
//	l.WithColumns(ink.FrTrack(1), ink.FrTrack(1), ink.FixedTrack(10))
func (l Layout) WithColumns(tracks ...TrackSize) Layout {
	l.columns = tracks
	return l
}

// WithRows sets the grid row track definitions.
func (l Layout) WithRows(tracks ...TrackSize) Layout {
	l.rows = tracks
	return l
}

// WithAutoFlow sets the direction for auto-placed grid items.
func (l Layout) WithAutoFlow(f AutoFlow) Layout {
	l.autoFlow = f
	return l
}

// WithColSpan sets how many grid columns this block spans (default 1).
// Values less than 1 are treated as 1.
func (l Layout) WithColSpan(n int) Layout {
	if n < 1 {
		n = 1
	}
	l.colSpan = n
	return l
}

// WithRowSpan sets how many grid rows this block spans (default 1).
// Values less than 1 are treated as 1.
func (l Layout) WithRowSpan(n int) Layout {
	if n < 1 {
		n = 1
	}
	l.rowSpan = n
	return l
}

// ---------------------------------------------------------------------------
// Gap setters
// ---------------------------------------------------------------------------

// WithGap sets both column-gap and row-gap to the same value.
// Equivalent to CSS gap: n.
func (l Layout) WithGap(n int) Layout {
	if n < 0 {
		n = 0
	}
	l.columnGap = n
	l.rowGap = n
	return l
}

// WithColumnGap sets the horizontal gap between flex/grid children.
// Equivalent to CSS column-gap.
func (l Layout) WithColumnGap(n int) Layout {
	if n < 0 {
		n = 0
	}
	l.columnGap = n
	return l
}

// WithRowGap sets the vertical gap between flex/grid children.
// Equivalent to CSS row-gap.
func (l Layout) WithRowGap(n int) Layout {
	if n < 0 {
		n = 0
	}
	l.rowGap = n
	return l
}

// ---------------------------------------------------------------------------
// Padding setters
// ---------------------------------------------------------------------------

// WithPadding sets padding on all four sides using CSS order:
// (top, right, bottom, left). Negative values are clamped to 0.
func (l Layout) WithPadding(top, right, bottom, left int) Layout {
	l.padding = Edges{
		Top:    max0(top),
		Right:  max0(right),
		Bottom: max0(bottom),
		Left:   max0(left),
	}
	return l
}

// WithPaddingX sets left and right padding to the same value.
// Equivalent to CSS padding-inline.
func (l Layout) WithPaddingX(n int) Layout {
	n = max0(n)
	l.padding.Left = n
	l.padding.Right = n
	return l
}

// WithPaddingY sets top and bottom padding to the same value.
// Equivalent to CSS padding-block.
func (l Layout) WithPaddingY(n int) Layout {
	n = max0(n)
	l.padding.Top = n
	l.padding.Bottom = n
	return l
}

// WithUniformPadding sets all four padding sides to the same value.
// Equivalent to CSS padding: n.
func (l Layout) WithUniformPadding(n int) Layout {
	n = max0(n)
	l.padding = uniform(n)
	return l
}

// WithPaddingTop sets the top padding. Negative values are clamped to 0.
func (l Layout) WithPaddingTop(n int) Layout {
	l.padding.Top = max0(n)
	return l
}

// WithPaddingRight sets the right padding. Negative values are clamped to 0.
func (l Layout) WithPaddingRight(n int) Layout {
	l.padding.Right = max0(n)
	return l
}

// WithPaddingBottom sets the bottom padding. Negative values are clamped to 0.
func (l Layout) WithPaddingBottom(n int) Layout {
	l.padding.Bottom = max0(n)
	return l
}

// WithPaddingLeft sets the left padding. Negative values are clamped to 0.
func (l Layout) WithPaddingLeft(n int) Layout {
	l.padding.Left = max0(n)
	return l
}

// ---------------------------------------------------------------------------
// Size constraint setters
// ---------------------------------------------------------------------------

// WithMinWidth sets the minimum rendered width in terminal columns.
// 0 means no minimum. Values less than 0 are clamped to 0.
func (l Layout) WithMinWidth(n int) Layout {
	l.minWidth = max0(n)
	return l
}

// WithMaxWidth sets the maximum rendered width in terminal columns.
// 0 means no limit. Values less than 0 are clamped to 0.
func (l Layout) WithMaxWidth(n int) Layout {
	l.maxWidth = max0(n)
	return l
}

// WithMinHeight sets the minimum rendered height in terminal rows.
// 0 means no minimum. Values less than 0 are clamped to 0.
func (l Layout) WithMinHeight(n int) Layout {
	l.minHeight = max0(n)
	return l
}

// WithMaxHeight sets the maximum rendered height in terminal rows.
// 0 means no limit. Values less than 0 are clamped to 0.
func (l Layout) WithMaxHeight(n int) Layout {
	l.maxHeight = max0(n)
	return l
}

// WithSize is a shorthand for setting both width constraints to the same value
// (min = max = n), effectively giving the block a fixed width.
func (l Layout) WithSize(width, height int) Layout {
	l.minWidth = max0(width)
	l.maxWidth = max0(width)
	l.minHeight = max0(height)
	l.maxHeight = max0(height)
	return l
}

// WithWidth is a shorthand for setting both min and max width to the same
// value, giving the block a fixed column width.
func (l Layout) WithWidth(n int) Layout {
	n = max0(n)
	l.minWidth = n
	l.maxWidth = n
	return l
}

// WithHeight is a shorthand for setting both min and max height to the same
// value, giving the block a fixed row height.
func (l Layout) WithHeight(n int) Layout {
	n = max0(n)
	l.minHeight = n
	l.maxHeight = n
	return l
}

// ---------------------------------------------------------------------------
// Overflow setter
// ---------------------------------------------------------------------------

// WithOverflow sets the overflow behaviour for both axes.
func (l Layout) WithOverflow(o Overflow) Layout {
	l.overflow = o
	return l
}

// ---------------------------------------------------------------------------
// Accessor methods
// ---------------------------------------------------------------------------

// Direction returns the main axis direction.
func (l Layout) Direction() Direction { return l.direction }

// JustifyContent returns the main-axis content distribution mode.
func (l Layout) JustifyContent() JustifyContent { return l.justifyContent }

// AlignItems returns the cross-axis alignment for children.
func (l Layout) AlignItems() AlignItems { return l.alignItems }

// AlignSelf returns the per-block cross-axis alignment override.
func (l Layout) AlignSelf() AlignSelf { return l.alignSelf }

// Wrap reports whether flex children may wrap onto additional lines.
func (l Layout) Wrap() bool { return l.wrap }

// Grow reports whether this block expands to fill available space.
func (l Layout) Grow() bool { return l.grow }

// Shrink reports whether this block may compress below its natural size.
func (l Layout) Shrink() bool { return l.shrink }

// Basis returns the explicit flex-basis in terminal columns (0 = auto).
func (l Layout) Basis() int { return l.basis }

// Columns returns the grid column track definitions.
func (l Layout) Columns() []TrackSize { return l.columns }

// Rows returns the grid row track definitions.
func (l Layout) Rows() []TrackSize { return l.rows }

// AutoFlow returns the grid auto-flow direction.
func (l Layout) AutoFlow() AutoFlow { return l.autoFlow }

// ColSpan returns the number of grid columns this block spans.
// Returns 1 when no span has been explicitly set.
func (l Layout) ColSpan() int {
	if l.colSpan < 1 {
		return 1
	}
	return l.colSpan
}

// RowSpan returns the number of grid rows this block spans.
// Returns 1 when no span has been explicitly set.
func (l Layout) RowSpan() int {
	if l.rowSpan < 1 {
		return 1
	}
	return l.rowSpan
}

// ColumnGap returns the horizontal gap between children in terminal columns.
func (l Layout) ColumnGap() int { return l.columnGap }

// RowGap returns the vertical gap between children in terminal rows.
func (l Layout) RowGap() int { return l.rowGap }

// Padding returns the four-sided padding as an Edges value.
func (l Layout) Padding() Edges { return l.padding }

// PaddingTop returns the top padding in terminal rows.
func (l Layout) PaddingTop() int { return l.padding.Top }

// PaddingRight returns the right padding in terminal columns.
func (l Layout) PaddingRight() int { return l.padding.Right }

// PaddingBottom returns the bottom padding in terminal rows.
func (l Layout) PaddingBottom() int { return l.padding.Bottom }

// PaddingLeft returns the left padding in terminal columns.
func (l Layout) PaddingLeft() int { return l.padding.Left }

// MinWidth returns the minimum rendered width in terminal columns (0 = none).
func (l Layout) MinWidth() int { return l.minWidth }

// MaxWidth returns the maximum rendered width in terminal columns (0 = none).
func (l Layout) MaxWidth() int { return l.maxWidth }

// MinHeight returns the minimum rendered height in terminal rows (0 = none).
func (l Layout) MinHeight() int { return l.minHeight }

// MaxHeight returns the maximum rendered height in terminal rows (0 = none).
func (l Layout) MaxHeight() int { return l.maxHeight }

// Overflow returns the overflow behaviour.
func (l Layout) Overflow() Overflow { return l.overflow }

// IsZero reports whether the Layout has no effect (zero value / no-op).
// A Layout is zero when it carries no padding, no size constraints, no gap,
// no grid tracks, and all enum fields are at their default (zero) values.
func (l Layout) IsZero() bool {
	return l.direction == DirectionColumn &&
		l.justifyContent == JustifyStart &&
		l.alignItems == AlignItemsStart &&
		l.alignSelf == AlignSelfAuto &&
		!l.wrap &&
		!l.grow &&
		!l.shrink &&
		l.basis == 0 &&
		len(l.columns) == 0 &&
		len(l.rows) == 0 &&
		l.autoFlow == AutoFlowRow &&
		l.colSpan == 0 &&
		l.rowSpan == 0 &&
		l.columnGap == 0 &&
		l.rowGap == 0 &&
		l.padding == (Edges{}) &&
		l.minWidth == 0 &&
		l.maxWidth == 0 &&
		l.minHeight == 0 &&
		l.maxHeight == 0 &&
		l.overflow == OverflowClip
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// max0 clamps v to a minimum of 0.
func max0(v int) int {
	if v < 0 {
		return 0
	}
	return v
}
