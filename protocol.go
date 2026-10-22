package charmimage

import (
	"os"
	"strings"
)

// Protocol identifica um protocolo de exibição de imagens suportado pelo
// terminal.
type Protocol int

const (
	// ProtocolNone indica que o terminal não deve receber sequências de imagem.
	// Nesse caso o chamador deve usar AsciiArt (o próprio Render faz isso).
	ProtocolNone Protocol = iota
	// ProtocolHalfBlock desenha a imagem com o caractere ▀ e cores truecolor.
	ProtocolHalfBlock
	// ProtocolSixel usa o protocolo Sixel (DCS).
	ProtocolSixel
	// ProtocolKitty usa o protocolo de imagens do Kitty (APC).
	ProtocolKitty
	// ProtocolITerm2 usa o protocolo de imagens inline do iTerm2 (OSC 1337).
	ProtocolITerm2
)

// String devolve o nome canônico do protocolo.
func (p Protocol) String() string {
	switch p {
	case ProtocolNone:
		return "none"
	case ProtocolHalfBlock:
		return "halfblock"
	case ProtocolSixel:
		return "sixel"
	case ProtocolKitty:
		return "kitty"
	case ProtocolITerm2:
		return "iterm2"
	default:
		return "unknown"
	}
}

// Detect descobre o melhor protocolo de imagem a partir das variáveis de
// ambiente do processo.
func Detect() Protocol {
	return DetectFrom(os.Getenv)
}

// DetectFrom descobre o melhor protocolo usando a função getenv fornecida.
//
// A precedência é: TERM=dumb ou NO_COLOR definem ProtocolNone; KITTY_WINDOW_ID
// define ProtocolKitty; TERM_PROGRAM em {iTerm.app, WezTerm, mintty} define
// ProtocolITerm2; TERM ou COLORTERM contendo "sixel" define ProtocolSixel;
// caso contrário, ProtocolHalfBlock. Se getenv for nil, devolve
// ProtocolHalfBlock.
func DetectFrom(getenv func(string) string) Protocol {
	if getenv == nil {
		return ProtocolHalfBlock
	}

	term := strings.TrimSpace(getenv("TERM"))
	if strings.EqualFold(term, "dumb") || strings.TrimSpace(getenv("NO_COLOR")) != "" {
		return ProtocolNone
	}

	if strings.TrimSpace(getenv("KITTY_WINDOW_ID")) != "" {
		return ProtocolKitty
	}

	switch strings.ToLower(strings.TrimSpace(getenv("TERM_PROGRAM"))) {
	case "iterm.app", "wezterm", "mintty":
		return ProtocolITerm2
	}

	colorTerm := strings.TrimSpace(getenv("COLORTERM"))
	if strings.Contains(strings.ToLower(term), "sixel") ||
		strings.Contains(strings.ToLower(colorTerm), "sixel") {
		return ProtocolSixel
	}

	return ProtocolHalfBlock
}
