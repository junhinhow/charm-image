package charmimage

import "image/color"

// Options controlam a renderização de uma imagem.
type Options struct {
	// Width é a largura desejada em células. Zero ou negativo usa a largura
	// natural, limitada a 80 células.
	Width int
	// Height é a altura desejada em células. Zero ou negativo calcula a altura
	// proporcional à largura, assumindo células aproximadamente duas vezes mais
	// altas do que largas.
	Height int
	// Protocol é o protocolo desejado. O valor zero (ProtocolNone) significa
	// "detectar automaticamente" com Detect.
	Protocol Protocol
	// Dither liga o dithering ordenado no fallback ASCII.
	Dither bool
	// Background é a cor usada para compor pixels transparentes. Se for nil,
	// usa preto.
	Background color.Color
}

// Option modifica um Options durante Render.
type Option func(*Options)

// WithWidth define a largura em células. Valores negativos são tratados como
// zero (largura natural limitada a 80).
func WithWidth(n int) Option {
	return func(o *Options) {
		o.Width = n
	}
}

// WithHeight define a altura em células. Valores negativos são tratados como
// zero (altura proporcional).
func WithHeight(n int) Option {
	return func(o *Options) {
		o.Height = n
	}
}

// WithProtocol força um protocolo, ignorando a detecção automática. Use
// ProtocolHalfBlock para o fallback com cores. O valor ProtocolNone volta a
// significar "detectar automaticamente" em Options.
func WithProtocol(p Protocol) Option {
	return func(o *Options) {
		o.Protocol = p
	}
}

// WithDither liga ou desliga o dithering ordenado (Bayer 4×4) aplicado à
// luminância no fallback ASCII. Não tem efeito nos protocolos coloridos.
func WithDither(b bool) Option {
	return func(o *Options) {
		o.Dither = b
	}
}

// normalize aplica as invariantes numéricas de Options.
func (o Options) normalize() Options {
	if o.Width < 0 {
		o.Width = 0
	}
	if o.Height < 0 {
		o.Height = 0
	}
	return o
}
