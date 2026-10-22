package charmimage

import (
	"encoding/base64"
	"fmt"
	"image"
	"strings"
)

// kittyChunkSize é o tamanho máximo, em bytes de base64, de cada sequência APC
// do protocolo Kitty.
const kittyChunkSize = 4096

// RenderKitty devolve a imagem no protocolo de imagens do Kitty: uma sequência
// APC ESC _ G ... ESC \ com o PNG em base64.
//
// width e height são células (zero/negativo segue a regra de Render) e viram os
// parâmetros c e r, que o terminal usa para dimensionar a imagem. Devolve erro
// se img for nil ou não tiver pixels.
func RenderKitty(img image.Image, width, height int) (string, error) {
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

	var sb strings.Builder
	if len(payload) <= kittyChunkSize {
		fmt.Fprintf(&sb, "\x1b_Gf=100,a=T,c=%d,r=%d,q=2;%s\x1b\\", cols, rows, payload)
		return sb.String(), nil
	}

	first := true
	for len(payload) > 0 {
		n := kittyChunkSize
		if len(payload) < n {
			n = len(payload)
		}

		chunk := payload[:n]
		payload = payload[n:]

		more := 0
		if len(payload) > 0 {
			more = 1
		}

		if first {
			fmt.Fprintf(&sb, "\x1b_Gf=100,a=T,c=%d,r=%d,q=2,m=%d;%s\x1b\\", cols, rows, more, chunk)
			first = false
			continue
		}

		fmt.Fprintf(&sb, "\x1b_Gm=%d;%s\x1b\\", more, chunk)
	}

	return sb.String(), nil
}
