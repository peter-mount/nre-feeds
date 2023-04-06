package render

import (
	"bytes"
	"fmt"
	"io"
)

func (b *builder) BuildAnsi() []byte {
	var out bytes.Buffer
	b.root().buildAnsi(&out)
	return out.Bytes()
}

func (b *builder) buildAnsi(w io.Writer) {
	for _, e := range b.chain {
		switch e.t {
		case root:
		case String:
			_, _ = w.Write(e.content)
		case Newline:
			e.buildAnsi(w)
			_, _ = w.Write([]byte{033, '[', 'm', '\n'})
		case Link:
			_, _ = w.Write(e.content)
		case Red, Green, Yellow, Blue, Magenta, Cyan, White:
			// 30+col for foreground, 40+col for background
			_, _ = w.Write([]byte(fmt.Sprintf("\033[%dm", e.t-Red+31)))
			e.buildAnsi(w)
			_, _ = w.Write([]byte{033, '[', 'm'})
		}
	}
}
