package ink

import (
	"strings"
	"unicode/utf8"
)

// Strip removes all ANSI/VT escape sequences from s and returns the plain
// text. It handles:
//   - SGR sequences:           ESC [ … m
//   - CSI sequences (general): ESC [ … <final byte 0x40–0x7E>
//   - OSC sequences:           ESC ] … ST  (ST = BEL or ESC \)
//   - Single-char escapes:     ESC <byte>
//
// Strip never panics on malformed or partial sequences.
// Calling Strip on a string that contains no ESC byte is a zero-allocation
// fast path that returns s unchanged.
func Strip(s string) string {
	if !strings.ContainsRune(s, '\x1b') {
		return s
	}

	var b strings.Builder
	b.Grow(len(s))

	i := 0
	for i < len(s) {
		if s[i] != '\x1b' {
			// Fast path: copy a full UTF-8 rune without allocating.
			_, size := utf8.DecodeRuneInString(s[i:])
			b.WriteString(s[i : i+size])
			i += size
			continue
		}

		// ESC — inspect the next byte to determine sequence type.
		i++ // consume ESC
		if i >= len(s) {
			break // trailing lone ESC — discard
		}

		switch s[i] {
		case '[': // CSI sequence: ESC [ <params> <final byte>
			i++ // consume '['
			for i < len(s) {
				c := s[i]
				i++
				// Final byte is in the range 0x40–0x7E (@–~).
				if c >= 0x40 && c <= 0x7E {
					break
				}
			}

		case ']': // OSC sequence: ESC ] … BEL  or  ESC ] … ESC \
			i++ // consume ']'
			for i < len(s) {
				c := s[i]
				i++
				if c == '\x07' { // BEL terminates
					break
				}
				if c == '\x1b' && i < len(s) && s[i] == '\\' {
					i++ // consume '\' of ST (ESC \)
					break
				}
			}

		default:
			// Two-character escape: ESC <byte> — discard both.
			i++
		}
	}

	return b.String()
}

// ---------------------------------------------------------------------------
// Rune-column width helpers (used by applyLayout and applyBorder in render.go)
// ---------------------------------------------------------------------------

// runeWidth returns the number of terminal columns occupied by s, counting
// wide (CJK / emoji) runes as 2 columns and everything else as 1.
// Any ANSI escape sequences in s are stripped before measurement so that
// invisible bytes do not inflate the count.
func runeWidth(s string) int {
	plain := Strip(s)
	w := 0
	for _, r := range plain {
		w += runeColumnWidth(r)
	}
	return w
}

// runeColumnWidth returns the number of terminal columns occupied by r:
//   - 0 for C0/C1 control characters, combining marks, and zero-width joiners
//   - 2 for wide Unicode characters (CJK, fullwidth forms, common emoji)
//   - 1 for everything else
func runeColumnWidth(r rune) int {
	switch {
	case r < 0x20: // C0 controls
		return 0
	case r == 0x7F: // DEL
		return 0
	case r >= 0x80 && r < 0xA0: // C1 controls
		return 0
	// Combining / zero-width characters
	case r >= 0x0300 && r <= 0x036F: // Combining Diacritical Marks
		return 0
	case r >= 0x200B && r <= 0x200F: // zero-width spaces / joiners
		return 0
	case r == 0xFEFF: // BOM / zero-width no-break space
		return 0
	// Wide blocks
	case r >= 0x1100 && r <= 0x115F: // Hangul Jamo
		return 2
	case r >= 0x2E80 && r <= 0x303E: // CJK Radicals / Kangxi / Bopomofo
		return 2
	case r >= 0x3040 && r <= 0x33FF: // Hiragana, Katakana, Bopomofo, CJK compat
		return 2
	case r >= 0x3400 && r <= 0x4DBF: // CJK Ext-A
		return 2
	case r >= 0x4E00 && r <= 0x9FFF: // CJK Unified Ideographs
		return 2
	case r >= 0xA000 && r <= 0xA4CF: // Yi
		return 2
	case r >= 0xAC00 && r <= 0xD7AF: // Hangul Syllables
		return 2
	case r >= 0xF900 && r <= 0xFAFF: // CJK Compat Ideographs
		return 2
	case r >= 0xFE10 && r <= 0xFE1F: // Vertical forms
		return 2
	case r >= 0xFE30 && r <= 0xFE6F: // CJK Compat Forms / Small Forms
		return 2
	case r >= 0xFF01 && r <= 0xFF60: // Fullwidth Latin / Katakana
		return 2
	case r >= 0xFFE0 && r <= 0xFFE6: // Fullwidth signs
		return 2
	case r >= 0x1B000 && r <= 0x1B0FF: // Kana Supplement
		return 2
	case r >= 0x1F004 && r <= 0x1F0CF: // Playing cards / mahjong
		return 2
	case r >= 0x1F300 && r <= 0x1F9FF: // Misc symbols, emoji
		return 2
	case r >= 0x20000 && r <= 0x2FFFD: // CJK Ext-B through G
		return 2
	case r >= 0x30000 && r <= 0x3FFFD: // CJK Ext-H+
		return 2
	default:
		return 1
	}
}

// truncateToWidth truncates line so that its visible column width does not
// exceed maxW. When truncation is necessary a single '…' (U+2026, 1 column)
// is appended, ensuring the result never exceeds maxW columns.
func truncateToWidth(line string, maxW int) string {
	if runeWidth(line) <= maxW {
		return line
	}
	// Reserve one column for the ellipsis.
	budget := maxW - 1
	var sb strings.Builder
	w := 0
	for _, r := range line {
		cw := runeColumnWidth(r)
		if w+cw > budget {
			break
		}
		sb.WriteRune(r)
		w += cw
	}
	sb.WriteRune('…')
	return sb.String()
}
