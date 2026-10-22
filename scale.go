package charmimage

import (
	"image"
	"image/color"
)

const (
	// defaultMaxWidth é a largura natural máxima, em células, quando o chamador
	// não informa uma largura.
	defaultMaxWidth = 80
	// cellAspect aproxima quantas vezes uma célula de terminal é mais alta do
	// que larga.
	cellAspect = 2
	// sixelCellWidth e sixelCellHeight definem a resolução usada para
	// rasterizar uma imagem Sixel por célula.
	sixelCellWidth  = 8
	sixelCellHeight = 16
)

// resolveSize traduz largura e altura em células para dimensões válidas.
//
// Largura zero ou negativa usa a largura natural da imagem limitada a 80.
// Altura zero ou negativa é derivada da largura preservando a proporção,
// considerando células duas vezes mais altas do que largas.
func resolveSize(srcWidth, srcHeight, width, height int) (int, int) {
	if srcWidth < 1 {
		srcWidth = 1
	}
	if srcHeight < 1 {
		srcHeight = 1
	}

	if width <= 0 {
		width = srcWidth
		if width > defaultMaxWidth {
			width = defaultMaxWidth
		}
	}
	if width < 1 {
		width = 1
	}

	if height <= 0 {
		numerator := int64(width) * int64(srcHeight)
		denominator := int64(srcWidth) * cellAspect
		height = int((numerator + denominator/2) / denominator)
	}
	if height < 1 {
		height = 1
	}

	return width, height
}

// sourceRect devolve o retângulo de pixels da imagem que alimenta a célula
// (col, row) de uma grade de cols×rows. O retângulo nunca é vazio.
func sourceRect(bounds image.Rectangle, col, row, cols, rows int) image.Rectangle {
	if cols < 1 {
		cols = 1
	}
	if rows < 1 {
		rows = 1
	}

	width := int64(bounds.Dx())
	height := int64(bounds.Dy())

	x0 := bounds.Min.X + int(int64(col)*width/int64(cols))
	x1 := bounds.Min.X + int(int64(col+1)*width/int64(cols))
	y0 := bounds.Min.Y + int(int64(row)*height/int64(rows))
	y1 := bounds.Min.Y + int(int64(row+1)*height/int64(rows))

	if x1 <= x0 {
		x1 = x0 + 1
	}
	if y1 <= y0 {
		y1 = y0 + 1
	}
	if x0 < bounds.Min.X {
		x0 = bounds.Min.X
	}
	if y0 < bounds.Min.Y {
		y0 = bounds.Min.Y
	}
	if x1 > bounds.Max.X {
		x1 = bounds.Max.X
	}
	if y1 > bounds.Max.Y {
		y1 = bounds.Max.Y
	}

	return image.Rect(x0, y0, x1, y1)
}

// backgroundRGBA converte a cor de fundo para RGBA opaco, compondo sobre preto
// quando ela própria tem transparência.
func backgroundRGBA(c color.Color) color.RGBA {
	if c == nil {
		return color.RGBA{A: 0xff}
	}

	n := color.NRGBAModel.Convert(c).(color.NRGBA)
	alpha := uint32(n.A)
	return color.RGBA{
		R: uint8((uint32(n.R)*alpha + 127) / 255),
		G: uint8((uint32(n.G)*alpha + 127) / 255),
		B: uint8((uint32(n.B)*alpha + 127) / 255),
		A: 0xff,
	}
}

// averageRect calcula a cor média do retângulo da imagem correspondente à
// célula (col, row), compondo cada pixel sobre bg.
func averageRect(img image.Image, bounds image.Rectangle, col, row, cols, rows int, bg color.RGBA) color.RGBA {
	r := sourceRect(bounds, col, row, cols, rows)

	var sumR, sumG, sumB, count uint64
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			n := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			alpha := uint32(n.A)
			inverse := 255 - alpha

			sumR += uint64((uint32(n.R)*alpha + uint32(bg.R)*inverse + 127) / 255)
			sumG += uint64((uint32(n.G)*alpha + uint32(bg.G)*inverse + 127) / 255)
			sumB += uint64((uint32(n.B)*alpha + uint32(bg.B)*inverse + 127) / 255)
			count++
		}
	}

	if count == 0 {
		return bg
	}

	return color.RGBA{
		R: uint8(sumR / count),
		G: uint8(sumG / count),
		B: uint8(sumB / count),
		A: 0xff,
	}
}

// luminance aproxima a luminância perceptual de uma cor, no intervalo 0..255.
func luminance(c color.RGBA) int {
	return (299*int(c.R) + 587*int(c.G) + 114*int(c.B) + 500) / 1000
}
