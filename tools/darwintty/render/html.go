package render

import (
	"io"
)

func (b *builder) BuildHtml(w io.Writer) {
	b.root().buildHtml(w)
}

func (b *builder) buildHtml(w io.Writer) {
	for _, e := range b.chain {
		switch e.t {
		case root:
		case String:
			_, _ = w.Write(e.content)
		case Newline:
			e.buildHtml(w)
			_, _ = w.Write([]byte{'\n'})
		case Link:
			_, _ = w.Write([]byte("<a href=\""))
			_, _ = w.Write(e.content)
			_, _ = w.Write([]byte{'"', '>'})
			_, _ = w.Write(e.content)
			_, _ = w.Write([]byte("</a>"))
		case Red, Green, Yellow, Blue, Magenta, Cyan, White:
			_, _ = w.Write([]byte("<span class=\"col"))
			_, _ = w.Write([]byte{byte(e.t-Red) + '1', '"', '>'})
			e.buildHtml(w)
			_, _ = w.Write([]byte("</span>"))
		}
	}
}
