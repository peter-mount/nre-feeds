package util

import (
	"encoding/xml"
	"errors"
	"strconv"
)

type XmlBuilder struct {
	e       *xml.Encoder
	start   xml.StartElement
	err     error
	started bool
	ended   bool
	ifLevel bool
	ifCond  bool
	parent  *XmlBuilder
}

func NewXmlBuilder(e *xml.Encoder, start xml.StartElement) *XmlBuilder {
	return &XmlBuilder{
		e:      e,
		start:  start,
		ifCond: true,
	}
}

func (b *XmlBuilder) If(v bool) *XmlBuilder {
	return &XmlBuilder{
		err:     b.err,
		start:   b.start,
		started: b.started,
		ended:   b.ended,
		ifLevel: true,
		ifCond:  v,
		parent:  b,
	}
}

func (b *XmlBuilder) EndIf() *XmlBuilder {
	if !b.ifLevel {
		b.err = errors.New("endIf called on root builder")
		return b
	} else {
		b.parent.start = b.start
		b.parent.started = b.started
		b.parent.ended = b.ended
		return b.parent
	}
}

func (b *XmlBuilder) startElement() *XmlBuilder {
	if !b.ifCond {
		return b
	}
	if b.err == nil && !b.started {
		b.started = true
		b.err = b.e.EncodeToken(b.start)
	}
	return b
}

func (b *XmlBuilder) Build() error {
	if b.err == nil && !b.started {
		b.startElement()
	}

	if b.err == nil && !b.ended {
		b.ended = true
		b.err = b.e.EncodeToken(xml.EndElement{Name: b.start.Name})
	}
	return b.err
}

func (b *XmlBuilder) Append(n xml.Name, v interface{}) *XmlBuilder {
	if !b.ifCond {
		return b
	}
	if b.err == nil && !b.started {
		b.startElement()
	}
	if b.err == nil {
		b.err = b.e.EncodeElement(v, xml.StartElement{Name: n})
	}
	return b
}

func (b *XmlBuilder) Run(f func(builder *XmlBuilder) error) *XmlBuilder {
	if b.err == nil {
		b.err = f(b)
	}
	return b
}

func (b *XmlBuilder) ElementIf(v bool, n xml.Name, f func(builder *XmlBuilder) error) *XmlBuilder {
	if v {
		return b.Element(n, f)
	}
	return b
}

func (b *XmlBuilder) Element(n xml.Name, f func(builder *XmlBuilder) error) *XmlBuilder {
	if !b.ifCond {
		return b
	}
	if !b.started {
		b.startElement()
	}
	if b.err == nil {
		b.err = NewXmlBuilder(b.e, xml.StartElement{Name: n}).
			Run(f).
			Build()
	}
	return b
}

func (b *XmlBuilder) AddAttribute(n xml.Name, v string) *XmlBuilder {
	if !b.ifCond {
		return b
	}
	if b.err == nil {
		if b.started {
			b.err = errors.New("add attr called after start")
		} else {
			b.start.Attr = append(b.start.Attr, xml.Attr{Name: n, Value: v})
		}
	}
	return b
}

func (b *XmlBuilder) AddAttributeIfSet(n xml.Name, v string) *XmlBuilder {
	if v == "" {
		return b
	}
	return b.AddAttribute(n, v)
}

func (b *XmlBuilder) AddBoolAttribute(n xml.Name, v bool) *XmlBuilder {
	if v {
		return b.AddAttribute(n, "true")
	}
	return b.AddAttribute(n, "false")
}

func (b *XmlBuilder) AddBoolAttributeIfSet(n xml.Name, v bool) *XmlBuilder {
	if v {
		return b.AddBoolAttribute(n, v)
	}
	return b
}

func (b *XmlBuilder) AddFloat32Attribute(n xml.Name, v float32) *XmlBuilder {
	return b.AddAttribute(n, strconv.FormatFloat(float64(v), 'f', 10, 32))
}

func (b *XmlBuilder) AddFloat32AttributeIfSet(n xml.Name, v float32) *XmlBuilder {
	if v != 0 {
		return b.AddFloat32Attribute(n, v)
	}
	return b
}

func (b *XmlBuilder) AddFloatAttribute(n xml.Name, v float64) *XmlBuilder {
	return b.AddAttribute(n, strconv.FormatFloat(v, 'f', 10, 64))
}

func (b *XmlBuilder) AddFloatAttributeIfSet(n xml.Name, v float64) *XmlBuilder {
	if v != 0 {
		return b.AddFloatAttribute(n, v)
	}
	return b
}
