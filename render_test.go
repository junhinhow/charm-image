package charmimage

import (
	"errors"
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestRenderNil(t *testing.T) {
	_, err := Render(nil)
	if !errors.Is(err, errNilImage) {
		t.Fatalf("Render(nil) erro = %v, quer errNilImage", err)
	}
}

func TestRenderEmptyImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))
	if _, err := Render(img); !errors.Is(err, errEmptyImage) {
		t.Fatalf("Render(imagem vazia) erro = %v, quer errEmptyImage", err)
	}
}

func TestRenderRespectsProtocol(t *testing.T) {
	img := gradientImage(8, 4)

	tests := []struct {
		name     string
		protocol Protocol
		prefix   string
		suffix   string
		contains string
	}{
		{"halfblock", ProtocolHalfBlock, "", "\x1b[0m", "▀"},
		{"sixel", ProtocolSixel, "\x1b[1;1H\x1bPq", "\x1b\\", "\x1b[1;1H"},
		{"kitty", ProtocolKitty, "\x1b_Gf=100,a=T,", "\x1b\\", "\x1b_Gf=100"},
		{"iterm2", ProtocolITerm2, "\x1b]1337;File=inline=1;", "\a", "\x1b]1337"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := Render(img, WithProtocol(tt.protocol), WithWidth(4), WithHeight(2))
			if err != nil {
				t.Fatalf("Render() erro: %v", err)
			}
			if tt.prefix != "" && !strings.HasPrefix(out, tt.prefix) {
				t.Fatalf("saída não começa com %q: %q", tt.prefix, out)
			}
			if !strings.HasSuffix(out, tt.suffix) {
				t.Fatalf("saída não termina com %q: %q", tt.suffix, out)
			}
			if !strings.Contains(out, tt.contains) {
				t.Fatalf("saída não contém %q", tt.contains)
			}
		})
	}
}

func TestRenderForcedProtocolIgnoresNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TERM", "xterm-256color")

	out, err := Render(gradientImage(8, 4),
		WithProtocol(ProtocolHalfBlock), WithWidth(4), WithHeight(2))
	if err != nil {
		t.Fatalf("Render() erro: %v", err)
	}
	if !strings.Contains(out, "▀") {
		t.Fatalf("protocolo forçado não foi respeitado: %q", out)
	}
}

func TestRenderDetectsAsciiFallback(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("KITTY_WINDOW_ID", "")

	img := gradientImage(8, 4)
	out, err := Render(img)
	if err != nil {
		t.Fatalf("Render() erro: %v", err)
	}
	if strings.Contains(out, "\x1b") {
		t.Fatalf("fallback ASCII contém ANSI: %q", out)
	}
	if want := AsciiArt(img, 0, 0, ""); out != want {
		t.Fatalf("Render() = %q, quer %q", out, want)
	}
}

func TestRenderDither(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("KITTY_WINDOW_ID", "")

	img := gradientImage(16, 8)
	plain, err := Render(img, WithWidth(16), WithHeight(8), WithDither(false))
	if err != nil {
		t.Fatalf("Render() erro: %v", err)
	}
	dithered, err := Render(img, WithWidth(16), WithHeight(8), WithDither(true))
	if err != nil {
		t.Fatalf("Render() erro: %v", err)
	}
	if plain == dithered {
		t.Fatal("WithDither(true) não alterou a saída ASCII")
	}
}

func TestRenderUnknownProtocol(t *testing.T) {
	_, err := Render(gradientImage(4, 4), WithProtocol(Protocol(42)))
	if err == nil {
		t.Fatal("Render() com protocolo desconhecido deveria devolver erro")
	}
}

func TestRenderZeroAndNegativeWidth(t *testing.T) {
	img := gradientImage(100, 50)

	for _, width := range []int{0, -1, -100} {
		out, err := Render(img, WithProtocol(ProtocolHalfBlock), WithWidth(width))
		if err != nil {
			t.Fatalf("Render(Width(%d)) erro: %v", width, err)
		}
		if out == "" {
			t.Fatalf("Render(Width(%d)) devolveu vazio", width)
		}
	}
}

func TestRenderScaleThroughOptions(t *testing.T) {
	img := gradientImage(100, 50)
	out, err := Render(img, WithProtocol(ProtocolHalfBlock), WithWidth(20))
	if err != nil {
		t.Fatalf("Render() erro: %v", err)
	}

	lines := strings.Split(out, "\n")
	if len(lines) != 5 {
		t.Fatalf("altura = %d linhas, quer 5", len(lines))
	}
	for i, line := range lines {
		if got := strings.Count(line, "▀"); got != 20 {
			t.Fatalf("linha %d tem %d células, quer 20", i, got)
		}
	}
}

func TestRenderBackgroundOption(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.SetNRGBA(0, 0, color.NRGBA{R: 0xff, A: 0x00})

	out, err := Render(img,
		WithProtocol(ProtocolHalfBlock),
		WithWidth(1),
		WithHeight(1),
		withBackground(color.RGBA{B: 0xff, A: 0xff}),
	)
	if err != nil {
		t.Fatalf("Render() erro: %v", err)
	}

	want := "\x1b[38;2;0;0;255m\x1b[48;2;0;0;255m▀\x1b[0m"
	if out != want {
		t.Fatalf("Render() = %q, quer %q", out, want)
	}
}

func TestRenderAsciiBackgroundOption(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("KITTY_WINDOW_ID", "")

	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.SetNRGBA(0, 0, color.NRGBA{R: 0xff, A: 0x00})

	out, err := Render(img,
		WithWidth(1),
		WithHeight(1),
		withBackground(color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}),
	)
	if err != nil {
		t.Fatalf("Render() erro: %v", err)
	}
	if out != "@" {
		t.Fatalf("Render() = %q, quer %q", out, "@")
	}
}

func TestRenderDeterministic(t *testing.T) {
	img := gradientImage(24, 12)

	for _, protocol := range []Protocol{
		ProtocolHalfBlock,
		ProtocolSixel,
		ProtocolKitty,
		ProtocolITerm2,
	} {
		t.Run(protocol.String(), func(t *testing.T) {
			first, err := Render(img, WithProtocol(protocol), WithWidth(6))
			if err != nil {
				t.Fatalf("primeira renderização: %v", err)
			}
			second, err := Render(img, WithProtocol(protocol), WithWidth(6))
			if err != nil {
				t.Fatalf("segunda renderização: %v", err)
			}
			if first != second {
				t.Fatal("renderizações diferentes para a mesma entrada")
			}
		})
	}
}
