package ink

import (
	"sort"
	"sync"
)

// Theme is a concurrent-safe named map of [Style] values. It allows
// a collection of named styles to be defined once and reused by name
// throughout an application.
//
// All read and write operations are safe for concurrent use.
//
// Build a theme and populate it:
//
//	t := ink.NewTheme()
//	t.Set("title", ink.New().WithForeground(ink.Cyan).WithBold(true))
//	t.Set("error", ink.New().WithForeground(ink.Danger).WithBold(true))
//
//	fmt.Println(t.Render("title", "Hello!"))
type Theme struct {
	mu     sync.RWMutex
	styles map[string]Style
}

// NewTheme returns an empty, ready-to-use Theme.
func NewTheme() *Theme {
	return &Theme{
		styles: make(map[string]Style),
	}
}

// Set stores the Style under name, replacing any existing value.
// It returns the receiver so calls can be chained.
//
//	t.Set("header", ink.New().WithBold(true)).
//	  Set("muted",  ink.New().WithForeground(ink.Muted))
func (t *Theme) Set(name string, s Style) *Theme {
	t.mu.Lock()
	t.styles[name] = s
	t.mu.Unlock()
	return t
}

// Get returns the Style stored under name and whether it was found.
//
//	if s, ok := t.Get("header"); ok {
//	    fmt.Println(s.Render("title"))
//	}
func (t *Theme) Get(name string) (Style, bool) {
	t.mu.RLock()
	s, ok := t.styles[name]
	t.mu.RUnlock()
	return s, ok
}

// Delete removes the entry for name. It is a no-op if name is not present.
// It returns the receiver so calls can be chained.
func (t *Theme) Delete(name string) *Theme {
	t.mu.Lock()
	delete(t.styles, name)
	t.mu.Unlock()
	return t
}

// Names returns a sorted slice of all style names currently in the theme.
func (t *Theme) Names() []string {
	t.mu.RLock()
	names := make([]string, 0, len(t.styles))
	for k := range t.styles {
		names = append(names, k)
	}
	t.mu.RUnlock()
	sort.Strings(names)
	return names
}

// Render applies the named style to text and returns the result.
// If name is not present in the theme, text is returned unchanged and no
// panic occurs.
//
//	output := t.Render("error", "something went wrong")
func (t *Theme) Render(name, text string) string {
	t.mu.RLock()
	s, ok := t.styles[name]
	t.mu.RUnlock()
	if !ok {
		return text
	}
	return s.Render(text)
}

// Clone returns a deep copy of the theme. Subsequent mutations to the clone
// — adding, updating, or deleting styles — do not affect the original, and
// vice versa.
func (t *Theme) Clone() *Theme {
	t.mu.RLock()
	clone := &Theme{
		styles: make(map[string]Style, len(t.styles)),
	}
	for k, v := range t.styles {
		clone.styles[k] = v
	}
	t.mu.RUnlock()
	return clone
}

// Merge applies all entries from patch on top of the receiver: keys present
// in patch overwrite the corresponding keys in the receiver; keys absent from
// patch are left untouched. The receiver is mutated in place.
// It returns the receiver so calls can be chained.
//
//	base.Merge(overrides)
func (t *Theme) Merge(patch *Theme) *Theme {
	patch.mu.RLock()
	entries := make(map[string]Style, len(patch.styles))
	for k, v := range patch.styles {
		entries[k] = v
	}
	patch.mu.RUnlock()

	t.mu.Lock()
	for k, v := range entries {
		t.styles[k] = v
	}
	t.mu.Unlock()
	return t
}
