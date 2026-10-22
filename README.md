# charm-image

Render images in the terminal using the best protocol the terminal supports — Kitty, iTerm2, Sixel, or a universal truecolor half-block fallback.

Part of [Charm Suite](https://github.com/junhinhow/charm-suite).

## Install

```bash
go get github.com/junhinhow/charm-image
```

No external dependencies: only the Go standard library.

## Usage

```go
package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	charmimage "github.com/junhinhow/charm-image"
)

func main() {
	file, err := os.Open("photo.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	out, err := charmimage.Render(img, charmimage.WithWidth(60))
	if err != nil {
		panic(err)
	}
	fmt.Print(out)
}
```

`Render` picks the protocol from the environment. To force one:

```go
out, err := charmimage.Render(img,
	charmimage.WithProtocol(charmimage.ProtocolHalfBlock),
	charmimage.WithWidth(40),
	charmimage.WithHeight(20),
)
```

For terminals without any color support, print `AsciiArt` instead:

```go
fmt.Print(charmimage.AsciiArt(img, 40, 0, ""))
```

## Terminal support

| Terminal | Protocol | How it is detected |
| --- | --- | --- |
| Kitty | Kitty | `KITTY_WINDOW_ID` is set |
| WezTerm | iTerm2 | `TERM_PROGRAM=WezTerm` |
| iTerm2 | iTerm2 | `TERM_PROGRAM=iTerm.app` |
| mintty | iTerm2 | `TERM_PROGRAM=mintty` |
| mlterm, xterm with sixel, tmux with sixel features | Sixel | `TERM` or `COLORTERM` contains `sixel` |
| Most modern terminals | Half-block | fallback (truecolor) |
| `TERM=dumb` or `NO_COLOR` | ASCII | no escape sequences; `Render` uses `AsciiArt` |

Detection is heuristic and based only on environment variables; no terminal query is sent. `WithProtocol` always wins over detection.

## API

| Symbol | Description |
| --- | --- |
| `Protocol`, `ProtocolNone`, `ProtocolHalfBlock`, `ProtocolSixel`, `ProtocolKitty`, `ProtocolITerm2` | Protocol enumeration |
| `Detect()`, `DetectFrom(getenv)` | Environment-based protocol detection |
| `Options`, `Option` | Render configuration and functional options |
| `WithWidth(n)`, `WithHeight(n)`, `WithProtocol(p)`, `WithDither(b)` | Option constructors |
| `Render(img, opts...)` | Renders with the detected or forced protocol |
| `RenderHalfBlock(img, width, height)` | Universal truecolor fallback using `▀` |
| `RenderSixel`, `RenderKitty`, `RenderITerm2` | Raw escape sequences for each protocol |
| `AsciiArt(img, width, height, ramp)` | Colorless fallback with a configurable ramp |

## Design notes

- **Package name:** the package is `charmimage`, not `image`, so it never shadows the standard library.
- **Aspect ratio:** a terminal cell is assumed to be roughly twice as tall as it is wide. When only `Width` is given, the height is derived with that ratio.
- **Natural width:** `width <= 0` uses the image width capped at 80 cells. Negative values never reach a slice or `strings.Repeat`.
- **Half-block:** each cell renders two pixel rows with `▀`; the top pixel is the foreground and the bottom pixel is the background, both in truecolor (`38;2` / `48;2`). Every row ends with `ESC[0m`.
- **Kitty / iTerm2:** the full-resolution image is encoded as PNG (stdlib) and sent as base64. Kitty receives `c`/`r` in cells; iTerm2 receives `width`/`height` in cells. The terminal scales while preserving the aspect ratio.
- **Sixel:** rasterized at 8×16 pixels per cell, quantized to a fixed 216-color cube plus 40 grays, and run-length encoded. The Sixel payload is wrapped in base64 inside the DCS envelope (`ESC P q … ESC \`), as required by this library's contract.
- **Cursor:** `Render` emits `ESC[1;1H` before a Sixel block. `RenderSixel` returns only the DCS sequence, so it can be composed by the caller.
- **Dithering:** `WithDither(true)` applies an ordered Bayer 4×4 offset to the luminance in the ASCII fallback. Truecolor protocols do not quantize colors, so they are unaffected.
- **Transparency:** in the half-block, ASCII and Sixel paths, pixels are composited over `Options.Background`, or black when it is nil. Kitty and iTerm2 send the PNG alpha channel and let the terminal compose it. There is no `WithBackground`; use a custom option:
  ```go
  charmimage.Option(func(o *charmimage.Options) {
      o.Background = color.White
  })
  ```
- **Determinism:** rendering never iterates over maps; the same image and options always produce the same string.

## Limitations

- The escape sequences are covered by unit tests, but this project was not exercised against real Kitty, iTerm2, WezTerm, mintty, or Sixel terminals in its development environment.
- Sixel data is base64-encoded inside the DCS envelope. Stock Sixel terminals expect unencoded Sixel data; decode the payload first if you need to feed a terminal directly.
- Detection does not query the terminal (no DA1/DA2), so unusual setups may need `WithProtocol`.
- Half-block output assumes truecolor support; terminals limited to 256 colors will approximate the colors.
- The half-block fallback averages source pixels per cell. There is no gamma-correct resampling or image rotation.

## License

MIT
