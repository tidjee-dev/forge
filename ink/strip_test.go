package ink

import (
	"strings"
	"sync"
	"testing"
)

// ---------------------------------------------------------------------------
// Strip — plain strings
// ---------------------------------------------------------------------------

func TestStrip_PlainString_Unchanged(t *testing.T) {
	inputs := []string{
		"hello",
		"hello, world!",
		"line1\nline2\nline3",
		"unicode: café résumé naïve",
		"CJK: 你好世界",
		"emoji: 🎉🔥",
		"  leading and trailing spaces  ",
	}
	for _, input := range inputs {
		got := Strip(input)
		if got != input {
			t.Errorf("Strip(%q) = %q, want unchanged", input, got)
		}
	}
}

func TestStrip_EmptyString(t *testing.T) {
	got := Strip("")
	if got != "" {
		t.Errorf("Strip(\"\") = %q, want \"\"", got)
	}
}

// ---------------------------------------------------------------------------
// Strip — output produced by Render recovers original
// ---------------------------------------------------------------------------

func TestStrip_RenderedOutput_RecoverOriginal(t *testing.T) {
	// Force colour on so Render actually emits SGR sequences.
	oldMode := globalColorMode.Load()
	t.Cleanup(func() {
		globalColorMode.Store(oldMode)
		envOnce = sync.Once{}
		envDisabled = false
		envForced = false
	})
	SetGlobalColorMode(colorModeAlways)

	cases := []struct {
		name  string
		text  string
		style Style
	}{
		{
			name:  "bold",
			text:  "hello",
			style: New().WithBold(true),
		},
		{
			name:  "italic",
			text:  "italic text",
			style: New().WithItalic(true),
		},
		{
			name:  "underline",
			text:  "underlined",
			style: New().WithUnderline(true),
		},
		{
			name:  "dim",
			text:  "dim text",
			style: New().WithDim(true),
		},
		{
			name:  "strikethrough",
			text:  "struck through",
			style: New().WithStrikethrough(true),
		},
		{
			name:  "inverse",
			text:  "inverted",
			style: New().WithInverse(true),
		},
		{
			name:  "fg RGB",
			text:  "red text",
			style: New().WithForeground(RGB(255, 0, 0)),
		},
		{
			name:  "bg RGB",
			text:  "blue background",
			style: New().WithBackground(RGB(0, 0, 255)),
		},
		{
			name:  "fg ANSI16",
			text:  "ansi16 text",
			style: New().WithForeground(ANSI16(1)),
		},
		{
			name:  "fg ANSI256",
			text:  "ansi256 text",
			style: New().WithForeground(ANSI256(200)),
		},
		{
			name:  "bg ANSI256",
			text:  "ansi256 bg",
			style: New().WithBackground(ANSI256(50)),
		},
		{
			name:  "bold + fg + bg combined",
			text:  "combined",
			style: New().WithBold(true).WithForeground(RGB(255, 128, 0)).WithBackground(RGB(30, 30, 30)),
		},
		{
			name:  "all SGR attributes combined",
			text:  "everything",
			style: New().WithBold(true).WithItalic(true).WithUnderline(true).WithStrikethrough(true).WithDim(true),
		},
		{
			name:  "multiline text",
			text:  "line one\nline two\nline three",
			style: New().WithBold(true).WithForeground(Blue),
		},
		{
			name:  "unicode content",
			text:  "café résumé",
			style: New().WithForeground(RGB(100, 200, 100)),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rendered := tc.style.Render(tc.text)
			got := Strip(rendered)
			if got != tc.text {
				t.Errorf("Strip(Render(%q)) = %q, want original text", tc.text, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Strip — known SGR sequences
// ---------------------------------------------------------------------------

func TestStrip_KnownSGRSequences(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "reset only",
			input: "\x1b[0m",
			want:  "",
		},
		{
			name:  "bold code",
			input: "\x1b[1mhello\x1b[0m",
			want:  "hello",
		},
		{
			name:  "fg 256-colour",
			input: "\x1b[38;5;200mtext\x1b[0m",
			want:  "text",
		},
		{
			name:  "fg RGB truecolour",
			input: "\x1b[38;2;255;128;0mtext\x1b[0m",
			want:  "text",
		},
		{
			name:  "bg RGB truecolour",
			input: "\x1b[48;2;0;64;128mtext\x1b[0m",
			want:  "text",
		},
		{
			name:  "multiple sequences",
			input: "\x1b[1m\x1b[3m\x1b[38;2;255;0;0mhello\x1b[0m",
			want:  "hello",
		},
		{
			name:  "sequence mid-string",
			input: "before\x1b[1mbold\x1b[0mafter",
			want:  "beforeboldafter",
		},
		{
			name:  "multiline with sequences",
			input: "\x1b[1mline1\x1b[0m\n\x1b[3mline2\x1b[0m",
			want:  "line1\nline2",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Strip(tc.input)
			if got != tc.want {
				t.Errorf("Strip(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Strip — OSC sequences (e.g. hyperlinks, window title)
// ---------------------------------------------------------------------------

func TestStrip_OSCSequences(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "OSC terminated by BEL",
			input: "\x1b]0;window title\x07text",
			want:  "text",
		},
		{
			name:  "OSC terminated by ST (ESC \\)",
			input: "\x1b]8;;https://example.com\x1b\\link text\x1b]8;;\x1b\\",
			want:  "link text",
		},
		{
			name:  "OSC between plain text",
			input: "before\x1b]0;title\x07after",
			want:  "beforeafter",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Strip(tc.input)
			if got != tc.want {
				t.Errorf("Strip(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Strip — malformed / partial sequences must not panic
// ---------------------------------------------------------------------------

func TestStrip_MalformedSequences_NoPanic(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "lone ESC at end",
			input: "text\x1b",
		},
		{
			name:  "ESC [ with no final byte",
			input: "\x1b[1",
		},
		{
			name:  "ESC [ empty",
			input: "\x1b[",
		},
		{
			name:  "ESC ] unterminated OSC",
			input: "\x1b]0;unterminated title",
		},
		{
			name:  "ESC followed by ESC",
			input: "\x1b\x1b[1mhi\x1b[0m",
		},
		{
			name:  "truncated RGB sequence",
			input: "\x1b[38;2;25",
		},
		{
			name:  "null bytes mixed with ESC",
			input: "\x00\x1b[1m\x00hi\x1b[0m\x00",
		},
		{
			name:  "only ESC bytes",
			input: "\x1b\x1b\x1b",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Must not panic — result is not validated beyond that.
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Strip panicked on %q: %v", tc.input, r)
				}
			}()
			_ = Strip(tc.input)
		})
	}
}

// ---------------------------------------------------------------------------
// Strip — idempotency
// ---------------------------------------------------------------------------

func TestStrip_Idempotent(t *testing.T) {
	// Stripping an already-stripped string should be a no-op.
	oldMode := globalColorMode.Load()
	t.Cleanup(func() {
		globalColorMode.Store(oldMode)
		envOnce = sync.Once{}
		envDisabled = false
		envForced = false
	})
	SetGlobalColorMode(colorModeAlways)

	texts := []string{"hello", "line1\nline2", "café 你好 🎉"}
	for _, text := range texts {
		rendered := New().WithBold(true).WithForeground(Red).Render(text)
		once := Strip(rendered)
		twice := Strip(once)
		if once != twice {
			t.Errorf("Strip not idempotent on %q: first=%q second=%q", text, once, twice)
		}
		if strings.ContainsRune(twice, '\x1b') {
			t.Errorf("Strip twice still contains ESC on %q", text)
		}
	}
}
