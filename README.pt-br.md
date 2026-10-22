# charm-image

Renderize imagens no terminal usando o melhor protocolo que o terminal suporta — Kitty, iTerm2, Sixel ou o fallback universal de half-blocks em truecolor.

Parte do [Charm Suite](https://github.com/junhinhow/charm-suite).

## Instalação

```bash
go get github.com/junhinhow/charm-image
```

Sem dependências externas: apenas a biblioteca padrão do Go.

## Uso

```go
package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	charmimage "github.com/junhinhow/charm-image"
)

func main() {
	file, err := os.Open("foto.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	out, err := charmimage.Render(img, charmimage.WithWidth(60))
	if err != nil {
		panic(err)
	}
	fmt.Print(out)
}
```

`Render` escolhe o protocolo a partir do ambiente. Para forçar um protocolo:

```go
out, err := charmimage.Render(img,
	charmimage.WithProtocol(charmimage.ProtocolHalfBlock),
	charmimage.WithWidth(40),
	charmimage.WithHeight(20),
)
```

Para terminais sem suporte a cor, imprima `AsciiArt`:

```go
fmt.Print(charmimage.AsciiArt(img, 40, 0, ""))
```

## Suporte por terminal

| Terminal | Protocolo | Como é detectado |
| --- | --- | --- |
| Kitty | Kitty | `KITTY_WINDOW_ID` definido |
| WezTerm | iTerm2 | `TERM_PROGRAM=WezTerm` |
| iTerm2 | iTerm2 | `TERM_PROGRAM=iTerm.app` |
| mintty | iTerm2 | `TERM_PROGRAM=mintty` |
| mlterm, xterm com sixel, tmux com recursos de sixel | Sixel | `TERM` ou `COLORTERM` contém `sixel` |
| A maioria dos terminais modernos | Half-block | fallback (truecolor) |
| `TERM=dumb` ou `NO_COLOR` | ASCII | sem sequências de escape; `Render` usa `AsciiArt` |

A detecção é heurística e usa apenas variáveis de ambiente; nenhuma consulta é enviada ao terminal. `WithProtocol` sempre vence a detecção.

## API

| Símbolo | Descrição |
| --- | --- |
| `Protocol`, `ProtocolNone`, `ProtocolHalfBlock`, `ProtocolSixel`, `ProtocolKitty`, `ProtocolITerm2` | Enumeração de protocolos |
| `Detect()`, `DetectFrom(getenv)` | Detecção de protocolo pelo ambiente |
| `Options`, `Option` | Configuração de renderização e opções funcionais |
| `WithWidth(n)`, `WithHeight(n)`, `WithProtocol(p)`, `WithDither(b)` | Construtores de opções |
| `Render(img, opts...)` | Renderiza com o protocolo detectado ou forçado |
| `RenderHalfBlock(img, width, height)` | Fallback universal em truecolor usando `▀` |
| `RenderSixel`, `RenderKitty`, `RenderITerm2` | Sequências de escape cruas de cada protocolo |
| `AsciiArt(img, width, height, ramp)` | Fallback sem cor com rampa configurável |

## Notas de projeto

- **Nome do pacote:** o pacote é `charmimage`, não `image`, para nunca sombrear a biblioteca padrão.
- **Proporção:** assume-se que uma célula de terminal é aproximadamente duas vezes mais alta do que larga. Quando só `Width` é informado, a altura é derivada com essa proporção.
- **Largura natural:** `width <= 0` usa a largura da imagem limitada a 80 células. Valores negativos nunca chegam a slices ou `strings.Repeat`.
- **Half-block:** cada célula desenha duas linhas de pixel com `▀`; o pixel de cima é o foreground e o de baixo é o background, ambos em truecolor (`38;2` / `48;2`). Toda linha termina com `ESC[0m`.
- **Kitty / iTerm2:** a imagem em resolução original é codificada em PNG (stdlib) e enviada em base64. Kitty recebe `c`/`r` em células; iTerm2 recebe `width`/`height` em células. O terminal escala preservando a proporção.
- **Sixel:** rasterizado com 8×16 pixels por célula, quantizado para um cubo fixo de 216 cores mais 40 tons de cinza e codificado com run-length encoding. O payload Sixel é embrulhado em base64 dentro do envelope DCS (`ESC P q … ESC \`), conforme o contrato desta biblioteca.
- **Cursor:** `Render` emite `ESC[1;1H` antes de um bloco Sixel. `RenderSixel` devolve apenas a sequência DCS, para ser composta pelo chamador.
- **Dithering:** `WithDither(true)` aplica um deslocamento ordenado de Bayer 4×4 à luminância no fallback ASCII. Protocolos truecolor não quantizam cores, então não são afetados.
- **Transparência:** nos caminhos half-block, ASCII e Sixel, pixels são compostos sobre `Options.Background`, ou preto quando ele é nil. Kitty e iTerm2 enviam o canal alfa do PNG e deixam o terminal compor. Não existe `WithBackground`; use uma opção personalizada:
  ```go
  charmimage.Option(func(o *charmimage.Options) {
      o.Background = color.White
  })
  ```
- **Determinismo:** a renderização nunca itera sobre maps; a mesma imagem e as mesmas opções sempre produzem a mesma string.

## Limitações

- As sequências de escape são cobertas por testes unitários, mas o projeto não foi exercitado contra terminais Kitty, iTerm2, WezTerm, mintty ou Sixel reais no ambiente de desenvolvimento.
- Os dados Sixel são codificados em base64 dentro do envelope DCS. Terminais Sixel tradicionais esperam dados Sixel sem codificação; decodifique o payload antes de alimentar um terminal diretamente.
- A detecção não consulta o terminal (sem DA1/DA2), então configurações incomuns podem precisar de `WithProtocol`.
- A saída half-block assume suporte a truecolor; terminais limitados a 256 cores vão aproximar as cores.
- O fallback half-block calcula a média dos pixels de origem por célula. Não há reamostragem com correção de gamma nem rotação de imagem.

## Licença

MIT
