package charmimage

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestRenderHalfBlockExactCell(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	red := color.RGBA{R: 0xff, A: 0xff}
	blue := color.RGBA{B: 0xff, A: 0xff}
	img.SetRGBA(0, 0, red)
	img.SetRGBA(1, 0, red)
	img.SetRGBA(0, 1, blue)
	img.SetRGBA(1, 1, blue)

	got := RenderHalfBlock(img, 1, 1)
	want := "\x1b[38;2;255;0;0m\x1b[48;2;0;0;255m▀\x1b[0m"
	if got != want {
		t.Fatalf("RenderHalfBlock() = %q, quer %q", got, want)
	}
	if !strings.HasSuffix(got, "\x1b[0m") {
		t.Fatalf("saída não termina com reset: %q", got)
	}
}

func TestRenderHalfBlockScalePreservesAspect(t *testing.T) {
	img := solidImage(100, 50, color.RGBA{R: 0x80, G: 0x40, B: 0x20, A: 0xff})
	out := RenderHalfBlock(img, 20, 0)

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

func TestRenderHalfBlockNaturalWidthCappedAt80(t *testing.T) {
	img := solidImage(200, 100, color.RGBA{R: 0xff, A: 0xff})
	out := RenderHalfBlock(img, 0, 0)

	lines := strings.Split(out, "\n")
	if len(lines) != 20 {
		t.Fatalf("altura = %d linhas, quer 20", len(lines))
	}
	if got := strings.Count(lines[0], "▀"); got != 80 {
		t.Fatalf("largura = %d células, quer 80", got)
	}
}

func TestRenderHalfBlockResetsEveryLine(t *testing.T) {
	img := solidImage(4, 4, color.RGBA{G: 0xff, A: 0xff})
	out := RenderHalfBlock(img, 2, 2)

	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		t.Fatalf("altura = %d linhas, quer 2", len(lines))
	}
	for i, line := range lines {
		if !strings.HasSuffix(line, "\x1b[0m") {
			t.Fatalf("linha %d não termina com reset: %q", i, line)
		}
	}
}

func TestRenderHalfBlockNegativeSizeDoesNotPanic(t *testing.T) {
	img := solidImage(10, 10, color.RGBA{R: 0xff, A: 0xff})
	out := RenderHalfBlock(img, -10, -20)
	if out == "" {
		t.Fatal("saída vazia para dimensões negativas")
	}
}

func TestRenderHalfBlockNil(t *testing.T) {
	if got := RenderHalfBlock(nil, 10, 10); got != "" {
		t.Fatalf("RenderHalfBlock(nil) = %q, quer vazio", got)
	}
}

func TestRenderHalfBlockTransparentBackground(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.SetNRGBA(0, 0, color.NRGBA{R: 0xff, A: 0x00})

	got := renderHalfBlock(img, 1, 1, backgroundRGBA(color.RGBA{B: 0xff, A: 0xff}))
	want := "\x1b[38;2;0;0;255m\x1b[48;2;0;0;255m▀\x1b[0m"
	if got != want {
		t.Fatalf("renderHalfBlock() = %q, quer %q", got, want)
	}
}

func TestRenderHalfBlockNonZeroBounds(t *testing.T) {
	img := image.NewRGBA(image.Rect(10, 20, 12, 22))
	img.SetRGBA(10, 20, color.RGBA{R: 0xff, A: 0xff})
	img.SetRGBA(11, 20, color.RGBA{R: 0xff, A: 0xff})
	img.SetRGBA(10, 21, color.RGBA{B: 0xff, A: 0xff})
	img.SetRGBA(11, 21, color.RGBA{B: 0xff, A: 0xff})

	got := RenderHalfBlock(img, 1, 1)
	want := "\x1b[38;2;255;0;0m\x1b[48;2;0;0;255m▀\x1b[0m"
	if got != want {
		t.Fatalf("RenderHalfBlock() com bounds deslocados = %q, quer %q", got, want)
	}
}
