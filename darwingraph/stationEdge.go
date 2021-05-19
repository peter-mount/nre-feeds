package darwingraph

import (
	"encoding/xml"
	"github.com/peter-mount/nre-feeds/util"
	"gonum.org/v1/gonum/graph"
	"strconv"
	"strings"
)

type StationEdge struct {
	F   int64         // from tiploc - used in XML parser only
	T   int64         // to  tiploc - used in XML parser only
	Src string        // Source of edge
	f   RailNode      // From node
	t   RailNode      // To node
	s   []*TiplocNode // Tiplocs forming this edge
	ss  []int64       // used in unmarshalling
}

// From returns the from-node of the edge.
func (e StationEdge) From() graph.Node { return e.f }

// To returns the to-node of the edge.
func (e StationEdge) To() graph.Node { return e.t }

// ReversedEdge returns a new Edge with the F and T fields swapped.
func (e StationEdge) ReversedEdge() graph.Edge { return StationEdge{f: e.t, t: e.f} }

func (e StationEdge) EdgeType() int {
	return EdgeStation
}

func (e *StationEdge) ForEachTiploc(f func(node *TiplocNode)) {
	if e != nil {
		for _, t := range e.s {
			f(t)
		}
	}
}

func (e *StationEdge) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	var s []string
	for _, v := range e.s {
		s = append(s, v.Tiploc)
	}
	return util.NewXmlBuilder(encoder, start).
		AddAttribute(xml.Name{Local: "from"}, strconv.FormatInt(e.f.ID(), IdBase)).
		AddAttribute(xml.Name{Local: "to"}, strconv.FormatInt(e.t.ID(), IdBase)).
		AddAttributeIfSet(xml.Name{Local: "src"}, e.Src).
		AddCharData(strings.Join(s, ",")).
		Build()
}

func (e *StationEdge) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
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

		switch tok := token.(type) {
		case xml.CharData:
			for _, s := range strings.Split(strings.ReplaceAll(string(tok), " ", ""), ",") {
				v, err := strconv.ParseInt(s, IdBase, 64)
				if err == nil {
					e.ss = append(e.ss, v)
				}
			}

		case xml.EndElement:
			return nil
		}
	}
}
