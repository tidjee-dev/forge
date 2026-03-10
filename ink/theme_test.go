package ink

import (
	"reflect"
	"slices"
	"sync"
	"testing"
)

// ---------------------------------------------------------------------------
// NewTheme / zero state
// ---------------------------------------------------------------------------

func TestNewTheme_Empty(t *testing.T) {
	th := NewTheme()
	if len(th.Names()) != 0 {
		t.Errorf("NewTheme().Names() = %v, want empty", th.Names())
	}
}

// ---------------------------------------------------------------------------
// Set / Get
// ---------------------------------------------------------------------------

func TestTheme_Set_Get_Found(t *testing.T) {
	th := NewTheme()
	s := New().WithBold(true).WithForeground(Red)
	th.Set("header", s)

	got, ok := th.Get("header")
	if !ok {
		t.Fatal("Get(\"header\") ok = false, want true")
	}
	if !reflect.DeepEqual(got, s) {
		t.Errorf("Get(\"header\") = %v, want %v", got, s)
	}
}

func TestTheme_Get_Missing(t *testing.T) {
	th := NewTheme()
	_, ok := th.Get("nonexistent")
	if ok {
		t.Error("Get on missing key returned ok = true, want false")
	}
}

func TestTheme_Set_Overwrites(t *testing.T) {
	th := NewTheme()
	first := New().WithBold(true)
	second := New().WithItalic(true)

	th.Set("key", first)
	th.Set("key", second)

	got, ok := th.Get("key")
	if !ok {
		t.Fatal("Get after overwrite: ok = false")
	}
	if !reflect.DeepEqual(got, second) {
		t.Errorf("Get after overwrite = %v, want second style", got)
	}
}

func TestTheme_Set_Chaining(t *testing.T) {
	th := NewTheme()
	returned := th.Set("a", New()).Set("b", New().WithBold(true))
	if returned != th {
		t.Error("Set chaining did not return the receiver")
	}
	if _, ok := th.Get("a"); !ok {
		t.Error("chained Set did not store 'a'")
	}
	if _, ok := th.Get("b"); !ok {
		t.Error("chained Set did not store 'b'")
	}
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestTheme_Delete_RemovesEntry(t *testing.T) {
	th := NewTheme()
	th.Set("x", New().WithBold(true))
	th.Delete("x")

	if _, ok := th.Get("x"); ok {
		t.Error("Get after Delete returned ok = true, want false")
	}
}

func TestTheme_Delete_MissingKey_NoOp(t *testing.T) {
	th := NewTheme()
	// Must not panic.
	th.Delete("does-not-exist")
}

func TestTheme_Delete_Chaining(t *testing.T) {
	th := NewTheme()
	th.Set("a", New()).Set("b", New())
	returned := th.Delete("a").Delete("b")
	if returned != th {
		t.Error("Delete chaining did not return the receiver")
	}
	if len(th.Names()) != 0 {
		t.Errorf("after chained Delete, Names() = %v, want empty", th.Names())
	}
}

// ---------------------------------------------------------------------------
// Names
// ---------------------------------------------------------------------------

func TestTheme_Names_Sorted(t *testing.T) {
	th := NewTheme()
	th.Set("zebra", New())
	th.Set("apple", New())
	th.Set("mango", New())

	names := th.Names()
	want := []string{"apple", "mango", "zebra"}
	if !slices.Equal(names, want) {
		t.Errorf("Names() = %v, want %v", names, want)
	}
}

func TestTheme_Names_EmptyTheme(t *testing.T) {
	th := NewTheme()
	names := th.Names()
	if len(names) != 0 {
		t.Errorf("Names() on empty theme = %v, want []", names)
	}
}

func TestTheme_Names_AfterDelete(t *testing.T) {
	th := NewTheme()
	th.Set("a", New()).Set("b", New()).Set("c", New())
	th.Delete("b")

	names := th.Names()
	want := []string{"a", "c"}
	if !slices.Equal(names, want) {
		t.Errorf("Names() after delete = %v, want %v", names, want)
	}
}

// ---------------------------------------------------------------------------
// Render
// ---------------------------------------------------------------------------

func TestTheme_Render_MissingKey_ReturnsTextUnchanged(t *testing.T) {
	th := NewTheme()
	got := th.Render("no-such-style", "hello")
	if got != "hello" {
		t.Errorf("Render missing key = %q, want %q", got, "hello")
	}
}

func TestTheme_Render_MissingKey_NoPanic(t *testing.T) {
	th := NewTheme()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Render on missing key panicked: %v", r)
		}
	}()
	_ = th.Render("ghost", "text")
}

func TestTheme_Render_ZeroStyle_ReturnsText(t *testing.T) {
	// Force colour on so that a non-zero style would produce SGR sequences.
	old := globalColorMode.Load()
	t.Cleanup(func() {
		globalColorMode.Store(old)
		envOnce = sync.Once{}
		envDisabled = false
		envForced = false
	})
	SetGlobalColorMode(colorModeAlways)

	th := NewTheme()
	th.Set("zero", New()) // zero style — no attributes
	got := th.Render("zero", "hello")
	if got != "hello" {
		t.Errorf("Render zero style = %q, want plain %q", got, "hello")
	}
}

func TestTheme_Render_AppliesStyle(t *testing.T) {
	// Force colour on so Render emits SGR.
	old := globalColorMode.Load()
	t.Cleanup(func() {
		globalColorMode.Store(old)
		envOnce = sync.Once{}
		envDisabled = false
		envForced = false
	})
	SetGlobalColorMode(colorModeAlways)

	th := NewTheme()
	th.Set("bold-red", New().WithBold(true).WithForeground(RGB(255, 0, 0)))

	got := th.Render("bold-red", "hello")
	// Must contain ESC sequences and the original text.
	if Strip(got) != "hello" {
		t.Errorf("Strip(Render) = %q, want %q", Strip(got), "hello")
	}
	if got == "hello" {
		t.Error("Render with bold-red style produced no SGR sequences")
	}
}

// ---------------------------------------------------------------------------
// Clone
// ---------------------------------------------------------------------------

func TestTheme_Clone_IsIndependent(t *testing.T) {
	original := NewTheme()
	original.Set("a", New().WithBold(true))
	original.Set("b", New().WithItalic(true))

	clone := original.Clone()

	// Mutate the clone — original must be unaffected.
	clone.Set("a", New().WithDim(true))
	clone.Set("c", New().WithUnderline(true))

	got, ok := original.Get("a")
	if !ok {
		t.Fatal("original lost key 'a' after clone mutation")
	}
	if !got.IsBold() {
		t.Error("original 'a' style was mutated by clone.Set")
	}

	if _, ok := original.Get("c"); ok {
		t.Error("original gained key 'c' from clone.Set")
	}
}

func TestTheme_Clone_ContainsAllOriginalEntries(t *testing.T) {
	original := NewTheme()
	original.Set("x", New().WithBold(true))
	original.Set("y", New().WithItalic(true))

	clone := original.Clone()

	for _, name := range []string{"x", "y"} {
		orig, _ := original.Get(name)
		got, ok := clone.Get(name)
		if !ok {
			t.Errorf("clone missing key %q", name)
			continue
		}
		if !reflect.DeepEqual(got, orig) {
			t.Errorf("clone[%q] = %v, want %v", name, got, orig)
		}
	}
}

func TestTheme_Clone_MutateOriginal_DoesNotAffectClone(t *testing.T) {
	original := NewTheme()
	original.Set("shared", New().WithBold(true))

	clone := original.Clone()

	// Now mutate the original.
	original.Set("shared", New().WithItalic(true))
	original.Set("new-key", New())

	// Clone must still see the old value of "shared".
	got, ok := clone.Get("shared")
	if !ok {
		t.Fatal("clone lost 'shared' after original mutation")
	}
	if !got.IsBold() {
		t.Error("clone 'shared' was affected by mutation of original")
	}

	// Clone must not see the new key added to the original.
	if _, ok := clone.Get("new-key"); ok {
		t.Error("clone gained 'new-key' from original.Set")
	}
}

// ---------------------------------------------------------------------------
// Merge
// ---------------------------------------------------------------------------

func TestTheme_Merge_PatchOverwritesBase(t *testing.T) {
	base := NewTheme()
	base.Set("shared", New().WithBold(true))
	base.Set("base-only", New().WithItalic(true))

	patch := NewTheme()
	patch.Set("shared", New().WithDim(true))
	patch.Set("patch-only", New().WithUnderline(true))

	base.Merge(patch)

	// "shared" must now be the patch value.
	got, ok := base.Get("shared")
	if !ok {
		t.Fatal("base lost 'shared' after Merge")
	}
	if !got.IsDim() {
		t.Error("Merge did not overwrite 'shared' with patch value")
	}

	// "base-only" must be preserved.
	baseOnly, ok := base.Get("base-only")
	if !ok {
		t.Fatal("Merge removed 'base-only'")
	}
	if !baseOnly.IsItalic() {
		t.Error("Merge corrupted 'base-only' style")
	}

	// "patch-only" must now be in base.
	patchOnly, ok := base.Get("patch-only")
	if !ok {
		t.Fatal("Merge did not add 'patch-only' to base")
	}
	if !patchOnly.IsUnderline() {
		t.Error("Merge did not preserve 'patch-only' style")
	}
}

func TestTheme_Merge_EmptyPatch_BaseUnchanged(t *testing.T) {
	base := NewTheme()
	base.Set("a", New().WithBold(true))

	base.Merge(NewTheme())

	got, ok := base.Get("a")
	if !ok {
		t.Fatal("Merge with empty patch removed 'a'")
	}
	if !got.IsBold() {
		t.Error("Merge with empty patch changed 'a'")
	}
}

func TestTheme_Merge_Chaining(t *testing.T) {
	base := NewTheme()
	patch := NewTheme().Set("x", New())
	returned := base.Merge(patch)
	if returned != base {
		t.Error("Merge did not return the receiver")
	}
}

// ---------------------------------------------------------------------------
// Concurrent safety (race detector)
// ---------------------------------------------------------------------------

func TestTheme_ConcurrentSetGet_RaceClean(t *testing.T) {
	th := NewTheme()
	th.Set("key", New().WithBold(true))

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	// Concurrent writers
	for i := range goroutines {
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				th.Set("key", New().WithBold(true))
			} else {
				th.Set("other", New().WithItalic(true))
			}
		}(i)
	}

	// Concurrent readers
	for range goroutines {
		go func() {
			defer wg.Done()
			_, _ = th.Get("key")
			_ = th.Names()
		}()
	}

	wg.Wait()
}

func TestTheme_ConcurrentSetDelete_RaceClean(t *testing.T) {
	th := NewTheme()
	for i := range 20 {
		th.Set(string(rune('a'+i)), New())
	}

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	keys := []string{"a", "b", "c", "d", "e"}

	// Concurrent deleters
	for i := range goroutines {
		go func(i int) {
			defer wg.Done()
			th.Delete(keys[i%len(keys)])
		}(i)
	}

	// Concurrent setters
	for i := range goroutines {
		go func(i int) {
			defer wg.Done()
			th.Set(keys[i%len(keys)], New().WithBold(true))
		}(i)
	}

	wg.Wait()
}

func TestTheme_ConcurrentNames_RaceClean(t *testing.T) {
	th := NewTheme()

	const goroutines = 40
	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	// Concurrent writers
	for i := range goroutines {
		go func(i int) {
			defer wg.Done()
			th.Set(string(rune('a'+i%26)), New())
		}(i)
	}

	// Concurrent Names() readers
	for range goroutines {
		go func() {
			defer wg.Done()
			_ = th.Names()
		}()
	}

	wg.Wait()
}
