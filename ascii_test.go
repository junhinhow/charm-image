package charmimage

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestAsciiArtEmptyRampUsesDefault(t *testing.T) {
	white := solidImage(1, 1, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	black := solidImage(1, 1, color.RGBA{A: 0xff})

	if got := AsciiArt(white, 1, 1, ""); got != "@" {
		t.Fatalf("AsciiArt(branco, rampa vazia) = %q, quer %q", got, "@")
	}
	if got := AsciiArt(black, 1, 1, ""); got != " " {
		t.Fatalf("AsciiArt(preto, rampa vazia) = %q, quer %q", got, " ")
	}
}

func TestAsciiArtCustomRamp(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.SetRGBA(0, 0, color.RGBA{A: 0xff})
	img.SetRGBA(1, 0, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})

	if got := AsciiArt(img, 2, 1, "#."); got != "#." {
		t.Fatalf("AsciiArt() = %q, quer %q", got, "#.")
	}
}

func TestAsciiArtSingleRuneRamp(t *testing.T) {
	img := gradientImage(8, 4)
	got := AsciiArt(img, 8, 4, "X")
	if got != strings.Repeat("X", 8)+"\n"+strings.Repeat("X", 8)+"\n"+
		strings.Repeat("X", 8)+"\n"+strings.Repeat("X", 8) {
		t.Fatalf("AsciiArt() com rampa de um caractere = %q", got)
	}
}

func TestAsciiArtWidthZero(t *testing.T) {
	img := solidImage(4, 2, color.RGBA{A: 0xff})
	got := AsciiArt(img, 0, 0, " .")
	lines := strings.Split(got, "\n")
	if len(lines) != 1 {
		t.Fatalf("altura = %d linhas, quer 1", len(lines))
	}
	if len([]rune(lines[0])) != 4 {
		t.Fatalf("largura = %d caracteres, quer 4", len([]rune(lines[0])))
	}
}

func TestAsciiArtNegativeSizeDoesNotPanic(t *testing.T) {
	img := solidImage(3, 3, color.RGBA{R: 0xff, A: 0xff})
	if got := AsciiArt(img, -1, -1, " .#"); got == "" {
		t.Fatal("saída vazia para dimensões negativas")
	}
}

func TestAsciiArtNil(t *testing.T) {
	if got := AsciiArt(nil, 10, 10, " .#"); got != "" {
		t.Fatalf("AsciiArt(nil) = %q, quer vazio", got)
	}
}

func TestAsciiArtEmptyImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))
	if got := AsciiArt(img, 10, 10, " .#"); got != "" {
		t.Fatalf("AsciiArt(imagem vazia) = %q, quer vazio", got)
	}
}

func TestAsciiArtHasNoANSI(t *testing.T) {
	img := gradientImage(20, 10)
	got := AsciiArt(img, 20, 10, "")
	if strings.Contains(got, "\x1b") {
		t.Fatalf("AsciiArt() contém sequência ANSI: %q", got)
	}
}

func TestAsciiArtDitherChangesOutput(t *testing.T) {
	img := gradientImage(16, 8)
	plain := asciiArt(img, 16, 8, defaultRamp, false)
	dithered := asciiArt(img, 16, 8, defaultRamp, true)
	if plain == dithered {
		t.Fatal("dithering não alterou a saída")
	}
}
