// Package charmimage renderiza imagens em terminais, escolhendo entre os
// protocolos Kitty, iTerm2, Sixel e o fallback universal de half-blocks.
//
// O pacote não tem dependências externas: usa apenas a biblioteca padrão do Go.
// A saída de cada protocolo é uma string pronta para ser impressa no terminal.
//
// Exemplo:
//
//	out, err := charmimage.Render(img, charmimage.WithWidth(40))
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Print(out)
package charmimage
