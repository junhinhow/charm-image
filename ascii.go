package charmimage

import (
	"image"
	"image/color"
	"strings"
)

// defaultRamp é a rampa de caracteres usada quando nenhuma é informada.
const defaultRamp = " .:-=+*#%@"

// bayer4 é a matriz de Bayer 4×4 usada no dithering ordenado.
var bayer4 = [16]int{
	0, 8, 2, 10,
	12, 4, 14, 6,
	3, 11, 1, 9,
	15, 7, 13, 5,
}

// AsciiArt desenha a imagem usando apenas caracteres, sem cores. A rampa vai do
// tom mais escuro (primeiro caractere) ao mais claro (último). Rampa vazia usa a
// rampa padrão " .:-=+*#%@".
//
// Largura ou altura zero/negativa são resolvidas como em Render. Devolve string
// vazia se img for nil ou não tiver pixels.
func AsciiArt(img image.Image, width, height int, ramp string) string {
	return asciiArt(img, width, height, ramp, false)
}

// asciiArt é a implementação interna com fundo preto, usada por AsciiArt e
// pelos testes.
func asciiArt(img image.Image, width, height int, ramp string, dither bool) string {
	return asciiArtWithBackground(img, width, height, ramp, dither, backgroundRGBA(nil))
}

// asciiArtWithBackground permite compor pixels transparentes sobre bg.
func asciiArtWithBackground(img image.Image, width, height int, ramp string, dither bool, bg color.RGBA) string {
	if img == nil {
		return ""
	}

	bounds := img.Bounds()
	if bounds.Empty() {
		return ""
	}

	runes := []rune(ramp)
	if len(runes) == 0 {
		runes = []rune(defaultRamp)
	}

	cols, rows := resolveSize(bounds.Dx(), bounds.Dy(), width, height)
	last := len(runes) - 1

	var sb strings.Builder
	sb.Grow(cols*rows + rows)

	for row := 0; row < rows; row++ {
		if row > 0 {
			sb.WriteByte('\n')
		}
		for col := 0; col < cols; col++ {
			c := averageRect(img, bounds, col, row, cols, rows, bg)
			lum := luminance(c)

			if dither && last > 0 {
				step := 255 / last
				lum += (bayer4[(row&3)*4+(col&3)] - 8) * step / 8
				if lum < 0 {
					lum = 0
				} else if lum > 255 {
					lum = 255
				}
			}

			index := 0
			if last > 0 {
				index = lum * last / 255
			}
			if index < 0 {
				index = 0
			} else if index > last {
				index = last
			}

			sb.WriteRune(runes[index])
		}
	}

	return sb.String()
}
