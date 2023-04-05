package telstar

import (
	"fmt"
	"strings"
)

type FrameBuilder struct {
	response           *Response
	header             string
	content            []string
	contentType        string
	navMessage         string
	navMessageNotFound string
	routingTable       []int
}

func (r *Response) NewFrame() *FrameBuilder {
	return &FrameBuilder{
		response:     r,
		contentType:  "markup",
		routingTable: make([]int, 11),
	}
}

func (b *FrameBuilder) Route(i, p int) *FrameBuilder {
	if i >= 0 && i <= 10 {
		b.routingTable[i] = p
	}
	return b
}

func (b *FrameBuilder) Header(s string, a ...interface{}) *FrameBuilder {
	b.header = fmt.Sprintf(s, a...)
	return b
}

func (b *FrameBuilder) NavMessage(s string, a ...interface{}) *FrameBuilder {
	b.navMessage = fmt.Sprintf(s, a...)
	return b
}

func (b *FrameBuilder) NavMessageNotFound(s string, a ...interface{}) *FrameBuilder {
	b.navMessageNotFound = fmt.Sprintf(s, a...)
	return b
}

func (b *FrameBuilder) ContentType(s string) *FrameBuilder {
	b.contentType = s
	return b
}

func (b *FrameBuilder) Print(s string) *FrameBuilder {
	b.content = append(b.content, s)
	return b
}

func (b *FrameBuilder) Println(s string) *FrameBuilder {
	return b.Print(s).Print("\r\n")
}

func (b *FrameBuilder) Printf(f string, a ...interface{}) *FrameBuilder {
	return b.Print(fmt.Sprintf(f, a...))
}

func (b *FrameBuilder) Build() *Response {
	if b.contentType == "" {
		b.contentType = "markup"
	}

	f := &Frame{

		HeaderText: b.header,
		FrameType:  "information",
		Content: Content{
			Data: strings.Join(b.content, ""),
			Type: b.contentType,
		},
		RoutingTable:       b.routingTable,
		Cursor:             false,
		NavMessage:         b.navMessage,
		NavMessageNotFound: b.navMessageNotFound,
	}

	return b.response.addFrame(f)
}

func (b *FrameBuilder) Red() *FrameBuilder     { return b.Print("[R]") }
func (b *FrameBuilder) Green() *FrameBuilder   { return b.Print("[G]") }
func (b *FrameBuilder) Yellow() *FrameBuilder  { return b.Print("[Y]") }
func (b *FrameBuilder) Blue() *FrameBuilder    { return b.Print("[B]") }
func (b *FrameBuilder) Magenta() *FrameBuilder { return b.Print("[M]") }
func (b *FrameBuilder) Cyan() *FrameBuilder    { return b.Print("[C]") }
func (b *FrameBuilder) White() *FrameBuilder   { return b.Print("[W]") }

func (b *FrameBuilder) Flash() *FrameBuilder  { return b.Print("[F]") }
func (b *FrameBuilder) Steady() *FrameBuilder { return b.Print("[S]") }

func (b *FrameBuilder) NormalHeight() *FrameBuilder { return b.Print("[N]") }
func (b *FrameBuilder) DoubleHeight() *FrameBuilder { return b.Print("[D]") }

func (b *FrameBuilder) NewBackground() *FrameBuilder { return b.Print("[n]") }

func (b *FrameBuilder) MosaicRed() *FrameBuilder     { return b.Print("[r]") }
func (b *FrameBuilder) MosaicGreen() *FrameBuilder   { return b.Print("[g]") }
func (b *FrameBuilder) MosaicYellow() *FrameBuilder  { return b.Print("[y]") }
func (b *FrameBuilder) MosaicBlue() *FrameBuilder    { return b.Print("[b]") }
func (b *FrameBuilder) MosaicMagenta() *FrameBuilder { return b.Print("[m]") }
func (b *FrameBuilder) MosaicCyan() *FrameBuilder    { return b.Print("[c]") }
func (b *FrameBuilder) MosaicWhite() *FrameBuilder   { return b.Print("[w]") }

func (b *FrameBuilder) SepGraphDotsHigh() *FrameBuilder { return b.Print("[h.]") }
func (b *FrameBuilder) SepGraphDotsMid() *FrameBuilder  { return b.Print("[m.]") }
func (b *FrameBuilder) SepGraphDotsLow() *FrameBuilder  { return b.Print("[l.]") }

func (b *FrameBuilder) SepGraphSolidHigh() *FrameBuilder { return b.Print("[h-]") }
func (b *FrameBuilder) SepGraphSolidMid() *FrameBuilder  { return b.Print("[m-]") }
func (b *FrameBuilder) SepGraphSolidLow() *FrameBuilder  { return b.Print("[l-]") }

func (b *FrameBuilder) NewLine() *FrameBuilder { return b.Print("\r\n") }
