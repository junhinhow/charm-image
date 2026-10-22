package charmimage

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"strings"
	"testing"
)

// decodeKittyPNG extrai e decodifica o PNG de uma sequência Kitty sem chunks.
func decodeKittyPNG(t *testing.T, out string) image.Image {
	t.Helper()

	const suffix = "\x1b\\"
	if !strings.HasPrefix(out, "\x1b_Gf=100,a=T,") {
		t.Fatalf("sequência não começa com APC do Kitty: %q", out)
	}
	if !strings.HasSuffix(out, suffix) {
		t.Fatalf("sequência não termina com ST: %q", out)
	}

	body := strings.TrimSuffix(out, suffix)
	semicolon := strings.LastIndexByte(body, ';')
	if semicolon < 0 {
		t.Fatalf("sequência sem payload: %q", out)
	}

	payload := body[semicolon+1:]
	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		t.Fatalf("payload não é base64 válido: %v", err)
	}
	return decodePNG(t, data)
}

// decodePNG decodifica bytes PNG em uma imagem.
func decodePNG(t *testing.T, data []byte) image.Image {
	t.Helper()

	decoded, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("payload não é PNG válido: %v", err)
	}
	return decoded
}

func TestRenderKittySequence(t *testing.T) {
	img := solidImage(4, 3, solidColor())
	out, err := RenderKitty(img, 2, 3)
	if err != nil {
		t.Fatalf("RenderKitty() erro: %v", err)
	}

	const prefix = "\x1b_Gf=100,a=T,c=2,r=3,q=2;"
	if !strings.HasPrefix(out, prefix) {
		t.Fatalf("sequência não começa com %q: %q", prefix, out)
	}

	decoded := decodeKittyPNG(t, out)
	bounds := decoded.Bounds()
	if bounds.Dx() != 4 || bounds.Dy() != 3 {
		t.Fatalf("PNG decodificado tem %dx%d, quer 4x3", bounds.Dx(), bounds.Dy())
	}
}

func TestRenderKittyNaturalSize(t *testing.T) {
	img := solidImage(100, 50, solidColor())
	out, err := RenderKitty(img, 0, 0)
	if err != nil {
		t.Fatalf("RenderKitty() erro: %v", err)
	}

	const prefix = "\x1b_Gf=100,a=T,c=80,r=20,q=2;"
	if !strings.HasPrefix(out, prefix) {
		t.Fatalf("sequência não usa a largura natural limitada: %q", out[:80])
	}
}

func TestRenderKittyChunksLargePayload(t *testing.T) {
	img := noisyImage(96, 96)
	out, err := RenderKitty(img, 48, 48)
	if err != nil {
		t.Fatalf("RenderKitty() erro: %v", err)
	}

	if count := strings.Count(out, "\x1b_G"); count < 2 {
		t.Fatalf("payload grande não foi dividido em chunks: %d sequências", count)
	}
	if !strings.Contains(out, "m=1;") {
		t.Fatal("faltou o marcador m=1 nos chunks intermediários")
	}

	var payload strings.Builder
	for _, part := range strings.Split(out, "\x1b_G") {
		if part == "" {
			continue
		}
		part = strings.TrimSuffix(part, "\x1b\\")
		semicolon := strings.IndexByte(part, ';')
		if semicolon < 0 {
			t.Fatalf("chunk sem payload: %q", part)
		}
		payload.WriteString(part[semicolon+1:])
	}

	data, err := base64.StdEncoding.DecodeString(payload.String())
	if err != nil {
		t.Fatalf("chunks não formam base64 válido: %v", err)
	}
	decoded := decodePNG(t, data)
	if decoded.Bounds().Dx() != 96 || decoded.Bounds().Dy() != 96 {
		t.Fatalf("PNG reagrupado tem %dx%d, quer 96x96", decoded.Bounds().Dx(), decoded.Bounds().Dy())
	}
}

func TestRenderKittyNil(t *testing.T) {
	if _, err := RenderKitty(nil, 10, 10); err == nil {
		t.Fatal("RenderKitty(nil) deveria devolver erro")
	}
}

func TestRenderKittyEmptyImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))
	if _, err := RenderKitty(img, 10, 10); err == nil {
		t.Fatal("RenderKitty(imagem vazia) deveria devolver erro")
	}
}

func TestRenderKittyDeterministic(t *testing.T) {
	img := gradientImage(20, 10)
	first, err := RenderKitty(img, 5, 3)
	if err != nil {
		t.Fatalf("RenderKitty() erro: %v", err)
	}
	second, err := RenderKitty(img, 5, 3)
	if err != nil {
		t.Fatalf("RenderKitty() erro: %v", err)
	}
	if first != second {
		t.Fatal("RenderKitty() não é determinístico")
	}
}
