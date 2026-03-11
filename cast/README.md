# cast

Ready-made terminal display components for Go, built on top of [`ink`](../ink).
Every component renders to a plain `string` — no runtime, no state, no event loop.

## Install

```sh
go get github.com/tidjee-dev/forge/cast
```

## Quick start

```go
import (
    "fmt"
    "github.com/tidjee-dev/forge/cast"
    "github.com/tidjee-dev/forge/ink"
)

// Badge
fmt.Println(cast.NewBadge("OK").Success().Render())

// Divider with a centred label
fmt.Println(cast.NewDivider().Label("Results").Width(40).Render())

// Table
fmt.Println(
    cast.NewTable().
        Headers("Name", "Status", "Port").
        Border(ink.BorderRounded()).
        AddRow("api", "running", "8080").
        AddRow("db",  "stopped", "5432").
        Render(),
)

// Animated spinner
s := cast.NewSpinner().WithLabel("Loading…")
s.Start()
// … do work …
s.Stop()
```

---

## Components

- [Badge](#badge)
- [Banner](#banner)
- [Divider](#divider)
- [List](#list)
- [Spinner](#spinner)
- [Table](#table)
- [Tree](#tree)

---

## Badge

A short inline label with a styled background — useful for status indicators.

```
 OK    WARN    ERROR    INFO    v1.2.3
```

```go
// Semantic shorthands
cast.NewBadge("OK").Success().Render()
cast.NewBadge("WARN").Warning().Render()
cast.NewBadge("ERR").Danger().Render()
cast.NewBadge("INFO").Info().Render()
cast.NewBadge("v1.2.3").Neutral().Render()

// Custom style
cast.NewBadge("custom").
    Style(ink.New().WithBackground(ink.Blue).WithForeground(ink.White)).
    Render()
```

| Setter / Method   | Description                                                  |
| ----------------- | ------------------------------------------------------------ |
| `NewBadge(label)` | Create a badge with the given label                          |
| `Style(s)`        | Full style override                                          |
| `Success()`       | Bright-green background, black foreground, bold              |
| `Warning()`       | Amber background, black foreground, bold                     |
| `Danger()`        | Bright-red background, white foreground, bold                |
| `Info()`          | Bright-blue background, white foreground, bold               |
| `Neutral()`       | Muted background, white foreground, bold                     |
| `Render()`        | Returns the badge string, padded with one space on each side |

---

## Banner

A full-width styled block — useful for section headers or prominent messages.

```
╭─────────────────────────────╮
│  Server started on :8080    │
╰─────────────────────────────╯
```

```go
cast.NewBanner("Server started on :8080").
    Style(ink.New().WithForeground(ink.BrightWhite).WithBold(true)).
    Border(ink.BorderRounded()).
    Width(40).
    Align(ink.JustifyCenter).
    Render()
```

| Setter / Method   | Description                                                                |
| ----------------- | -------------------------------------------------------------------------- |
| `NewBanner(text)` | Create a banner with the given text                                        |
| `Style(s)`        | Text style (foreground, background, bold, …)                               |
| `Border(bs)`      | Border glyph set; use `ink.NoBorder()` to remove                           |
| `Width(n)`        | Fixed render width in terminal columns (0 = auto)                          |
| `Align(a)`        | Horizontal text alignment: `JustifyStart` / `JustifyCenter` / `JustifyEnd` |
| `Render()`        | Returns the banner string                                                  |

---

## Divider

A horizontal rule with an optional inline label.

```
────────────────────────────────────────────────────────────────────────────────
─────────────────────────────────── Results ────────────────────────────────────
─── Results ────────────────────────────────────────────────────────────────────
──────────────────────────────────────────────────────────────────── Results ───
```

```go
// Plain rule
cast.NewDivider().Render()

// Labelled — default alignment is center
cast.NewDivider().Label("Results").Width(60).Render()

// Label alignment
cast.NewDivider().Label("Start").LabelAlign(ink.JustifyStart).Width(60).Render()
cast.NewDivider().Label("End").LabelAlign(ink.JustifyEnd).Width(60).Render()

// Custom char and style
cast.NewDivider().
    Char("═").
    Style(ink.New().WithForeground(ink.Muted)).
    Width(60).
    Render()
```

| Setter / Method | Description                                                                          |
| --------------- | ------------------------------------------------------------------------------------ |
| `NewDivider()`  | Create a divider (`"─"` fill, 80-column default width)                               |
| `Label(text)`   | Inline label placed within the rule                                                  |
| `LabelAlign(a)` | Label position: `JustifyStart` / `JustifyCenter` (default) / `JustifyEnd`            |
| `NearFill(n)`   | Number of fill chars on the short side for `JustifyStart` / `JustifyEnd` (default 3) |
| `Char(s)`       | Fill character (default `"─"`)                                                       |
| `Width(n)`      | Fixed render width in terminal columns (default 80)                                  |
| `Style(s)`      | Style applied to the fill characters                                                 |
| `LabelStyle(s)` | Style applied to the label only; falls back to `Style` when unset                    |
| `Render()`      | Returns the divider string                                                           |

---

## List

A vertical list of items with an optional bullet or sequential numbering.

```
• Buy carrots        1. Buy carrots
• Buy celery         2. Buy celery
• Buy kohlrabi       3. Buy kohlrabi
```

```go
// Bulleted
cast.NewList("Buy carrots", "Buy celery", "Buy kohlrabi").
    Bullet("•").
    ItemStyle(ink.New().WithForeground(ink.White)).
    BulletStyle(ink.New().WithForeground(ink.Muted)).
    Render()

// Numbered
cast.NewList("one", "two", "three").Numbered().Render()

// No bullet
cast.NewList("line one", "line two").Bullet("").Render()

// With indent
cast.NewList("alpha", "beta").Indent(4).Render()
```

| Setter / Method   | Description                                                      |
| ----------------- | ---------------------------------------------------------------- |
| `NewList(items…)` | Create a list pre-populated with items; default bullet `"•"`     |
| `AddItem(item)`   | Append a single item                                             |
| `Bullet(s)`       | Custom bullet character; `""` suppresses it                      |
| `Numbered()`      | Use `1.`, `2.`, … (takes precedence over `Bullet`)               |
| `ItemStyle(s)`    | Style applied to each item's text                                |
| `BulletStyle(s)`  | Style applied to the bullet or number; falls back to `ItemStyle` |
| `Indent(n)`       | Leading spaces prepended to every line                           |
| `Render()`        | Returns the list as a newline-separated string                   |

---

## Spinner

A self-animating spinner that runs in a background goroutine, advancing through
its frame set at a fixed interval and overwriting the same terminal line with
each tick.

```
⠋ Loading…   ⠙ Loading…   ⠹ Loading…   …
```

```go
s := cast.NewSpinner().
    WithFrames(cast.SpinnerDots).
    WithLabel("Loading…").
    WithStyle(ink.New().WithForeground(ink.Info)).
    WithInterval(80 * time.Millisecond)

s.Start()
defer s.Stop() // clears the line and blocks until the goroutine exits
```

### Built-in frame sets

| Constant        | Frames                            |
| --------------- | --------------------------------- |
| `SpinnerDots`   | `⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏` (10 frames) |
| `SpinnerLine`   | `- \ \| /` (4 frames)             |
| `SpinnerCircle` | `◐ ◓ ◑ ◒` (4 frames)              |
| `SpinnerArrow`  | `← ↖ ↑ ↗ → ↘ ↓ ↙` (8 frames)      |
| `SpinnerBounce` | `⣾ ⣽ ⣻ ⢿ ⡿ ⣟ ⣯ ⣷` (8 frames)      |

### API

| Setter / Method     | Description                                                                       |
| ------------------- | --------------------------------------------------------------------------------- |
| `NewSpinner()`      | Create a spinner (`SpinnerDots`, 80 ms interval, output to `os.Stderr`)           |
| `WithFrames(f)`     | Replace the frame set; ignored if `f` is empty                                    |
| `WithLabel(text)`   | Text displayed alongside the glyph                                                |
| `WithStyle(s)`      | Style applied to the glyph (and label when no `LabelStyle` is set)                |
| `WithLabelStyle(s)` | Style applied to the label only; falls back to `WithStyle`                        |
| `WithInterval(d)`   | Frame advance interval; values ≤ 0 are ignored                                    |
| `WithWriter(w)`     | Redirect output (default `os.Stderr`); useful in tests                            |
| `Start()`           | Start the animation goroutine; no-op if already running                           |
| `Stop()`            | Stop animation, clear the line, block until goroutine exits; no-op if not running |
| `IsRunning()`       | Report whether the goroutine is currently active                                  |

`DefaultSpinnerInterval` is exported as a `time.Duration` constant (`80ms`).

---

## Table

A bordered, aligned tabular data renderer.

```
╭──────┬─────────┬──────╮
│ Name │ Status  │ Port │
├──────┼─────────┼──────┤
│ api  │ running │ 8080 │
│ db   │ stopped │ 5432 │
╰──────┴─────────┴──────╯
```

```go
cast.NewTable().
    Headers("Name", "Status", "Port").
    Border(ink.BorderRounded()).
    HeaderStyle(ink.New().WithBold(true).WithForeground(ink.BrightWhite)).
    AltRowStyle(ink.New().WithDim(true)).
    ColumnAlign(2, ink.JustifyEnd).
    AddRow("api", "running", "8080").
    AddRow("db",  "stopped", "5432").
    Render()
```

| Setter / Method          | Description                                                                 |
| ------------------------ | --------------------------------------------------------------------------- |
| `NewTable()`             | Create an empty table with no headers, rows, or border                      |
| `Headers(cols…)`         | Set column headers                                                          |
| `AddRow(cells…)`         | Append a data row                                                           |
| `AddRows(rows)`          | Append multiple data rows at once                                           |
| `Border(bs)`             | Border glyph set; use `ink.NoBorder()` to remove                            |
| `NoBorder()`             | Remove any border                                                           |
| `HeaderStyle(s)`         | Style applied to every header cell                                          |
| `RowStyle(s)`            | Style applied to all data rows                                              |
| `AltRowStyle(s)`         | Style applied to alternating rows (zebra striping)                          |
| `ColumnAlign(col, a)`    | Alignment for column `col`: `JustifyStart` / `JustifyCenter` / `JustifyEnd` |
| `ColumnStyle(col, s)`    | Style for all cells in column `col`; overrides row style                    |
| `ColumnMinWidth(col, n)` | Minimum column width in terminal columns                                    |
| `Render()`               | Returns the table string; empty string when no headers or rows              |

Column widths are auto-sized to the widest cell. ANSI sequences are stripped
before measuring so styled input never inflates column widths.

---

## Tree

A hierarchical tree renderer.

```
forge/
├── core/
│   ├── forge.go
│   └── renderer.go
├── cast/
│   └── table.go
└── ink/
    └── style.go
```

```go
root := cast.NewTree("forge/").
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
    )

fmt.Println(root.Render())
```

### Custom connectors

```go
root := cast.NewTree("root").
    WithConnectors("├── ", "└── ", "│   ", "    ").
    AddChild(cast.NewTree("child"))
```

### Styling

```go
root := cast.NewTree("forge/").
    RootStyle(ink.New().WithForeground(ink.BrightWhite).WithBold(true)).
    NodeStyle(ink.New().WithForeground(ink.White)).
    LeafStyle(ink.New().WithForeground(ink.Muted)).
    ConnectorStyle(ink.New().WithForeground(ink.Muted))
```

Style precedence (highest to lowest): `RootStyle` → `LeafStyle` → `NodeStyle` → unstyled.
Styles and connector glyphs are inherited by children that have none explicitly set.

| Setter / Method              | Description                                                            |
| ---------------------------- | ---------------------------------------------------------------------- |
| `NewTree(label)`             | Create a tree node with default connectors and no children             |
| `AddChild(child)`            | Append a single child node                                             |
| `AddChildren(children…)`     | Append multiple child nodes                                            |
| `NodeStyle(s)`               | Style applied to all node labels                                       |
| `RootStyle(s)`               | Style applied to the root label only                                   |
| `LeafStyle(s)`               | Style applied to leaf node labels                                      |
| `ConnectorStyle(s)`          | Style applied to connector glyphs (`├── `, `└── `, `│   `)             |
| `WithConnectors(b, l, p, i)` | Override branch, last, pipe, and indent glyphs; `""` resets to default |
| `Render()`                   | Returns the tree as a multi-line string                                |

---

## Design notes

**Immutability** — every setter on every component returns a new copy; the
original is never mutated. This matches the `ink.Style` pattern and makes
components safe to store, share, and chain freely.

**ANSI-aware sizing** — all width calculations strip ANSI sequences before
measuring, so styled input never inflates padding or column widths.

**Unicode width** — full-width and combining characters are handled correctly;
CJK and emoji glyphs that occupy two terminal columns are measured as 2.

**Color fallback** — all components delegate to `ink.Style.Render()`, which
already respects `NO_COLOR`, `TERM=dumb`, `COLORTERM`, and TTY detection. No
additional color logic is needed at the `cast` level.

---

## License

Part of the [forge](https://github.com/tidjee-dev/forge) project.
