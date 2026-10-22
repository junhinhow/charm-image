package charmimage

import (
	"image"
	"image/color"
	"strconv"
	"strings"
)

// halfBlockRune é o caractere usado para desenhar dois pixels por célula.
const halfBlockRune = "▀"

// RenderHalfBlock desenha a imagem com o caractere ▀, usando a cor do pixel de
// cima como foreground e a do pixel de baixo como background. Cada linha de
// célula representa duas linhas de pixel.
//
// Largura ou altura zero/negativa são resolvidas como em Render: largura
// natural limitada a 80 e altura proporcional. Pixels transparentes são
// compostos sobre preto. Devolve string vazia se img for nil ou não tiver
// pixels.
func RenderHalfBlock(img image.Image, width, height int) string {
	return renderHalfBlock(img, width, height, backgroundRGBA(nil))
}

// renderHalfBlock é a implementação usada por Render e permite configurar a cor
// de fundo.
func renderHalfBlock(img image.Image, width, height int, bg color.RGBA) string {
	if img == nil {
		return ""
	}

	bounds := img.Bounds()
	if bounds.Empty() {
		return ""
	}

	cols, cellRows := resolveSize(bounds.Dx(), bounds.Dy(), width, height)
	pixelRows := cellRows * 2

	var sb strings.Builder
	sb.Grow(cols * cellRows * 40)

	for row := 0; row < cellRows; row++ {
		if row > 0 {
			sb.WriteByte('\n')
		}
		for col := 0; col < cols; col++ {
			top := averageRect(img, bounds, col, row*2, cols, pixelRows, bg)
			bottom := averageRect(img, bounds, col, row*2+1, cols, pixelRows, bg)

			writeTruecolor(&sb, 38, top)
			writeTruecolor(&sb, 48, bottom)
			sb.WriteString(halfBlockRune)
		}
		sb.WriteString("\x1b[0m")
	}

	return sb.String()
}

// writeTruecolor escreve uma sequência SGR truecolor para a camada 38
// (foreground) ou 48 (background).
func writeTruecolor(sb *strings.Builder, layer int, c color.RGBA) {
	sb.WriteString("\x1b[")
	sb.WriteString(strconv.Itoa(layer))
	sb.WriteString(";2;")
	sb.WriteString(strconv.Itoa(int(c.R)))
	sb.WriteByte(';')
	sb.WriteString(strconv.Itoa(int(c.G)))
	sb.WriteByte(';')
	sb.WriteString(strconv.Itoa(int(c.B)))
	sb.WriteByte('m')
}
