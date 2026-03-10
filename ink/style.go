package ink

// Style describes the visual appearance of a piece of terminal text.
// All fields are unexported; use the fluent setters to build a style and
// the accessor methods to read it back.
// The zero value is a valid, no-op style.
//
// Build styles fluently:
//
//	s := ink.New().
//	    WithForeground(ink.Red).
//	    WithBackground(ink.Black).
//	    WithBold(true).
//	    WithLayout(ink.NewLayout().WithUniformPadding(1).WithJustifyContent(ink.JustifyCenter))
//
//	fmt.Println(s.Render("hello"))
type Style struct {
	fgColor       Color
	bgColor       Color
	bold          bool
	italic        bool
	underline     bool
	dim           bool
	strikethrough bool
	inverse       bool
	layout        Layout      // spatial sizing, padding, flex/grid alignment
	borderStyle   BorderStyle // glyph set.
	borderColor   Color       // foreground color of the border glyphs.
	borderSides   BorderSide  // which sides of the border to render.
}

// New returns a zero-value Style ready for configuration.
func New() Style {
	return Style{}
}

// Setters
func (s Style) WithForeground(c Color) Style {
	s.fgColor = c
	return s
}

func (s Style) WithBackground(c Color) Style {
	s.bgColor = c
	return s
}

func (s Style) WithBold(b bool) Style {
	s.bold = b
	return s
}

func (s Style) WithItalic(i bool) Style {
	s.italic = i
	return s
}

func (s Style) WithUnderline(u bool) Style {
	s.underline = u
	return s
}

func (s Style) WithDim(d bool) Style {
	s.dim = d
	return s
}

func (s Style) WithStrikethrough(st bool) Style {
	s.strikethrough = st
	return s
}

func (s Style) WithInverse(i bool) Style {
	s.inverse = i
	return s
}

// WithLayout attaches a Layout to the style, controlling padding, flex/grid
// alignment, size constraints, and overflow behaviour.
func (s Style) WithLayout(l Layout) Style {
	s.layout = l
	return s
}

func (s Style) WithBorder(style BorderStyle) Style {
	s.borderStyle = style
	s.borderSides = BorderSideAll
	return s
}

func (s Style) WithBorderStyle(style BorderStyle) Style {
	s.borderStyle = style
	if s.borderSides == BorderSideNone && !style.IsZero() {
		s.borderSides = BorderSideAll
	}
	return s
}

// WithBorderColor sets the foreground color applied to all border glyphs.
// Pass a zero [Color] to unset.
func (s Style) WithBorderColor(c Color) Style {
	s.borderColor = c
	return s
}

func (s Style) WithBorderSide(side BorderSide) Style {
	s.borderSides = side
	return s
}

// ---------------------------------------------------------------------------
// Accessors
// ---------------------------------------------------------------------------

// GetForeground returns the foreground color of the style.
func (s Style) GetForeground() Color {
	return s.fgColor
}

// GetBackground returns the background color of the style.
func (s Style) GetBackground() Color {
	return s.bgColor
}

// IsBold returns true if the style has bold text.
func (s Style) IsBold() bool {
	return s.bold
}

// IsItalic returns true if the style has italic text.
func (s Style) IsItalic() bool {
	return s.italic
}

// IsUnderline returns true if the style has underlined text.
func (s Style) IsUnderline() bool {
	return s.underline
}

// IsDim returns true if the style has dimmed text.
func (s Style) IsDim() bool {
	return s.dim
}

// IsStrikethrough returns true if the style has strikethrough text.
func (s Style) IsStrikethrough() bool {
	return s.strikethrough
}

// IsInverse returns true if the style has inverse colors (swapped foreground and background).
func (s Style) IsInverse() bool {
	return s.inverse
}

// GetLayout returns the Layout attached to this style.
// The zero Layout is returned when no layout has been set.
func (s Style) GetLayout() Layout {
	return s.layout
}

// GetBorderStyle returns the border glyph set of the style.
func (s Style) GetBorderStyle() BorderStyle {
	return s.borderStyle
}

// GetBorderColor returns the foreground color of the border glyphs.
func (s Style) GetBorderColor() Color {
	return s.borderColor
}

// GetBorderSide returns which sides of the border is render for the style.
func (s Style) GetBorderSide() BorderSide {
	return s.borderSides
}
