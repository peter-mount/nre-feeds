package darwingraph

import (
	"encoding/xml"
	"github.com/peter-mount/nre-feeds/util"
	"gonum.org/v1/gonum/graph"
	"strconv"
)

type TiplocEdge struct {
	F int64
	T int64
	f *TiplocNode
	t *TiplocNode
}

// From returns the from-node of the edge.
func (e TiplocEdge) From() graph.Node { return e.f }

// To returns the to-node of the edge.
func (e TiplocEdge) To() graph.Node { return e.t }

// ReversedEdge returns a new Edge with the F and T fields swapped.
func (e TiplocEdge) ReversedEdge() graph.Edge { return TiplocEdge{f: e.t, t: e.f} }

func (e *TiplocEdge) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return util.NewXmlBuilder(encoder, start).
		AddAttribute(xml.Name{Local: "from"}, strconv.FormatInt(e.f.ID(), IdBase)).
		AddAttribute(xml.Name{Local: "to"}, strconv.FormatInt(e.t.ID(), IdBase)).
		Build()
}

func (e *TiplocEdge) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		var err error
		switch attr.Name.Local {
		case "from":
			e.F, err = strconv.ParseInt(attr.Value, IdBase, 64)
		case "to":
			e.T, err = strconv.ParseInt(attr.Value, IdBase, 64)
		}
		if err != nil {
			return err
		}
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch token.(type) {
		case xml.EndElement:
			return nil
		}
	}
}
