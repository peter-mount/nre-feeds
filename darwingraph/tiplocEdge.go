package darwingraph

import (
	"encoding/xml"
	"github.com/peter-mount/nre-feeds/util"
	"gonum.org/v1/gonum/graph"
	"strconv"
)

// LinkTiplocs links two tiplocs together
// Returns the new TiplocEdge or nil if one already exists
func (d *RailGraph) LinkTiplocs(a, b string) *TiplocEdge {
	aT := d.GetNode(a)
	bT := d.GetNode(b)

	if aT != nil && bT != nil && aT.NodeType() == NodeTiploc && bT.NodeType() == NodeTiploc {
		aI := aT.ID()
		bI := bT.ID()
		if aI != bI && !d.graph.HasEdgeBetween(aI, bI) {
			edge := &TiplocEdge{f: aT.(*TiplocNode), t: bT.(*TiplocNode)}
			d.graph.SetEdge(edge)
			return edge
		}
	}

	return nil
}

type TiplocEdge struct {
	F   int64       // from tiploc - used in XML parser only
	T   int64       // to  tiploc - used in XML parser only
	Src string      // Source of edge
	f   *TiplocNode // From node
	t   *TiplocNode // To node
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
		AddAttributeIfSet(xml.Name{Local: "src"}, e.Src).
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
		case "src":
			e.Src = attr.Value
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
