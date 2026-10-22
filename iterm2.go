package charmimage

import (
	"encoding/base64"
	"fmt"
	"image"
)

// RenderITerm2 devolve a imagem no protocolo inline do iTerm2:
// OSC 1337;File=inline=1;width=...;height=...:<base64> BEL.
//
// width e height são células (zero/negativo segue a regra de Render); o
// terminal ajusta a imagem à caixa preservando a proporção. Devolve erro se img
// for nil ou não tiver pixels.
func RenderITerm2(img image.Image, width, height int) (string, error) {
	if img == nil {
		return "", errNilImage
	}

	bounds := img.Bounds()
	if bounds.Empty() {
		return "", errEmptyImage
	}

	cols, rows := resolveSize(bounds.Dx(), bounds.Dy(), width, height)
	data, err := encodePNG(img)
	if err != nil {
		return "", err
	}

	payload := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf(
		"\x1b]1337;File=inline=1;width=%d;height=%d:%s\a",
		cols, rows, payload,
	), nil
}
