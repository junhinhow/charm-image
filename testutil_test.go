package charmimage

import (
	"image"
	"image/color"
	"image/draw"
)

// solidImage cria uma imagem preenchida com uma cor.
func solidImage(width, height int, c color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: c}, image.Point{}, draw.Src)
	return img
}

// gradientImage cria um degradê horizontal determinístico.
func gradientImage(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			value := uint8(0)
			if width > 1 {
				value = uint8(x * 255 / (width - 1))
			}
			img.SetRGBA(x, y, color.RGBA{R: value, G: value, B: value, A: 0xff})
		}
	}
	return img
}

// withBackground é um Option de teste para a cor de fundo.
func withBackground(c color.Color) Option {
	return func(o *Options) {
		o.Background = c
	}
}

// solidColor é uma cor opaca usada nos testes.
func solidColor() color.RGBA {
	return color.RGBA{R: 0x33, G: 0x66, B: 0x99, A: 0xff}
}

// noisyImage cria uma imagem de alta frequência, difícil de comprimir, para
// exercitar o particionamento do protocolo Kitty.
func noisyImage(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetRGBA(x, y, color.RGBA{
				R: uint8((x*73 + y*41) % 256),
				G: uint8((x*29 + y*151) % 256),
				B: uint8((x*197 + y*17) % 256),
				A: 0xff,
			})
		}
	}
	return img
}
