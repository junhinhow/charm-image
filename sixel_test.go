package charmimage

import (
	"encoding/base64"
	"image"
	"image/color"
	"strings"
	"testing"
)

// decodeSixel extrai e decodifica o payload base64 de uma sequência Sixel.
func decodeSixel(t *testing.T, out string) string {
	t.Helper()

	const prefix = "\x1bPq"
	const suffix = "\x1b\\"
	if !strings.HasPrefix(out, prefix) {
		t.Fatalf("sequência não começa com DCS: %q", out)
	}
	if !strings.HasSuffix(out, suffix) {
		t.Fatalf("sequência não termina com ST: %q", out)
	}

	body := strings.TrimSuffix(strings.TrimPrefix(out, prefix), suffix)
	decoded, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		t.Fatalf("payload não é base64 válido: %v", err)
	}
	if len(decoded) == 0 {
		t.Fatal("payload decodificado vazio")
	}
	return string(decoded)
}

func TestRenderSixelSequence(t *testing.T) {
	img := solidImage(4, 4, color.RGBA{R: 0xff, A: 0xff})
	out, err := RenderSixel(img, 2, 2)
	if err != nil {
		t.Fatalf("RenderSixel() erro: %v", err)
	}

	payload := decodeSixel(t, out)
	if !strings.HasPrefix(payload, `"1;1;16;32`) {
		t.Fatalf("payload não começa com os atributos de raster esperados: %q", payload)
	}
	if !strings.Contains(payload, "#180;2;100;0;0") {
		t.Fatalf("payload não define a cor vermelha da paleta: %q", payload)
	}
	if !strings.Contains(payload, "!") {
		t.Fatalf("payload não usa run-length encoding: %q", payload)
	}
	if !strings.HasSuffix(payload, "-") {
		t.Fatalf("payload não encerra a última banda: %q", payload)
	}
}

func TestRenderSixelSinglePixel(t *testing.T) {
	img := solidImage(1, 1, color.RGBA{B: 0xff, A: 0xff})
	out, err := RenderSixel(img, 0, 0)
	if err != nil {
		t.Fatalf("RenderSixel() erro: %v", err)
	}

	payload := decodeSixel(t, out)
	if !strings.HasPrefix(payload, `"1;1;8;16`) {
		t.Fatalf("payload 1×1 inesperado: %q", payload)
	}
}

func TestRenderSixelNegativeSizeDoesNotPanic(t *testing.T) {
	img := solidImage(3, 3, color.RGBA{G: 0xff, A: 0xff})
	out, err := RenderSixel(img, -5, -5)
	if err != nil {
		t.Fatalf("RenderSixel() erro: %v", err)
	}
	decodeSixel(t, out)
}

func TestRenderSixelDeterministic(t *testing.T) {
	img := gradientImage(20, 10)
	first, err := RenderSixel(img, 5, 3)
	if err != nil {
		t.Fatalf("RenderSixel() erro: %v", err)
	}
	second, err := RenderSixel(img, 5, 3)
	if err != nil {
		t.Fatalf("RenderSixel() erro: %v", err)
	}
	if first != second {
		t.Fatal("RenderSixel() não é determinístico")
	}
}

func TestRenderSixelBackground(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.SetNRGBA(0, 0, color.NRGBA{R: 0xff, A: 0x00})

	out, err := renderSixel(img, 1, 1, backgroundRGBA(color.RGBA{B: 0xff, A: 0xff}))
	if err != nil {
		t.Fatalf("renderSixel() erro: %v", err)
	}

	payload := decodeSixel(t, out)
	if !strings.Contains(payload, "#5;2;0;0;100") {
		t.Fatalf("payload não usa a cor de fundo azul: %q", payload)
	}
}

func TestRenderSixelNil(t *testing.T) {
	if _, err := RenderSixel(nil, 10, 10); err == nil {
		t.Fatal("RenderSixel(nil) deveria devolver erro")
	}
}

func TestRenderSixelEmptyImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))
	if _, err := RenderSixel(img, 10, 10); err == nil {
		t.Fatal("RenderSixel(imagem vazia) deveria devolver erro")
	}
}

func TestBuildSixelPayloadRLE(t *testing.T) {
	pixels := make([]uint8, 16*6)
	var used [256]bool
	used[0] = true
	payload := buildSixelPayload(16, 6, pixels, used)

	if !strings.Contains(payload, "#0;2;0;0;0") {
		t.Fatalf("payload sem definição de paleta: %q", payload)
	}
	if !strings.Contains(payload, "!16") {
		t.Fatalf("payload sem RLE de 16 colunas: %q", payload)
	}
}
