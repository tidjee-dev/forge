package ink

import (
	"fmt"
	"math"
)

type colorMode uint8

const (
	colorNone colorMode = iota
	colorANSI16
	colorANSI256
	colorRGB
)

// Color represents a terminal color. It can be an RGB true-color value, an
// ANSI 256-color palette entry, an ANSI 16-color code, or the zero value
// (meaning "no color / inherit from terminal default").
//
// Construct a Color with [RGB], [ANSI16], [ANSI256], or [Hex].
// The zero value of Color is valid and means "no color".
type Color struct {
	r, g, b uint8
	code    uint8
	mode    colorMode
}

// RGB constructs a true-color (24-bit) Color from red, green, and blue
// components in the range 0–255.
func RGB(r, g, b uint8) Color {
	return Color{r: r, g: g, b: b, mode: colorRGB}
}

// Hex parses a CSS-style hex color string and returns the corresponding
// RGB Color. Both the long form ("#rrggbb") and the short form ("#rgb") are
// accepted; the hash prefix is required. Letter case is ignored.
//
// If the input is empty, missing the "#" prefix, the wrong length, or contains
// non-hex characters, a zero Color (no color) is returned — Hex never
// panics.
func Hex(hex string) Color {
	var r, g, b uint8

	if len(hex) == 7 && hex[0] == '#' {
		n, _ := fmt.Sscanf(hex[1:], "%02x%02x%02x", &r, &g, &b)
		if n != 3 {
			return Color{}
		}
	} else if len(hex) == 4 && hex[0] == '#' {
		var r4, g4, b4 uint8
		n, _ := fmt.Sscanf(hex[1:], "%1x%1x%1x", &r4, &g4, &b4)
		if n != 3 {
			return Color{}
		}
		r = r4 * 0x11
		g = g4 * 0x11
		b = b4 * 0x11
	} else {
		return Color{mode: colorNone}
	}
	return RGB(r, g, b)
}

// ANSI16 constructs a Color that uses one of the 16 standard ANSI palette
// codes (0–7 for normal colors, 8–15 for bright variants). The mapping from
// code to actual color is terminal-defined.
func ANSI16(code uint8) Color {
	return Color{code: code, mode: colorANSI16}
}

// ANSI256 constructs a Color that uses the 256-color xterm palette (codes
// 0–255). Codes 0–15 are the standard ANSI colors, 16–231 a 6×6×6 color
// cube, and 232–255 a grayscale ramp.
func ANSI256(code uint8) Color {
	return Color{code: code, mode: colorANSI256}
}

// IsZeroColor reports whether c is the zero Color, i.e. "no color". A zero
// Color carries no color information and causes [Style] to omit the
// corresponding SGR parameter entirely.
func (c Color) IsZeroColor() bool {
	return c.mode == colorNone
}

// String returns a human-readable representation of c, suitable for debugging
// and test output. Examples: "RGB(255, 0, 0)", "ANSI16(1)", "ANSI256(200)",
// "Color(None)".
func (c Color) String() string {
	switch c.mode {
	case colorRGB:
		return fmt.Sprintf("RGB(%d, %d, %d)", c.r, c.g, c.b)
	case colorANSI16:
		return fmt.Sprintf("ANSI16(%d)", c.code)
	case colorANSI256:
		return fmt.Sprintf("ANSI256(%d)", c.code)
	default:
		return "Color(None)"
	}
}

// fgParams returns the SGR parameter strings needed to set c as the foreground
// color in an ANSI escape sequence. Returns nil for a zero Color.
func (c Color) fgParams() []string {
	switch c.mode {
	case colorRGB:
		return []string{"38", "2", fmt.Sprintf("%d", c.r), fmt.Sprintf("%d", c.g), fmt.Sprintf("%d", c.b)}
	case colorANSI16:
		if c.code < 8 {
			return []string{fmt.Sprintf("%d", 30+c.code)}
		}
		return []string{fmt.Sprintf("%d", 82+c.code)} // 90–97 bright
	case colorANSI256:
		return []string{"38", "5", fmt.Sprintf("%d", c.code)}
	}
	return nil
}

// bgParams returns the SGR parameter strings needed to set c as the background
// color in an ANSI escape sequence. Returns nil for a zero Color.
func (c Color) bgParams() []string {
	switch c.mode {
	case colorRGB:
		return []string{"48", "2", fmt.Sprintf("%d", c.r), fmt.Sprintf("%d", c.g), fmt.Sprintf("%d", c.b)}
	case colorANSI16:
		if c.code < 8 {
			return []string{fmt.Sprintf("%d", 40+c.code)}
		}
		return []string{fmt.Sprintf("%d", 92+c.code)} // 100–107 bright
	case colorANSI256:
		return []string{"48", "5", fmt.Sprintf("%d", c.code)}
	}
	return nil
}

// toRGB extracts the red, green, and blue components of c when c is an RGB
// color. ok is false for ANSI16, ANSI256, and zero colors.
func (c Color) toRGB() (r, g, b uint8, ok bool) {
	if c.mode == colorRGB {
		return c.r, c.g, c.b, true
	}
	return 0, 0, 0, false
}

func clampF(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func lerpU8(a, b uint8, t float64) uint8 {
	return uint8(math.Round(float64(a) + t*(float64(b)-float64(a))))
}

// Lighten returns a brighter version of c by mixing it toward white. factor
// is clamped to [0, 1]: 0 returns c unchanged, 1 returns pure white. Only
// RGB colors are modified; ANSI16 and ANSI256 colors are returned unchanged.
func Lighten(c Color, factor float64) Color {
	r, g, b, ok := c.toRGB()
	if !ok {
		return c
	}
	factor = clampF(factor)
	r = lerpU8(r, 255, factor)
	g = lerpU8(g, 255, factor)
	b = lerpU8(b, 255, factor)
	return RGB(r, g, b)
}

// Darken returns a darker version of c by mixing it toward black. factor is
// clamped to [0, 1]: 0 returns c unchanged, 1 returns pure black. Only RGB
// colors are modified; ANSI16 and ANSI256 colors are returned unchanged.
func Darken(c Color, factor float64) Color {
	r, g, b, ok := c.toRGB()
	if !ok {
		return c
	}
	factor = clampF(factor)
	r = lerpU8(r, 0, factor)
	g = lerpU8(g, 0, factor)
	b = lerpU8(b, 0, factor)
	return RGB(r, g, b)
}

// Mix blends c1 and c2 by linear interpolation. factor is clamped to [0, 1]:
// 0 returns c1, 1 returns c2, 0.5 returns the midpoint. If either color is
// not an RGB color, c1 is returned unchanged.
func Mix(c1, c2 Color, factor float64) Color {
	r1, g1, b1, ok1 := c1.toRGB()
	r2, g2, b2, ok2 := c2.toRGB()
	if !ok1 || !ok2 {
		return c1
	}
	factor = clampF(factor)
	r := lerpU8(r1, r2, factor)
	g := lerpU8(g1, g2, factor)
	b := lerpU8(b1, b2, factor)
	return RGB(r, g, b)
}

func relativeLuminance(r, g, b uint8) float64 {
	linearize := func(v uint8) float64 {
		s := float64(v) / 255.0
		if s <= 0.04045 {
			return s / 12.92
		}
		return math.Pow((s+0.055)/1.055, 2.4)
	}
	return 0.2126*linearize(r) + 0.7152*linearize(g) + 0.0722*linearize(b)
}

// ContrastRatio returns the WCAG 2.1 contrast ratio between c1 and c2,
// which is a value in the range [1.0, 21.0]. The order of the arguments does
// not affect the result. If either color is not an RGB color, 1.0 is returned.
//
// WCAG AA requires a ratio of at least 4.5:1 for normal text and 3.0:1 for
// large text. WCAG AAA requires 7.0:1.
func ContrastRatio(c1, c2 Color) float64 {
	r1, g1, b1, ok1 := c1.toRGB()
	r2, g2, b2, ok2 := c2.toRGB()
	if !ok1 || !ok2 {
		return 1.0
	}
	l1 := relativeLuminance(r1, g1, b1)
	l2 := relativeLuminance(r2, g2, b2)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

// ContrastedColor returns either white or black — whichever achieves the
// higher contrast ratio against bg. It is useful for picking a legible
// foreground color when the background color is known.
//
// If bg is not an RGB color, white is returned.
func ContrastedColor(bg Color) Color {
	r, g, b, ok := bg.toRGB()
	if !ok {
		return RGB(255, 255, 255)
	}
	lum := relativeLuminance(r, g, b)
	if lum < 0.179 {
		return RGB(255, 255, 255)
	}
	return RGB(0, 0, 0)
}

// ContrastedColorWith returns fg if its contrast ratio against bg meets or
// exceeds minRatio. Otherwise it falls back to [ContrastedColor](bg) —
// whichever of white or black gives the better contrast.
//
// If minRatio is ≤ 0 it defaults to 4.5 (WCAG AA for normal text).
func ContrastedColorWith(bg, fg Color, minRatio float64) Color {
	if minRatio <= 0 {
		minRatio = 4.5
	}
	if ContrastRatio(bg, fg) >= minRatio {
		return fg
	}
	return ContrastedColor(bg)
}
