package ink

// ---------------------------------------------------------------------------
// Override helper types
//
// Each override type pairs a boolean "set" sentinel with the actual value so
// that Style.Merge can distinguish "not mentioned in the patch" from
// "explicitly set to the zero/false value".
// ---------------------------------------------------------------------------

// layoutOverride carries an explicit Layout override for Merge.
type layoutOverride struct {
	set    bool
	layout Layout
}

// boolOverride carries an explicit boolean override so that Merge can
// distinguish "not set" from "explicitly set to false".
type boolOverride struct {
	set   bool
	value bool
}

// borderOverride carries an explicit border override for Merge.
type borderOverride struct {
	set   bool
	style BorderStyle
}

// borderColorOverride carries an explicit border color override for Merge.
type borderColorOverride struct {
	set   bool
	color Color
}

// borderSidesOverride carries an explicit border sides override for Merge.
type borderSidesOverride struct {
	set   bool
	sides BorderSide
}

// ---------------------------------------------------------------------------
// StyleOverride
// ---------------------------------------------------------------------------

// StyleOverride is a style patch where boolean attributes can be explicitly
// set to false, allowing [Style.Merge] to unset a flag in the base style.
// Use [Override] to build one fluently, then pass it to [Style.Merge].
//
//	base := ink.New().Bold().Foreground(ink.Red)
//	patch := ink.Override().NoBold().Foreground(ink.Blue)
//	result := base.Merge(patch) // not bold, blue fg
type StyleOverride struct {
	fg          Color
	bg          Color
	bold        boolOverride
	italic      boolOverride
	underline   boolOverride
	dim         boolOverride
	strike      boolOverride
	inverse     boolOverride
	layout      layoutOverride
	border      borderOverride
	borderColor borderColorOverride
	borderSides borderSidesOverride
}

// Override begins building a StyleOverride that can be passed to [Style.Merge].
func Override() StyleOverride { return StyleOverride{} }

// ---------------------------------------------------------------------------
// StyleOverride fluent setters
// ---------------------------------------------------------------------------

// Foreground sets the foreground color override.
func (o StyleOverride) Foreground(c Color) StyleOverride {
	o.fg = c
	return o
}

// Background sets the background color override.
func (o StyleOverride) Background(c Color) StyleOverride {
	o.bg = c
	return o
}

// Bold explicitly enables bold.
func (o StyleOverride) Bold() StyleOverride {
	o.bold = boolOverride{true, true}
	return o
}

// NoBold explicitly disables bold.
func (o StyleOverride) NoBold() StyleOverride {
	o.bold = boolOverride{true, false}
	return o
}

// Italic explicitly enables italic.
func (o StyleOverride) Italic() StyleOverride {
	o.italic = boolOverride{true, true}
	return o
}

// NoItalic explicitly disables italic.
func (o StyleOverride) NoItalic() StyleOverride {
	o.italic = boolOverride{true, false}
	return o
}

// Underline explicitly enables underline.
func (o StyleOverride) Underline() StyleOverride {
	o.underline = boolOverride{true, true}
	return o
}

// NoUnderline explicitly disables underline.
func (o StyleOverride) NoUnderline() StyleOverride {
	o.underline = boolOverride{true, false}
	return o
}

// Dim explicitly enables dim.
func (o StyleOverride) Dim() StyleOverride {
	o.dim = boolOverride{true, true}
	return o
}

// NoDim explicitly disables dim.
func (o StyleOverride) NoDim() StyleOverride {
	o.dim = boolOverride{true, false}
	return o
}

// Strike explicitly enables strike-through.
func (o StyleOverride) Strike() StyleOverride {
	o.strike = boolOverride{true, true}
	return o
}

// NoStrike explicitly disables strike-through.
func (o StyleOverride) NoStrike() StyleOverride {
	o.strike = boolOverride{true, false}
	return o
}

// Inverse explicitly enables inverse.
func (o StyleOverride) Inverse() StyleOverride {
	o.inverse = boolOverride{true, true}
	return o
}

// NoInverse explicitly disables inverse.
func (o StyleOverride) NoInverse() StyleOverride {
	o.inverse = boolOverride{true, false}
	return o
}

// WithLayout sets the layout override. The entire Layout is replaced; use
// WithLayout(ink.NewLayout()) to explicitly clear a layout from a base style.
func (o StyleOverride) WithLayout(l Layout) StyleOverride {
	o.layout = layoutOverride{true, l}
	return o
}

// WithBorder sets the border style override.
func (o StyleOverride) WithBorder(style BorderStyle) StyleOverride {
	o.border = borderOverride{true, style}
	return o
}

// WithBorderColor sets the border foreground color override.
func (o StyleOverride) WithBorderColor(color Color) StyleOverride {
	o.borderColor = borderColorOverride{true, color}
	return o
}

// WithBorderSides sets the border sides override.
func (o StyleOverride) WithBorderSides(sides BorderSide) StyleOverride {
	o.borderSides = borderSidesOverride{true, sides}
	return o
}

// WithNoBorder explicitly removes the border (sets a zero BorderStyle,
// BorderNone sides, and clears the border color). Use this to strip a
// border from a base style via Merge.
func (o StyleOverride) WithNoBorder() StyleOverride {
	o.border = borderOverride{true, BorderStyle{}}
	o.borderSides = borderSidesOverride{true, BorderSideNone}
	o.borderColor = borderColorOverride{true, Color{}}
	return o
}

// ---------------------------------------------------------------------------
// Style.Merge
// ---------------------------------------------------------------------------

// Merge applies the given [StyleOverride] on top of the receiver and returns
// the patched style. Any attribute explicitly set in the override (including
// an explicit false) replaces the corresponding attribute in the base; unset
// attributes are left untouched. The receiver is never mutated.
//
//	base := ink.New().Bold().Foreground(ink.Red)
//	patch := ink.Override().NoBold().Foreground(ink.Blue)
//	result := base.Merge(patch) // not bold, blue fg
func (s Style) Merge(o StyleOverride) Style {
	if !o.fg.IsZeroColor() {
		s.fgColor = o.fg
	}
	if !o.bg.IsZeroColor() {
		s.bgColor = o.bg
	}
	if o.bold.set {
		s.bold = o.bold.value
	}
	if o.italic.set {
		s.italic = o.italic.value
	}
	if o.underline.set {
		s.underline = o.underline.value
	}
	if o.dim.set {
		s.dim = o.dim.value
	}
	if o.strike.set {
		s.strikethrough = o.strike.value
	}
	if o.inverse.set {
		s.inverse = o.inverse.value
	}
	if o.layout.set {
		s.layout = o.layout.layout
	}
	if o.border.set {
		s.borderStyle = o.border.style
	}
	if o.borderColor.set {
		s.borderColor = o.borderColor.color
	}
	if o.borderSides.set {
		s.borderSides = o.borderSides.sides
	}
	return s
}
