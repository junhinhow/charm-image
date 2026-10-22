package charmimage

import (
	"encoding/base64"
	"image"
	"strings"
	"testing"
)

func TestRenderITerm2Sequence(t *testing.T) {
	img := solidImage(4, 3, solidColor())
	out, err := RenderITerm2(img, 2, 3)
	if err != nil {
		t.Fatalf("RenderITerm2() erro: %v", err)
	}

	const prefix = "\x1b]1337;File=inline=1;width=2;height=3:"
	if !strings.HasPrefix(out, prefix) {
		t.Fatalf("sequência não começa com %q: %q", prefix, out)
	}
	if !strings.HasSuffix(out, "\a") {
		t.Fatalf("sequência não termina com BEL: %q", out)
	}

	payload := strings.TrimSuffix(strings.TrimPrefix(out, prefix), "\a")
	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		t.Fatalf("payload não é base64 válido: %v", err)
	}

	decoded := decodePNG(t, data)
	bounds := decoded.Bounds()
	if bounds.Dx() != 4 || bounds.Dy() != 3 {
		t.Fatalf("PNG decodificado tem %dx%d, quer 4x3", bounds.Dx(), bounds.Dy())
	}
}

func TestRenderITerm2NaturalSize(t *testing.T) {
	img := solidImage(100, 50, solidColor())
	out, err := RenderITerm2(img, 0, 0)
	if err != nil {
		t.Fatalf("RenderITerm2() erro: %v", err)
	}

	const prefix = "\x1b]1337;File=inline=1;width=80;height=20:"
	if !strings.HasPrefix(out, prefix) {
		t.Fatalf("sequência não usa a largura natural limitada: %q", out[:100])
	}
}

func TestRenderITerm2SinglePixel(t *testing.T) {
	img := solidImage(1, 1, solidColor())
	out, err := RenderITerm2(img, 0, 0)
	if err != nil {
		t.Fatalf("RenderITerm2() erro: %v", err)
	}

	prefix := "\x1b]1337;File=inline=1;width=1;height=1:"
	payload := strings.TrimSuffix(strings.TrimPrefix(out, prefix), "\a")
	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		t.Fatalf("payload não é base64 válido: %v", err)
	}
	decoded := decodePNG(t, data)
	if decoded.Bounds().Dx() != 1 || decoded.Bounds().Dy() != 1 {
		t.Fatalf("PNG decodificado tem %dx%d, quer 1x1", decoded.Bounds().Dx(), decoded.Bounds().Dy())
	}
}

func TestRenderITerm2Nil(t *testing.T) {
	if _, err := RenderITerm2(nil, 10, 10); err == nil {
		t.Fatal("RenderITerm2(nil) deveria devolver erro")
	}
}

func TestRenderITerm2EmptyImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))
	if _, err := RenderITerm2(img, 10, 10); err == nil {
		t.Fatal("RenderITerm2(imagem vazia) deveria devolver erro")
	}
}

func TestRenderITerm2Deterministic(t *testing.T) {
	img := gradientImage(20, 10)
	first, err := RenderITerm2(img, 5, 3)
	if err != nil {
		t.Fatalf("RenderITerm2() erro: %v", err)
	}
	second, err := RenderITerm2(img, 5, 3)
	if err != nil {
		t.Fatalf("RenderITerm2() erro: %v", err)
	}
	if first != second {
		t.Fatal("RenderITerm2() não é determinístico")
	}
}
