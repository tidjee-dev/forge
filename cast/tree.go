package cast

import (
	"strings"

	"github.com/tidjee-dev/forge/ink"
)

// Tree is a recursive hierarchical tree renderer. Each node has a label and
// zero or more children, which are also Tree values.
//
// Every setter returns a new Tree — the original is never mutated.
//
// Basic usage:
//
//	root := cast.NewTree("forge/").
//	    AddChild(
//	        cast.NewTree("core/").
//	            AddChild(cast.NewTree("forge.go")).
//	            AddChild(cast.NewTree("renderer.go")),
//	    ).
//	    AddChild(
//	        cast.NewTree("cast/").
//	            AddChild(cast.NewTree("table.go")),
//	    )
//
//	fmt.Println(root.Render())
type Tree struct {
	label          string
	children       []Tree
	nodeStyle      ink.Style
	rootStyle      ink.Style
	leafStyle      ink.Style
	connectorStyle ink.Style

	// connector glyphs
	branchGlyph string // e.g. "├── "
	lastGlyph   string // e.g. "└── "
	pipeGlyph   string // e.g. "│   "
	indentGlyph string // e.g. "    "
}

// default connector glyphs
const (
	defaultBranch = "├── "
	defaultLast   = "└── "
	defaultPipe   = "│   "
	defaultIndent = "    "
)

// NewTree returns a Tree with the given label, default connector glyphs, and
// no children or styling. Calling Render() on the result immediately produces
// valid output.
func NewTree(label string) Tree {
	return Tree{
		label:       label,
		branchGlyph: defaultBranch,
		lastGlyph:   defaultLast,
		pipeGlyph:   defaultPipe,
		indentGlyph: defaultIndent,
	}
}

// ---------------------------------------------------------------------------
// Setters
// ---------------------------------------------------------------------------

// AddChild appends a single child node and returns a new Tree.
func (t Tree) AddChild(child Tree) Tree {
	children := make([]Tree, len(t.children)+1)
	copy(children, t.children)
	children[len(t.children)] = child
	t.children = children
	return t
}

// AddChildren appends multiple child nodes and returns a new Tree.
func (t Tree) AddChildren(children ...Tree) Tree {
	all := make([]Tree, len(t.children)+len(children))
	copy(all, t.children)
	copy(all[len(t.children):], children)
	t.children = all
	return t
}

// NodeStyle sets the ink.Style applied to every node label (root, branch, and
// leaf). More specific styles (RootStyle, LeafStyle) take precedence over
// NodeStyle when set.
func (t Tree) NodeStyle(s ink.Style) Tree {
	t.nodeStyle = s
	return t
}

// RootStyle sets the ink.Style applied to the root node label only.
func (t Tree) RootStyle(s ink.Style) Tree {
	t.rootStyle = s
	return t
}

// LeafStyle sets the ink.Style applied to leaf node labels (nodes with no
// children).
func (t Tree) LeafStyle(s ink.Style) Tree {
	t.leafStyle = s
	return t
}

// ConnectorStyle sets the ink.Style applied to the connector glyphs
// ("├── ", "└── ", "│   ").
func (t Tree) ConnectorStyle(s ink.Style) Tree {
	t.connectorStyle = s
	return t
}

// WithConnectors overrides the four connector glyph strings used when drawing
// the tree:
//
//   - branch — prefix for non-last children (default "├── ")
//   - last   — prefix for the last child     (default "└── ")
//   - pipe   — continuation line             (default "│   ")
//   - indent — blank continuation line       (default "    ")
//
// Pass empty strings to reset individual glyphs to their defaults.
func (t Tree) WithConnectors(branch, last, pipe, indent string) Tree {
	if branch == "" {
		branch = defaultBranch
	}
	if last == "" {
		last = defaultLast
	}
	if pipe == "" {
		pipe = defaultPipe
	}
	if indent == "" {
		indent = defaultIndent
	}
	t.branchGlyph = branch
	t.lastGlyph = last
	t.pipeGlyph = pipe
	t.indentGlyph = indent
	return t
}

// ---------------------------------------------------------------------------
// Render
// ---------------------------------------------------------------------------

// Render returns the tree as a multi-line string. The root label appears on
// the first line; each child is indented and prefixed with connector glyphs.
//
// Style precedence for node labels (highest to lowest):
//  1. RootStyle  — root node only
//  2. LeafStyle  — leaf nodes only
//  3. NodeStyle  — all nodes
//  4. (unstyled) — zero-value Style
func (t Tree) Render() string {
	var sb strings.Builder
	// labelPrefix — the complete string written before this node's label.
	//               For the root it is empty; for deeper nodes it is the
	//               accumulated ancestor continuation + this level's connector.
	// childBase   — plain-text prefix that will be prepended before each
	//               child's connector on the next level. It represents the
	//               "pipe" or "indent" continuation glyphs from all ancestors.
	treeRenderNode(t, &sb, "", "", true)
	result := sb.String()
	// Trim the trailing newline left by the last WriteByte('\n').
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result
}

// treeRenderNode writes one node's label and recursively all its descendants.
//
// Parameters:
//
//   - node        — node to render
//   - sb          — output builder
//   - labelPrefix — complete prefix written before this node's own label
//     (already includes ancestor continuation + connector)
//   - childBase   — plain-text base used to construct each child's labelPrefix
//     on the next level; it is the plain continuation accumulated
//     from all ancestor levels above this node
//   - isRoot      — true for the very first call so RootStyle is applied
//
// Key invariant: labelPrefix may contain ANSI escape sequences (styled
// connector), while childBase is always plain text (no ANSI). This prevents
// styled connector glyphs from leaking onto grandchild lines.
func treeRenderNode(node Tree, sb *strings.Builder, labelPrefix string, childBase string, isRoot bool) {
	// Write this node's label line.
	label := treeStyledLabel(node, isRoot)
	sb.WriteString(labelPrefix)
	sb.WriteString(label)
	sb.WriteByte('\n')

	for i, child := range node.children {
		isLast := i == len(node.children)-1

		// Propagate glyphs and styles from the parent to the child.
		child = treeInheritFrom(child, node)

		// connector   — the glyph written before the child's label
		// continuation — the plain-text glyph that represents this level on
		//                all subsequent lines (replaces the connector visually)
		var connector string
		var continuation string

		if isLast {
			connector = node.lastGlyph
			continuation = node.indentGlyph
		} else {
			connector = node.branchGlyph
			continuation = node.pipeGlyph
		}

		// Apply the connector style. childBase is always plain, so styledConn
		// may have ANSI but that is fine — it only appears on the child's own
		// label line, not on grandchild lines.
		styledConn := node.connectorStyle.Render(connector)

		// childLabelPrefix — what appears before the child's own label:
		//   plain ancestor continuation  +  styled connector
		childLabelPrefix := childBase + styledConn

		// grandchildBase — what all lines of the child's subtree use as their
		// plain-text base before they add their own connectors:
		//   plain ancestor continuation  +  plain continuation for this level
		grandchildBase := childBase + continuation

		treeRenderNode(child, sb, childLabelPrefix, grandchildBase, false)
	}
}

// treeStyledLabel returns node.label with the appropriate style applied.
// Priority (highest to lowest): RootStyle (isRoot only) > LeafStyle (leaf
// only) > NodeStyle > unstyled.
func treeStyledLabel(node Tree, isRoot bool) string {
	isLeaf := len(node.children) == 0

	switch {
	case isRoot && isStyleSet(node.rootStyle):
		return node.rootStyle.Render(node.label)
	case isRoot && isStyleSet(node.nodeStyle):
		return node.nodeStyle.Render(node.label)
	case isLeaf && isStyleSet(node.leafStyle):
		return node.leafStyle.Render(node.label)
	case isStyleSet(node.nodeStyle):
		return node.nodeStyle.Render(node.label)
	default:
		return node.label
	}
}

// treeInheritFrom copies connector glyphs and styles from parent into child,
// but only for fields that are still at their zero/default values. This lets
// the whole tree share consistent glyphs and styles unless a subtree
// explicitly overrides them.
//
// Note: rootStyle is intentionally NOT inherited — it applies only to the
// actual root node, not to sub-tree roots.
func treeInheritFrom(child Tree, parent Tree) Tree {
	// Inherit connector glyphs only when the child still has the default values.
	if child.branchGlyph == defaultBranch {
		child.branchGlyph = parent.branchGlyph
	}
	if child.lastGlyph == defaultLast {
		child.lastGlyph = parent.lastGlyph
	}
	if child.pipeGlyph == defaultPipe {
		child.pipeGlyph = parent.pipeGlyph
	}
	if child.indentGlyph == defaultIndent {
		child.indentGlyph = parent.indentGlyph
	}

	// Inherit styles only if the child has none explicitly set.
	if !isStyleSet(child.nodeStyle) {
		child.nodeStyle = parent.nodeStyle
	}
	if !isStyleSet(child.leafStyle) {
		child.leafStyle = parent.leafStyle
	}
	if !isStyleSet(child.connectorStyle) {
		child.connectorStyle = parent.connectorStyle
	}

	return child
}
