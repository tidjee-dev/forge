# ink

A terminal styling library for Go — colors, text attributes, borders, layout, and themes, all through a fluent API.

## Install

```sh
go get github.com/tidjee-dev/forge/ink
```

## Quick start

```go
import "github.com/tidjee-dev/forge/ink"

// Bold cyan text
s := ink.New().WithForeground(ink.Cyan).WithBold(true)
fmt.Println(s.Render("Hello, world!"))

// Danger label with a rounded border
label := ink.New().
    WithForeground(ink.Danger).
    WithBorder(ink.BorderRounded()).
    WithBorderColor(ink.Danger)
fmt.Println(label.Render("error"))
```

---

## Features

- **True-color, ANSI 256, and ANSI 16** color support
- **Named color palette** — 100+ ready-to-use constants (`Red`, `Cyan`, `Emerald`, `RoyalBlue`, …) plus semantic aliases (`Success`, `Warning`, `Danger`, `Info`, `Muted`)
- **Text attributes** — bold, italic, underline, dim, strikethrough, inverse
- **Borders** — 8 built-in presets, per-side control, custom border colors
- **Layout** — padding, min/max width, horizontal alignment
- **Themes** — named style maps, concurrent-safe, cloneable and mergeable
- **Style overrides** — patch a base style without mutating it
- **ANSI strip** — remove escape sequences from any string
- **Color utilities** — `Lighten`, `Darken`, `Mix`, `ContrastRatio`, `ContrastedColor`, `Hex`
- **Auto color detection** — respects `NO_COLOR`, `COLORTERM`, `TERM=dumb`, and TTY state

---

## Color

Three ways to specify a color:

```go
ink.RGB(255, 87, 51)          // 24-bit true color
ink.Hex("#ff5733")       // CSS hex (long or short form)
ink.ANSI256(202)              // xterm 256-color palette
ink.ANSI16(1)                 // classic 16-color palette
```

Named constants are available directly:

```go
ink.Red    ink.Green    ink.Blue    ink.Cyan
ink.White  ink.Black    ink.Muted   ink.Danger
// …and 100+ more
```

### Color utilities

```go
// Adjust brightness
brighter := ink.Lighten(ink.Red, 0.3)   // 30% toward white
darker   := ink.Darken(ink.Blue, 0.2)   // 20% toward black

// Blend two colors
mid := ink.Mix(ink.Red, ink.Blue, 0.5)  // midpoint

// Accessibility
ratio := ink.ContrastRatio(ink.White, ink.Navy)   // WCAG contrast ratio
fg    := ink.ContrastedColor(ink.RoyalBlue)        // white or black for best legibility
```

---

## Style

`Style` is a value type. Every setter returns a new copy — the original is never mutated.

```go
s := ink.New().
    WithForeground(ink.Cyan).
    WithBackground(ink.RGB(18, 18, 30)).
    WithBold(true).
    WithItalic(false).
    WithUnderline(true)

fmt.Println(s.Render("styled text"))
```

Available setters:

| Setter                         | Description                          |
| ------------------------------ | ------------------------------------ |
| `WithForeground(Color)`        | Foreground (text) color              |
| `WithBackground(Color)`        | Background color                     |
| `WithBold(bool)`               | Bold weight                          |
| `WithItalic(bool)`             | Italic                               |
| `WithUnderline(bool)`          | Underline                            |
| `WithDim(bool)`                | Dim (reduced intensity)              |
| `WithStrikethrough(bool)`      | Strikethrough                        |
| `WithInverse(bool)`            | Swap foreground / background         |
| `WithLayout(Layout)`           | Padding, alignment, size constraints |
| `WithBorder(BorderStyle)`      | Border on all four sides             |
| `WithBorderStyle(BorderStyle)` | Border glyph set                     |
| `WithBorderColor(Color)`       | Border glyph color                   |
| `WithBorderSide(BorderSide)`   | Select which sides to draw           |

---

## Borders

### Built-in presets

```
BorderNormal()       BorderRounded()      BorderThick()
┌─────┐              ╭─────╮              ┏━━━━━┓
│  …  │              │  …  │              ┃  …  ┃
└─────┘              ╰─────╯              ┗━━━━━┛

BorderDouble()       BorderDashed()       BorderASCII()
╔═════╗              ┌┄┄┄┄┄┐              +-----+
║  …  ║              ┆  …  ┆              |  …  |
╚═════╝              └┄┄┄┄┄┘              +-----+

BorderBlock()        BorderHidden()       BorderInnerHalfBlock()
███████              (spaces)             ▄▄▄▄▄▄▄
█  …  █                 …                 ▌  …  ▐
███████              (spaces)             ▀▀▀▀▀▀▀
```

### Selective sides

```go
s := ink.New().
    WithBorderStyle(ink.BorderRounded()).
    WithBorderSide(ink.BorderSideTop | ink.BorderSideBottom) // top and bottom only
```

Available side constants:

```go
ink.BorderSideNone        // no border
ink.BorderSideTop
ink.BorderSideRight
ink.BorderSideBottom
ink.BorderSideLeft
ink.BorderSideAll         // all four sides (default)
ink.BorderSideVertical    // top + bottom
ink.BorderSideHorizontal  // left + right
```

### Custom border

```go
custom := ink.BorderStyle{
    TopLeft: "★", TopRight: "★", BottomLeft: "★", BottomRight: "★",
    Top: "─", Bottom: "─", Left: "│", Right: "│",
}
s := ink.New().WithBorder(custom).WithBorderColor(ink.Gold)
fmt.Println(s.Render("custom border"))
```

---

## Layout

`Layout` controls padding, horizontal alignment, and width constraints.

```go
l := ink.NewLayout().
    WithUniformPadding(1).          // 1 cell on all sides
    WithMinWidth(30).               // at least 30 columns
    WithJustifyContent(ink.JustifyCenter)

s := ink.New().WithLayout(l)
fmt.Println(s.Render("centered"))
```

### Padding

```go
ink.NewLayout().WithUniformPadding(1)           // all sides
ink.NewLayout().WithPaddingX(2)                 // left + right
ink.NewLayout().WithPaddingY(1)                 // top + bottom
ink.NewLayout().WithPadding(1, 2, 1, 2)         // top, right, bottom, left
ink.NewLayout().WithPaddingTop(1).WithPaddingLeft(2)
```

### Alignment (`JustifyContent`)

```go
ink.JustifyStart        // left-align (default)
ink.JustifyCenter       // center
ink.JustifyEnd          // right-align
```

### Width constraints

```go
ink.NewLayout().WithMinWidth(20)
ink.NewLayout().WithMaxWidth(80)   // truncates with … when exceeded
ink.NewLayout().WithWidth(40)      // sets both min and max
```

---

## Theme

`Theme` is a concurrent-safe, named map of styles.

```go
t := ink.NewTheme()
t.Set("title",   ink.New().WithForeground(ink.White).WithBold(true))
t.Set("success", ink.New().WithForeground(ink.Success).WithBold(true))
t.Set("error",   ink.New().WithForeground(ink.Danger).WithBold(true))
t.Set("muted",   ink.New().WithForeground(ink.Muted).WithDim(true))

fmt.Println(t.Render("title",   "Build complete"))
fmt.Println(t.Render("success", "All tests passed"))
fmt.Println(t.Render("error",   "Compilation failed"))
```

### Clone and merge

```go
// Deep-copy a theme; changes to the clone don't affect the original
dark := base.Clone()
dark.Set("title", ink.New().WithForeground(ink.BrightWhite).WithBold(true))

// Overlay one theme on top of another
base.Merge(overrides)
```

---

## Style overrides

`StyleOverride` lets you patch a base style — including explicitly setting a flag to `false` — without mutating the original.

```go
base := ink.New().WithForeground(ink.Red).WithBold(true)

patch := ink.Override().
    NoBold().
    Foreground(ink.Blue)

result := base.Merge(patch)  // not bold, blue fg
```

---

## Strip

Remove all ANSI escape sequences from a string:

```go
plain := ink.Strip(styled) // safe to measure length, write to files, etc.
```

`Strip` handles SGR sequences, OSC sequences, and single-character escapes. Calling it on a string that contains no `ESC` byte is a zero-allocation fast path.

---

## Color mode

Color output is auto-detected from the environment and TTY state. Override it globally if needed:

```go
ink.SetGlobalColorMode(ink.AlwaysColorMode) // force colors on (e.g. CI pipelines)
ink.SetGlobalColorMode(ink.NeverColorMode)  // force plain text
ink.SetGlobalColorMode(ink.AutoColorMode)   // default — detect automatically
```

The auto mode respects:

- `NO_COLOR` env var (disables color)
- `TERM=dumb` (disables color)
- `COLORTERM=truecolor` / `COLORTERM=24bit` (forces color)
- Whether stdout is a TTY

---

## License

Part of the [forge](https://github.com/tidjee-dev/forge) project.
