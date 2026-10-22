package charmimage

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
)

// encodePNG codifica a imagem em PNG. O formato é exigido pelos protocolos
// Kitty (f=100) e iTerm2, que recebem arquivos de imagem em base64.
func encodePNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	encoder := png.Encoder{CompressionLevel: png.BestSpeed}
	if err := encoder.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("charmimage: falha ao codificar imagem em PNG: %w", err)
	}
	return buf.Bytes(), nil
}
