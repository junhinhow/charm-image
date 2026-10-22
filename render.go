package charmimage

import (
	"errors"
	"fmt"
	"image"
)

var (
	// errNilImage indica que Render recebeu uma imagem nil.
	errNilImage = errors.New("charmimage: imagem nil")
	// errEmptyImage indica que a imagem não tem pixels.
	errEmptyImage = errors.New("charmimage: imagem sem pixels")
)

// Render renderiza a imagem com o melhor protocolo disponível e devolve a
// string pronta para ser impressa no terminal.
//
// Sem opções, o protocolo vem de Detect e a largura é a natural limitada a 80
// células. As opções são aplicadas na ordem em que aparecem. Um ProtocolNone em
// Options significa detecção automática; quando Detect devolve ProtocolNone
// (TERM=dumb ou NO_COLOR), a saída é AsciiArt com a rampa padrão.
//
// Devolve erro se img for nil, se a imagem não tiver pixels ou se o protocolo
// for desconhecido.
func Render(img image.Image, opts ...Option) (string, error) {
	if img == nil {
		return "", errNilImage
	}

	bounds := img.Bounds()
	if bounds.Empty() {
		return "", errEmptyImage
	}

	options := Options{Protocol: Detect()}
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}
	options = options.normalize()

	if options.Protocol == ProtocolNone {
		options.Protocol = Detect()
	}

	switch options.Protocol {
	case ProtocolHalfBlock:
		return renderHalfBlock(img, options.Width, options.Height, backgroundRGBA(options.Background)), nil
	case ProtocolSixel:
		out, err := renderSixel(img, options.Width, options.Height, backgroundRGBA(options.Background))
		if err != nil {
			return "", err
		}
		// Posiciona o cursor na origem antes do bloco Sixel. RenderSixel
		// devolve a sequência pura, sem esse posicionamento.
		return "\x1b[1;1H" + out, nil
	case ProtocolKitty:
		return RenderKitty(img, options.Width, options.Height)
	case ProtocolITerm2:
		return RenderITerm2(img, options.Width, options.Height)
	case ProtocolNone:
		return asciiArtWithBackground(img, options.Width, options.Height, defaultRamp, options.Dither, backgroundRGBA(options.Background)), nil
	default:
		return "", fmt.Errorf("charmimage: protocolo desconhecido: %d", int(options.Protocol))
	}
}
