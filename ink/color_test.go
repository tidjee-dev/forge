package ink

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// RGB / ANSI16 / ANSI256 constructors
// ---------------------------------------------------------------------------

func TestRGB_StoresComponents(t *testing.T) {
	c := RGB(10, 20, 30)
	if c.r != 10 || c.g != 20 || c.b != 30 {
		t.Errorf("RGB(10,20,30): r=%d g=%d b=%d, want 10 20 30", c.r, c.g, c.b)
	}
	if c.mode != colorRGB {
		t.Errorf("RGB mode = %v, want colorRGB", c.mode)
	}
	if c.IsZeroColor() {
		t.Error("RGB(...).IsZeroColor() = true, want false")
	}
}

func TestRGB_Extremes(t *testing.T) {
	black := RGB(0, 0, 0)
	white := RGB(255, 255, 255)
	if black.IsZeroColor() {
		t.Error("RGB(0,0,0).IsZeroColor() = true, want false (black is a valid color)")
	}
	if white.IsZeroColor() {
		t.Error("RGB(255,255,255).IsZeroColor() = true, want false")
	}
}

func TestANSI16_StoresCode(t *testing.T) {
	c := ANSI16(3)
	if c.code != 3 {
		t.Errorf("ANSI16(3).code = %d, want 3", c.code)
	}
	if c.mode != colorANSI16 {
		t.Errorf("ANSI16 mode = %v, want colorANSI16", c.mode)
	}
	if c.IsZeroColor() {
		t.Error("ANSI16(3).IsZeroColor() = true, want false")
	}
}

func TestANSI256_StoresCode(t *testing.T) {
	c := ANSI256(200)
	if c.code != 200 {
		t.Errorf("ANSI256(200).code = %d, want 200", c.code)
	}
	if c.mode != colorANSI256 {
		t.Errorf("ANSI256 mode = %v, want colorANSI256", c.mode)
	}
	if c.IsZeroColor() {
		t.Error("ANSI256(200).IsZeroColor() = true, want false")
	}
}

// ---------------------------------------------------------------------------
// IsZeroColor
// ---------------------------------------------------------------------------

func TestColor_IsZeroColor(t *testing.T) {
	cases := []struct {
		name   string
		c      Color
		isZero bool
	}{
		{"zero value", Color{}, true},
		{"RGB black", RGB(0, 0, 0), false},
		{"RGB white", RGB(255, 255, 255), false},
		{"RGB arbitrary", RGB(1, 2, 3), false},
		{"ANSI16(0)", ANSI16(0), false},
		{"ANSI16(15)", ANSI16(15), false},
		{"ANSI256(0)", ANSI256(0), false},
		{"ANSI256(255)", ANSI256(255), false},
		{"HexToRGB valid", HexToRGB("#ff0000"), false},
		{"HexToRGB invalid", HexToRGB("bad"), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.c.IsZeroColor()
			if got != tc.isZero {
				t.Errorf("IsZeroColor() = %v, want %v", got, tc.isZero)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Color.String
// ---------------------------------------------------------------------------

func TestColor_String(t *testing.T) {
	cases := []struct {
		name string
		c    Color
		want string
	}{
		{"zero color", Color{}, "Color(None)"},
		{"RGB", RGB(255, 0, 0), "RGB(255, 0, 0)"},
		{"RGB zeros", RGB(0, 0, 0), "RGB(0, 0, 0)"},
		{"ANSI16", ANSI16(1), "ANSI16(1)"},
		{"ANSI16 bright", ANSI16(9), "ANSI16(9)"},
		{"ANSI256", ANSI256(200), "ANSI256(200)"},
		{"ANSI256 zero", ANSI256(0), "ANSI256(0)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.c.String()
			if got != tc.want {
				t.Errorf("Color.String() = %q, want %q", got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// ANSI16 SGR codes — normal (30–37) vs bright (90–97) branches
// ---------------------------------------------------------------------------

func TestANSI16_FgParams_NormalRange(t *testing.T) {
	// codes 0–7 → SGR 30–37
	for code := uint8(0); code < 8; code++ {
		c := ANSI16(code)
		params := c.fgParams()
		if len(params) != 1 {
			t.Errorf("ANSI16(%d) fgParams len = %d, want 1", code, len(params))
			continue
		}
		want := strings.Join([]string{fmt.Sprintf("%d", 30+int(code))}, "")
		_ = want
		// Verify the numeric value: 30 + code
		expected := 30 + int(code)
		var got int
		if _, err := fmt.Sscanf(params[0], "%d", &got); err != nil {
			t.Errorf("ANSI16(%d) fgParams[0] parse error: %v", code, err)
			continue
		}
		if got != expected {
			t.Errorf("ANSI16(%d) fgParams[0] = %q, want %d", code, params[0], expected)
		}
	}
}

func TestANSI16_FgParams_BrightRange(t *testing.T) {
	// codes 8–15 → SGR 90–97  (formula: 82 + code)
	for code := uint8(8); code < 16; code++ {
		c := ANSI16(code)
		params := c.fgParams()
		if len(params) != 1 {
			t.Errorf("ANSI16(%d) fgParams len = %d, want 1", code, len(params))
			continue
		}
		expected := 82 + int(code) // e.g. code=8 → 90
		var got int
		if _, err := fmt.Sscanf(params[0], "%d", &got); err != nil {
			t.Errorf("ANSI16(%d) fgParams[0] parse error: %v", code, err)
			continue
		}
		if got != expected {
			t.Errorf("ANSI16(%d) fgParams[0] = %q, want %d", code, params[0], expected)
		}
	}
}

func TestANSI16_BgParams_NormalRange(t *testing.T) {
	// codes 0–7 → SGR 40–47
	for code := uint8(0); code < 8; code++ {
		c := ANSI16(code)
		params := c.bgParams()
		if len(params) != 1 {
			t.Errorf("ANSI16(%d) bgParams len = %d, want 1", code, len(params))
			continue
		}
		expected := 40 + int(code)
		var got int
		if _, err := fmt.Sscanf(params[0], "%d", &got); err != nil {
			t.Errorf("ANSI16(%d) bgParams[0] parse error: %v", code, err)
			continue
		}
		if got != expected {
			t.Errorf("ANSI16(%d) bgParams[0] = %q, want %d", code, params[0], expected)
		}
	}
}

func TestANSI16_BgParams_BrightRange(t *testing.T) {
	// codes 8–15 → SGR 100–107  (formula: 92 + code)
	for code := uint8(8); code < 16; code++ {
		c := ANSI16(code)
		params := c.bgParams()
		if len(params) != 1 {
			t.Errorf("ANSI16(%d) bgParams len = %d, want 1", code, len(params))
			continue
		}
		expected := 92 + int(code)
		var got int
		if _, err := fmt.Sscanf(params[0], "%d", &got); err != nil {
			t.Errorf("ANSI16(%d) bgParams[0] parse error: %v", code, err)
			continue
		}
		if got != expected {
			t.Errorf("ANSI16(%d) bgParams[0] = %q, want %d", code, params[0], expected)
		}
	}
}

func TestANSI16_Render_BrightFg(t *testing.T) {
	// End-to-end: bright ANSI16 fg produces a sequence in the 90–97 range.
	old := globalColorMode.Load()
	t.Cleanup(func() { globalColorMode.Store(old) })
	SetGlobalColorMode(colorModeAlways)

	got := New().WithForeground(ANSI16(9)).Render("hi") // code 9 → SGR 91
	if !strings.Contains(got, "91") {
		t.Errorf("ANSI16(9) bright fg Render = %q, want SGR code 91", got)
	}
}

// ---------------------------------------------------------------------------
// ContrastedColor
// ---------------------------------------------------------------------------

func TestContrastedColor(t *testing.T) {
	white := RGB(255, 255, 255)
	black := RGB(0, 0, 0)

	cases := []struct {
		name string
		bg   Color
		want Color
	}{
		{
			name: "dark background returns white",
			bg:   RGB(0, 0, 0),
			want: white,
		},
		{
			name: "very dark navy returns white",
			bg:   HexToRGB("#1a1a2e"),
			want: white,
		},
		{
			name: "light background returns black",
			bg:   RGB(255, 255, 255),
			want: black,
		},
		{
			name: "light gray returns black",
			bg:   HexToRGB("#f0f0f0"),
			want: black,
		},
		{
			name: "non-RGB bg returns white (fallback)",
			bg:   ANSI16(0),
			want: white,
		},
		{
			name: "zero color bg returns white (fallback)",
			bg:   Color{},
			want: white,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ContrastedColor(tc.bg)
			if got != tc.want {
				t.Errorf("ContrastedColor(%v) = %v, want %v", tc.bg, got, tc.want)
			}
		})
	}
}

func TestContrastedColor_Symmetry_With_ContrastedColorWith(t *testing.T) {
	// ContrastedColor(bg) should behave identically to
	// ContrastedColorWith(bg, <any failing fg>, 4.5).
	bgs := []Color{
		RGB(0, 0, 0),
		RGB(255, 255, 255),
		HexToRGB("#1a1a2e"),
		HexToRGB("#f0f0f0"),
	}
	for _, bg := range bgs {
		// Pass a fg that will never pass WCAG AA so the fallback always fires.
		withFallback := ContrastedColorWith(bg, bg, 4.5) // same color → ratio 1.0 → fallback
		direct := ContrastedColor(bg)
		if withFallback != direct {
			t.Errorf("bg=%v: ContrastedColor=%v, ContrastedColorWith(fallback)=%v — mismatch",
				bg, direct, withFallback)
		}
	}
}

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   Color
		isZero bool
	}{
		{
			name:  "full hex lowercase",
			input: "#aabbcc",
			want:  RGB(0xAA, 0xBB, 0xCC),
		},
		{
			name:  "full hex uppercase",
			input: "#AABBCC",
			want:  RGB(0xAA, 0xBB, 0xCC),
		},
		{
			name:  "full hex mixed case",
			input: "#AaBbCc",
			want:  RGB(0xAA, 0xBB, 0xCC),
		},
		{
			name:  "full hex black",
			input: "#000000",
			want:  RGB(0x00, 0x00, 0x00),
		},
		{
			name:  "full hex white",
			input: "#ffffff",
			want:  RGB(0xFF, 0xFF, 0xFF),
		},
		{
			name:  "full hex red",
			input: "#ff0000",
			want:  RGB(0xFF, 0x00, 0x00),
		},
		{
			name:  "short hex lowercase",
			input: "#abc",
			want:  RGB(0xAA, 0xBB, 0xCC),
		},
		{
			name:  "short hex uppercase",
			input: "#ABC",
			want:  RGB(0xAA, 0xBB, 0xCC),
		},
		{
			name:  "short hex black",
			input: "#000",
			want:  RGB(0x00, 0x00, 0x00),
		},
		{
			name:  "short hex white",
			input: "#fff",
			want:  RGB(0xFF, 0xFF, 0xFF),
		},
		{
			name:  "short hex red",
			input: "#f00",
			want:  RGB(0xFF, 0x00, 0x00),
		},
		{
			name:   "invalid: empty string",
			input:  "",
			isZero: true,
		},
		{
			name:   "invalid: no hash",
			input:  "aabbcc",
			isZero: true,
		},
		{
			name:   "invalid: too short",
			input:  "#ab",
			isZero: true,
		},
		{
			name:   "invalid: too long",
			input:  "#aabbccdd",
			isZero: true,
		},
		{
			name:   "invalid: non-hex characters (long)",
			input:  "#xxyyzz",
			isZero: true,
		},
		{
			name:   "invalid: non-hex characters (short)",
			input:  "#xyz",
			isZero: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HexToRGB(tt.input)

			if tt.isZero {
				if !got.IsZeroColor() {
					t.Errorf("HexToRGB(%q) = %+v, want zero Color", tt.input, got)
				}
				return
			}

			if got != tt.want {
				t.Errorf("HexToRGB(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestContrastRatio(t *testing.T) {
	tests := []struct {
		name    string
		c1, c2  Color
		wantMin float64
		wantMax float64
	}{
		{
			name:    "black on white = max contrast (21:1)",
			c1:      RGB(0, 0, 0),
			c2:      RGB(255, 255, 255),
			wantMin: 20.9,
			wantMax: 21.1,
		},
		{
			name:    "white on black = same as black on white (order-independent)",
			c1:      RGB(255, 255, 255),
			c2:      RGB(0, 0, 0),
			wantMin: 20.9,
			wantMax: 21.1,
		},
		{
			name:    "same color = min contrast (1:1)",
			c1:      RGB(128, 128, 128),
			c2:      RGB(128, 128, 128),
			wantMin: 0.99,
			wantMax: 1.01,
		},
		{
			name:    "non-RGB fallback returns 1.0",
			c1:      ANSI16(1),
			c2:      RGB(255, 255, 255),
			wantMin: 0.99,
			wantMax: 1.01,
		},
		{
			name:    "dark bg vs white fg passes WCAG AA (>= 4.5)",
			c1:      HexToRGB("#1a1a2e"),
			c2:      RGB(255, 255, 255),
			wantMin: 4.5,
			wantMax: 21.1,
		},
		{
			name:    "dark bg vs dark fg fails WCAG AA (< 4.5)",
			c1:      HexToRGB("#1a1a2e"),
			c2:      HexToRGB("#222244"),
			wantMin: 1.0,
			wantMax: 4.4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContrastRatio(tt.c1, tt.c2)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("ContrastRatio(%v, %v) = %.4f, want in [%.2f, %.2f]",
					tt.c1, tt.c2, got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestContrastedColorWith(t *testing.T) {
	darkBg := HexToRGB("#1a1a2e")
	lightBg := HexToRGB("#f0f0f0")
	white := RGB(255, 255, 255)
	black := RGB(0, 0, 0)

	tests := []struct {
		name     string
		bg       Color
		fg       Color
		minRatio float64
		want     Color
	}{
		{
			name:     "fg passes threshold — returned as-is",
			bg:       darkBg,
			fg:       white,
			minRatio: 4.5,
			want:     white,
		},
		{
			name:     "fg fails threshold — fallback to white on dark bg",
			bg:       darkBg,
			fg:       black,
			minRatio: 4.5,
			want:     white,
		},
		{
			name:     "fg fails threshold — fallback to black on light bg",
			bg:       lightBg,
			fg:       white,
			minRatio: 4.5,
			want:     black,
		},
		{
			name:     "fg passes threshold — returned as-is on light bg",
			bg:       lightBg,
			fg:       black,
			minRatio: 4.5,
			want:     black,
		},
		{
			name:     "zero minRatio defaults to 4.5",
			bg:       darkBg,
			fg:       white,
			minRatio: 0,
			want:     white,
		},
		{
			name:     "low minRatio (3.0) accepts fg that 4.5 would reject",
			bg:       lightBg,
			fg:       HexToRGB("#767676"), // ~4.48 contrast on white, fails 4.5 but passes 3.0
			minRatio: 3.0,
			want:     HexToRGB("#767676"),
		},
		{
			name:     "non-RGB fg with passing contrast",
			bg:       darkBg,
			fg:       ANSI16(7),
			minRatio: 4.5,
			// ANSI16 can't be measured, ContrastRatio returns 1.0 → fallback
			want: white,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContrastedColorWith(tt.bg, tt.fg, tt.minRatio)
			if got != tt.want {
				t.Errorf("ContrastedColorWith(%v, %v, %.1f) = %v, want %v",
					tt.bg, tt.fg, tt.minRatio, got, tt.want)
			}
		})
	}
}

func TestLighten(t *testing.T) {
	tests := []struct {
		name   string
		input  Color
		factor float64
		want   Color
	}{
		{
			name:   "factor=0 returns original color unchanged",
			input:  RGB(100, 150, 200),
			factor: 0,
			want:   RGB(100, 150, 200),
		},
		{
			name:   "factor=1 returns white",
			input:  RGB(100, 150, 200),
			factor: 1,
			want:   RGB(255, 255, 255),
		},
		{
			name:   "factor=0.5 moves halfway to white",
			input:  RGB(0, 0, 0),
			factor: 0.5,
			want:   RGB(128, 128, 128),
		},
		{
			name:   "factor clamped below 0",
			input:  RGB(100, 150, 200),
			factor: -1,
			want:   RGB(100, 150, 200),
		},
		{
			name:   "factor clamped above 1",
			input:  RGB(100, 150, 200),
			factor: 2,
			want:   RGB(255, 255, 255),
		},
		{
			name:   "already white stays white",
			input:  RGB(255, 255, 255),
			factor: 0.5,
			want:   RGB(255, 255, 255),
		},
		{
			name:   "non-RGB input returned unchanged",
			input:  ANSI16(3),
			factor: 0.5,
			want:   ANSI16(3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Lighten(tt.input, tt.factor)
			if got != tt.want {
				t.Errorf("Lighten(%v, %v) = %v, want %v", tt.input, tt.factor, got, tt.want)
			}
		})
	}
}

func TestDarken(t *testing.T) {
	tests := []struct {
		name   string
		input  Color
		factor float64
		want   Color
	}{
		{
			name:   "factor=0 returns original color unchanged",
			input:  RGB(100, 150, 200),
			factor: 0,
			want:   RGB(100, 150, 200),
		},
		{
			name:   "factor=1 returns black",
			input:  RGB(100, 150, 200),
			factor: 1,
			want:   RGB(0, 0, 0),
		},
		{
			name:   "factor=0.5 moves halfway to black",
			input:  RGB(255, 255, 255),
			factor: 0.5,
			want:   RGB(128, 128, 128),
		},
		{
			name:   "factor clamped below 0",
			input:  RGB(100, 150, 200),
			factor: -1,
			want:   RGB(100, 150, 200),
		},
		{
			name:   "factor clamped above 1",
			input:  RGB(100, 150, 200),
			factor: 2,
			want:   RGB(0, 0, 0),
		},
		{
			name:   "already black stays black",
			input:  RGB(0, 0, 0),
			factor: 0.5,
			want:   RGB(0, 0, 0),
		},
		{
			name:   "non-RGB input returned unchanged",
			input:  ANSI256(200),
			factor: 0.5,
			want:   ANSI256(200),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Darken(tt.input, tt.factor)
			if got != tt.want {
				t.Errorf("Darken(%v, %v) = %v, want %v", tt.input, tt.factor, got, tt.want)
			}
		})
	}
}

func TestMix(t *testing.T) {
	tests := []struct {
		name   string
		c1, c2 Color
		factor float64
		want   Color
	}{
		{
			name:   "factor=0 returns c1 unchanged",
			c1:     RGB(255, 0, 0),
			c2:     RGB(0, 0, 255),
			factor: 0,
			want:   RGB(255, 0, 0),
		},
		{
			name:   "factor=1 returns c2",
			c1:     RGB(255, 0, 0),
			c2:     RGB(0, 0, 255),
			factor: 1,
			want:   RGB(0, 0, 255),
		},
		{
			name:   "factor=0.5 returns midpoint",
			c1:     RGB(0, 0, 0),
			c2:     RGB(255, 255, 255),
			factor: 0.5,
			want:   RGB(128, 128, 128),
		},
		{
			name:   "factor clamped below 0 returns c1",
			c1:     RGB(255, 0, 0),
			c2:     RGB(0, 0, 255),
			factor: -1,
			want:   RGB(255, 0, 0),
		},
		{
			name:   "factor clamped above 1 returns c2",
			c1:     RGB(255, 0, 0),
			c2:     RGB(0, 0, 255),
			factor: 2,
			want:   RGB(0, 0, 255),
		},
		{
			name:   "same color mixed with itself returns same color",
			c1:     RGB(128, 64, 32),
			c2:     RGB(128, 64, 32),
			factor: 0.5,
			want:   RGB(128, 64, 32),
		},
		{
			name:   "non-RGB c1 returns c1 unchanged",
			c1:     ANSI16(1),
			c2:     RGB(0, 255, 0),
			factor: 0.5,
			want:   ANSI16(1),
		},
		{
			name:   "non-RGB c2 returns c1 unchanged",
			c1:     RGB(255, 0, 0),
			c2:     ANSI16(2),
			factor: 0.5,
			want:   RGB(255, 0, 0),
		},
		{
			name:   "both non-RGB returns c1 unchanged",
			c1:     ANSI16(1),
			c2:     ANSI16(2),
			factor: 0.5,
			want:   ANSI16(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Mix(tt.c1, tt.c2, tt.factor)
			if got != tt.want {
				t.Errorf("Mix(%v, %v, %v) = %v, want %v", tt.c1, tt.c2, tt.factor, got, tt.want)
			}
		})
	}
}

func TestContrastRatioSymmetry(t *testing.T) {
	pairs := [][2]Color{
		{Black, White},
		{Red, Blue},
		{HexToRGB("#3a86ff"), HexToRGB("#ffbe0b")},
	}
	for _, p := range pairs {
		r1 := ContrastRatio(p[0], p[1])
		r2 := ContrastRatio(p[1], p[0])
		if math.Abs(r1-r2) > 1e-9 {
			t.Errorf("ContrastRatio not symmetric: (%v,%v)=%.6f vs (%v,%v)=%.6f",
				p[0], p[1], r1, p[1], p[0], r2)
		}
	}
}
