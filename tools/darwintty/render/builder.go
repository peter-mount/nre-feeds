package render

import (
	"fmt"
	"io"
	"strings"
)

type Builder interface {
	BuildAnsi() []byte
	BuildHtml(io.Writer)

	End() Builder
	NewLine() Builder

	Repeat(string, int) Builder

	Print(string) Builder
	Printf(string, ...interface{}) Builder
	Println(string) Builder
	Link(string) Builder

	Red() Builder
	Green() Builder
	Yellow() Builder
	Blue() Builder
	Magenta() Builder
	Cyan() Builder
	White() Builder
}

type builder struct {
	p       *builder   // Parent
	t       int        // type
	content []byte     // not root: Content in this builder
	chain   []*builder // root only, order of builders
}

const (
	root = iota
	Newline
	String
	Link
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func New() Builder {
	b := &builder{t: root}
	return b.NewLine()
}

func (b *builder) new(t int) *builder {
	c := &builder{p: b, t: t}
	b.chain = append(b.chain, c)
	return c
}

func (b *builder) root() *builder {
	n := b
	for n.p != nil {
		n = n.p
	}
	return n
}

func (b *builder) End() Builder {
	if b.p != nil {
		return b.p
	}
	return b
}

func (b *builder) NewLine() Builder {
	// Newline ends the current line hence added to the root
	return b.root().new(Newline)
}

func (b *builder) Print(s string) Builder {
	c := b.new(String)
	c.content = append(c.content, s...)
	return b
}

func (b *builder) Printf(s string, a ...interface{}) Builder {
	return b.Print(fmt.Sprintf(s, a...))
}

func (b *builder) Println(s string) Builder {
	return b.Print(s).NewLine()
}

// Link adds a url - do not use End() after this as this is inline
func (b *builder) Link(s string) Builder {
	c := b.new(Link)
	c.content = append(c.content, s...)
	return b
}

func (b *builder) Repeat(s string, c int) Builder {
	return b.Print(strings.Repeat(s, c))
}
func (b *builder) Red() Builder { return b.new(Red) }

func (b *builder) Green() Builder { return b.new(Green) }

func (b *builder) Yellow() Builder { return b.new(Yellow) }

func (b *builder) Blue() Builder { return b.new(Blue) }

func (b *builder) Magenta() Builder { return b.new(Magenta) }

func (b *builder) Cyan() Builder { return b.new(Cyan) }

func (b *builder) White() Builder { return b.new(White) }
