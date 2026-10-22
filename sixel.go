package charmimage

import (
	"encoding/base64"
	"image"
	"image/color"
	"strconv"
	"strings"
)

// RenderSixel devolve a imagem codificada no protocolo Sixel, embrulhada em um
// DCS (ESC P q ... ESC \) com o payload em base64.
//
// A imagem é rasterizada com 8×16 pixels por célula (a mesma proporção de
// célula usada no restante do pacote). Largura ou altura zero/negativa são
// resolvidas como em Render. Devolve erro se img for nil ou não tiver pixels.
func RenderSixel(img image.Image, width, height int) (string, error) {
	return renderSixel(img, width, height, backgroundRGBA(nil))
}

// renderSixel é a implementação usada por Render e permite configurar a cor de
// fundo usada na composição de transparências.
func renderSixel(img image.Image, width, height int, bg color.RGBA) (string, error) {
	if img == nil {
		return "", errNilImage
	}

	bounds := img.Bounds()
	if bounds.Empty() {
		return "", errEmptyImage
	}

	cols, rows := resolveSize(bounds.Dx(), bounds.Dy(), width, height)
	pixelWidth := cols * sixelCellWidth
	pixelHeight := rows * sixelCellHeight

	pixels := make([]uint8, pixelWidth*pixelHeight)
	var used [256]bool

	for y := 0; y < pixelHeight; y++ {
		for x := 0; x < pixelWidth; x++ {
			c := averageRect(img, bounds, x, y, pixelWidth, pixelHeight, bg)
			index := quantizeSixel(c)
			pixels[y*pixelWidth+x] = index
			used[index] = true
		}
	}

	payload := buildSixelPayload(pixelWidth, pixelHeight, pixels, used)
	return "\x1bPq" + base64.StdEncoding.EncodeToString([]byte(payload)) + "\x1b\\", nil
}

// quantizeSixel mapeia uma cor para um índice da paleta fixa de 256 cores:
// 216 cores do cubo 6×6×6 seguidas de 40 tons de cinza.
func quantizeSixel(c color.RGBA) uint8 {
	r := int(c.R)
	g := int(c.G)
	b := int(c.B)

	// Cubo 6×6×6.
	ri := (r*5 + 127) / 255
	gi := (g*5 + 127) / 255
	bi := (b*5 + 127) / 255
	cubeIndex := ri*36 + gi*6 + bi
	cubeR, cubeG, cubeB := ri*51, gi*51, bi*51

	// Tons de cinza.
	grayLevel := (luminance(c)*39 + 127) / 255
	gray := grayLevel * 255 / 39

	cubeDistance := square(r-cubeR) + square(g-cubeG) + square(b-cubeB)
	grayDistance := square(r-gray) + square(g-gray) + square(b-gray)

	if grayDistance < cubeDistance {
		return uint8(216 + grayLevel)
	}
	return uint8(cubeIndex)
}

// square devolve n*n.
func square(n int) int {
	return n * n
}

// sixelColor devolve a cor RGB de um índice da paleta fixa.
func sixelColor(index uint8) color.RGBA {
	if index < 216 {
		r := int(index) / 36
		g := (int(index) % 36) / 6
		b := int(index) % 6
		return color.RGBA{R: uint8(r * 51), G: uint8(g * 51), B: uint8(b * 51), A: 0xff}
	}

	gray := (int(index) - 216) * 255 / 39
	return color.RGBA{R: uint8(gray), G: uint8(gray), B: uint8(gray), A: 0xff}
}

// buildSixelPayload monta os dados Sixel (sem o DCS nem o ST): atributos de
// raster, definições de paleta e bandas de 6 linhas com RLE.
func buildSixelPayload(width, height int, pixels []uint8, used [256]bool) string {
	var sb strings.Builder
	sb.Grow(width*height/4 + 256)

	sb.WriteString("\"1;1;")
	sb.WriteString(strconv.Itoa(width))
	sb.WriteByte(';')
	sb.WriteString(strconv.Itoa(height))

	for index := 0; index < 256; index++ {
		if !used[index] {
			continue
		}
		c := sixelColor(uint8(index))
		sb.WriteByte('#')
		sb.WriteString(strconv.Itoa(index))
		sb.WriteString(";2;")
		sb.WriteString(strconv.Itoa(int(c.R) * 100 / 255))
		sb.WriteByte(';')
		sb.WriteString(strconv.Itoa(int(c.G) * 100 / 255))
		sb.WriteByte(';')
		sb.WriteString(strconv.Itoa(int(c.B) * 100 / 255))
	}

	bands := (height + 5) / 6
	for band := 0; band < bands; band++ {
		if band > 0 {
			sb.WriteByte('-')
		}

		var present [256]bool
		start := band * 6
		end := start + 6
		if end > height {
			end = height
		}
		for y := start; y < end; y++ {
			row := pixels[y*width : (y+1)*width]
			for _, index := range row {
				present[index] = true
			}
		}

		first := true
		for index := 0; index < 256; index++ {
			if !present[index] {
				continue
			}
			if !first {
				sb.WriteByte('$')
			}
			first = false

			sb.WriteByte('#')
			sb.WriteString(strconv.Itoa(index))
			writeSixelBand(&sb, width, height, band, pixels, uint8(index))
		}
	}

	return sb.String()
}

// writeSixelBand escreve as colunas de uma banda para uma cor, aplicando
// run-length encoding.
func writeSixelBand(sb *strings.Builder, width, height, band int, pixels []uint8, index uint8) {
	runChar := byte(0)
	runLength := 0

	flush := func() {
		if runLength == 0 {
			return
		}
		if runLength > 3 {
			sb.WriteByte('!')
			sb.WriteString(strconv.Itoa(runLength))
			sb.WriteByte(runChar)
		} else {
			for i := 0; i < runLength; i++ {
				sb.WriteByte(runChar)
			}
		}
		runLength = 0
	}

	for x := 0; x < width; x++ {
		bits := 0
		for bit := 0; bit < 6; bit++ {
			y := band*6 + bit
			if y < height && pixels[y*width+x] == index {
				bits |= 1 << bit
			}
		}

		ch := byte('?') + byte(bits)
		if ch == runChar {
			runLength++
			continue
		}

		flush()
		runChar = ch
		runLength = 1
	}

	flush()
}
