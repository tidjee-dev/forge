package cast_test

import (
	"strings"
	"testing"

	"github.com/tidjee-dev/forge/cast"
	"github.com/tidjee-dev/forge/ink"
)

// ---------------------------------------------------------------------------
// Single node
// ---------------------------------------------------------------------------

func TestTree_SingleNode(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTree("root").Render()
	if got != "root" {
		t.Errorf("expected \"root\", got %q", got)
	}
}

func TestTree_SingleNodeNoTrailingNewline(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTree("node").Render()
	if strings.HasSuffix(got, "\n") {
		t.Errorf("expected no trailing newline, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// Flat list (single level of children)
// ---------------------------------------------------------------------------

func TestTree_FlatList(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTree("root").
		AddChild(cast.NewTree("a")).
		AddChild(cast.NewTree("b")).
		AddChild(cast.NewTree("c")).
		Render()

	lines := strings.Split(got, "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d:\n%s", len(lines), got)
	}
	if lines[0] != "root" {
		t.Errorf("line 0: expected \"root\", got %q", lines[0])
	}
	// First two children use branch connector.
	if !strings.Contains(lines[1], "├") || !strings.Contains(lines[1], "a") {
		t.Errorf("line 1: expected branch+\"a\", got %q", lines[1])
	}
	if !strings.Contains(lines[2], "├") || !strings.Contains(lines[2], "b") {
		t.Errorf("line 2: expected branch+\"b\", got %q", lines[2])
	}
	// Last child uses last connector.
	if !strings.Contains(lines[3], "└") || !strings.Contains(lines[3], "c") {
		t.Errorf("line 3: expected last+\"c\", got %q", lines[3])
	}
}

func TestTree_FlatSingleChild(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTree("root").
		AddChild(cast.NewTree("only")).
		Render()

	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d:\n%s", len(lines), got)
	}
	// Only child → last connector.
	if !strings.Contains(lines[1], "└") || !strings.Contains(lines[1], "only") {
		t.Errorf("line 1: expected last+\"only\", got %q", lines[1])
	}
}

// ---------------------------------------------------------------------------
// Deep nesting
// ---------------------------------------------------------------------------

func TestTree_DeepNesting(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// forge/
	// └── core/
	//     └── forge.go
	got := cast.NewTree("forge/").
		AddChild(
			cast.NewTree("core/").
				AddChild(cast.NewTree("forge.go")),
		).
		Render()

	if !strings.Contains(got, "forge/") {
		t.Errorf("expected \"forge/\" in output:\n%s", got)
	}
	if !strings.Contains(got, "core/") {
		t.Errorf("expected \"core/\" in output:\n%s", got)
	}
	if !strings.Contains(got, "forge.go") {
		t.Errorf("expected \"forge.go\" in output:\n%s", got)
	}

	lines := strings.Split(got, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines for deep nesting, got %d:\n%s", len(lines), got)
	}
}

func TestTree_FullExample(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// forge/
	// ├── core/
	// │   ├── forge.go
	// │   └── renderer.go
	// ├── cast/
	// │   └── table.go
	// └── ink/
	//     └── style.go
	got := cast.NewTree("forge/").
		AddChild(
			cast.NewTree("core/").
				AddChild(cast.NewTree("forge.go")).
				AddChild(cast.NewTree("renderer.go")),
		).
		AddChild(
			cast.NewTree("cast/").
				AddChild(cast.NewTree("table.go")),
		).
		AddChild(
			cast.NewTree("ink/").
				AddChild(cast.NewTree("style.go")),
		).
		Render()

	for _, want := range []string{
		"forge/", "core/", "forge.go", "renderer.go",
		"cast/", "table.go", "ink/", "style.go",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q in output:\n%s", want, got)
		}
	}

	lines := strings.Split(got, "\n")
	// root + 3 dirs + 4 files = 8 lines
	if len(lines) != 8 {
		t.Fatalf("expected 8 lines, got %d:\n%s", len(lines), got)
	}

	// Root line
	if lines[0] != "forge/" {
		t.Errorf("line 0: expected \"forge/\", got %q", lines[0])
	}
	// core/ is first child → branch connector
	if !strings.Contains(lines[1], "├") {
		t.Errorf("line 1 (core/): expected branch connector ├:\n%s", got)
	}
	// ink/ is last child → last connector
	lastLine := lines[len(lines)-2] // ink/ is second-to-last overall
	if !strings.Contains(lastLine, "└") {
		t.Errorf("ink/ line: expected last connector └, got %q", lastLine)
	}
}

func TestTree_MixedDepth(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// a
	// ├── b
	// │   ├── d
	// │   └── e
	// └── c
	got := cast.NewTree("a").
		AddChild(
			cast.NewTree("b").
				AddChild(cast.NewTree("d")).
				AddChild(cast.NewTree("e")),
		).
		AddChild(cast.NewTree("c")).
		Render()

	for _, want := range []string{"a", "b", "c", "d", "e"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q in output:\n%s", want, got)
		}
	}

	lines := strings.Split(got, "\n")
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines, got %d:\n%s", len(lines), got)
	}
}

// ---------------------------------------------------------------------------
// AddChildren (variadic)
// ---------------------------------------------------------------------------

func TestTree_AddChildren(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTree("root").
		AddChildren(
			cast.NewTree("x"),
			cast.NewTree("y"),
			cast.NewTree("z"),
		).
		Render()

	lines := strings.Split(got, "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d:\n%s", len(lines), got)
	}
	if !strings.Contains(lines[3], "z") {
		t.Errorf("last child \"z\" not found on last line: %q", lines[3])
	}
}

// ---------------------------------------------------------------------------
// Connectors
// ---------------------------------------------------------------------------

func TestTree_DefaultConnectors(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTree("r").
		AddChild(cast.NewTree("a")).
		AddChild(cast.NewTree("b")).
		Render()

	if !strings.Contains(got, "├── ") {
		t.Errorf("expected default branch connector \"├── \" in:\n%s", got)
	}
	if !strings.Contains(got, "└── ") {
		t.Errorf("expected default last connector \"└── \" in:\n%s", got)
	}
}

func TestTree_CustomConnectors(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	got := cast.NewTree("root").
		WithConnectors("+-- ", "\\-- ", "|   ", "    ").
		AddChild(cast.NewTree("a")).
		AddChild(cast.NewTree("b")).
		Render()

	if !strings.Contains(got, "+-- ") {
		t.Errorf("expected custom branch connector \"+-- \" in:\n%s", got)
	}
	if !strings.Contains(got, "\\-- ") {
		t.Errorf("expected custom last connector \"\\-- \" in:\n%s", got)
	}
}

func TestTree_CustomConnectorsEmpty_FallbackToDefault(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Empty strings should fall back to defaults.
	got := cast.NewTree("root").
		WithConnectors("", "", "", "").
		AddChild(cast.NewTree("a")).
		Render()

	if !strings.Contains(got, "└── ") {
		t.Errorf("expected default last connector on empty WithConnectors, got:\n%s", got)
	}
}

// ---------------------------------------------------------------------------
// Pipe continuation in deep trees
// ---------------------------------------------------------------------------

func TestTree_PipeContinuation(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// When a non-last child has children, subsequent siblings at the parent
	// level should be connected via a pipe "│" continuation line.
	//
	// root
	// ├── a
	// │   └── a1
	// └── b
	got := cast.NewTree("root").
		AddChild(
			cast.NewTree("a").AddChild(cast.NewTree("a1")),
		).
		AddChild(cast.NewTree("b")).
		Render()

	if !strings.Contains(got, "│") {
		t.Errorf("expected pipe continuation \"│\" in:\n%s", got)
	}

	lines := strings.Split(got, "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d:\n%s", len(lines), got)
	}
}

// ---------------------------------------------------------------------------
// Styles
// ---------------------------------------------------------------------------

func TestTree_NodeStyle(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	s := ink.New().WithForeground(ink.Cyan)
	got := cast.NewTree("root").
		NodeStyle(s).
		AddChild(cast.NewTree("child")).
		Render()

	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI sequences from NodeStyle, got:\n%s", got)
	}
	plain := ink.Strip(got)
	if !strings.Contains(plain, "root") || !strings.Contains(plain, "child") {
		t.Errorf("expected \"root\" and \"child\" in plain output:\n%s", plain)
	}
}

func TestTree_RootStyle(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	rootStyle := ink.New().WithForeground(ink.BrightWhite).WithBold(true)
	got := cast.NewTree("root").
		RootStyle(rootStyle).
		AddChild(cast.NewTree("child")).
		Render()

	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI in output with RootStyle, got:\n%s", got)
	}
	plain := ink.Strip(got)
	if !strings.Contains(plain, "root") {
		t.Errorf("expected \"root\" in plain output:\n%s", plain)
	}
}

func TestTree_LeafStyle(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	leafStyle := ink.New().WithForeground(ink.Muted)
	got := cast.NewTree("root").
		LeafStyle(leafStyle).
		AddChild(cast.NewTree("leaf")).
		Render()

	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI in output with LeafStyle, got:\n%s", got)
	}
	plain := ink.Strip(got)
	if !strings.Contains(plain, "leaf") {
		t.Errorf("expected \"leaf\" in plain output:\n%s", plain)
	}
}

func TestTree_ConnectorStyle(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	connStyle := ink.New().WithForeground(ink.DimGray)
	got := cast.NewTree("root").
		ConnectorStyle(connStyle).
		AddChild(cast.NewTree("a")).
		AddChild(cast.NewTree("b")).
		Render()

	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI in output with ConnectorStyle, got:\n%s", got)
	}
	plain := ink.Strip(got)
	if !strings.Contains(plain, "a") || !strings.Contains(plain, "b") {
		t.Errorf("expected children in plain output:\n%s", plain)
	}
}

func TestTree_NoColorMode(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	s := ink.New().WithForeground(ink.Danger)
	got := cast.NewTree("root").
		NodeStyle(s).
		ConnectorStyle(s).
		AddChild(cast.NewTree("child")).
		Render()

	if strings.Contains(got, "\x1b[") {
		t.Errorf("expected no ANSI in NeverColorMode, got:\n%s", got)
	}
	if !strings.Contains(got, "root") || !strings.Contains(got, "child") {
		t.Errorf("expected content in NeverColorMode output:\n%s", got)
	}
}

// ---------------------------------------------------------------------------
// Immutability
// ---------------------------------------------------------------------------

func TestTree_Immutability_AddChild(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	base := cast.NewTree("root")
	withChild := base.AddChild(cast.NewTree("child"))

	// base should have no children → single line
	baseLines := strings.Split(base.Render(), "\n")
	if len(baseLines) != 1 {
		t.Errorf("base tree was mutated: expected 1 line, got %d:\n%s", len(baseLines), base.Render())
	}

	// withChild should have 2 lines
	childLines := strings.Split(withChild.Render(), "\n")
	if len(childLines) != 2 {
		t.Errorf("withChild tree should have 2 lines, got %d:\n%s", len(childLines), withChild.Render())
	}
}

func TestTree_Immutability_NodeStyle(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	base := cast.NewTree("root").AddChild(cast.NewTree("leaf"))
	styled := base.NodeStyle(ink.New().WithForeground(ink.Cyan))

	// base should produce no ANSI (no style set).
	if strings.Contains(base.Render(), "\x1b[") {
		t.Errorf("base tree was mutated with NodeStyle")
	}
	// styled should produce ANSI.
	if !strings.Contains(styled.Render(), "\x1b[") {
		t.Errorf("styled tree missing ANSI sequences")
	}
}

func TestTree_Immutability_Connectors(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	base := cast.NewTree("r").AddChild(cast.NewTree("x"))
	custom := base.WithConnectors("+-- ", "\\-- ", "|   ", "    ")

	// base should still use default connectors
	if strings.Contains(base.Render(), "+-- ") {
		t.Errorf("base tree connectors were mutated")
	}
	// custom should use custom connectors
	if !strings.Contains(custom.Render(), "\\-- ") {
		t.Errorf("custom tree should use custom last connector")
	}
}

// ---------------------------------------------------------------------------
// Edge cases
// ---------------------------------------------------------------------------

func TestTree_EmptyLabel(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Empty label should not panic and should produce a single (empty) line.
	got := cast.NewTree("").Render()
	if got != "" {
		t.Errorf("expected empty string for empty-label node, got %q", got)
	}
}

func TestTree_VeryDeepNesting(t *testing.T) {
	ink.SetGlobalColorMode(ink.NeverColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// Build a 5-level deep chain: a → b → c → d → e
	tree := cast.NewTree("a").
		AddChild(
			cast.NewTree("b").
				AddChild(
					cast.NewTree("c").
						AddChild(
							cast.NewTree("d").
								AddChild(cast.NewTree("e")),
						),
				),
		)

	got := tree.Render()
	for _, label := range []string{"a", "b", "c", "d", "e"} {
		if !strings.Contains(got, label) {
			t.Errorf("expected %q in deep tree output:\n%s", label, got)
		}
	}

	lines := strings.Split(got, "\n")
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines for 5-deep chain, got %d:\n%s", len(lines), got)
	}
}

func TestTree_RootStyleNotInheritedBySubRoots(t *testing.T) {
	ink.SetGlobalColorMode(ink.AlwaysColorMode)
	defer ink.SetGlobalColorMode(ink.AutoColorMode)

	// RootStyle should only apply to the actual root label, not to
	// intermediate nodes that are themselves roots of sub-trees.
	rootStyle := ink.New().WithForeground(ink.BrightWhite).WithBold(true)
	nodeStyle := ink.New().WithForeground(ink.Gray)

	got := cast.NewTree("root").
		RootStyle(rootStyle).
		NodeStyle(nodeStyle).
		AddChild(
			cast.NewTree("sub-root").
				AddChild(cast.NewTree("leaf")),
		).
		Render()

	// Output should be renderable without panic.
	if got == "" {
		t.Error("expected non-empty output")
	}
	plain := ink.Strip(got)
	for _, want := range []string{"root", "sub-root", "leaf"} {
		if !strings.Contains(plain, want) {
			t.Errorf("expected %q in plain output:\n%s", want, plain)
		}
	}
}
