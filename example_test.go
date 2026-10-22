package charmimage

import (
	"fmt"
	"image"
	"image/color"
)

func ExampleAsciiArt() {
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.SetRGBA(0, 0, color.RGBA{A: 0xff})
	img.SetRGBA(1, 0, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})

	fmt.Println(AsciiArt(img, 2, 1, "#."))
	// Output:
	// #.
}

func ExampleDetectFrom() {
	env := map[string]string{"KITTY_WINDOW_ID": "1"}
	fmt.Println(DetectFrom(func(key string) string { return env[key] }))
	// Output:
	// kitty
}

func ExampleRender() {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.SetRGBA(0, 0, color.RGBA{R: 0xff, A: 0xff})
	img.SetRGBA(1, 0, color.RGBA{R: 0xff, A: 0xff})
	img.SetRGBA(0, 1, color.RGBA{B: 0xff, A: 0xff})
	img.SetRGBA(1, 1, color.RGBA{B: 0xff, A: 0xff})

	out, err := Render(img, WithProtocol(ProtocolHalfBlock), WithWidth(1), WithHeight(1))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q\n", out)
	// Output:
	// "\x1b[38;2;255;0;0m\x1b[48;2;0;0;255m▀\x1b[0m"
}
