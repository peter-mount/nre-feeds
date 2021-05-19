package darwingraph

import (
	"encoding/xml"
	"github.com/peter-mount/nre-feeds/util"
	"github.com/pkg/errors"
	"log"
	"strconv"
	"strings"
)

// StationNode is a node in StationGraph representing a station
type StationNode struct {
	id     int64         // Unique ID
	Crs    string        // CRS of this station
	Name   string        // Name of the station
	tiploc []*TiplocNode // Tiploc nodes for this entry
	graph  *RailGraph    // LinkTiplocs to parent graph
}

func (d *RailGraph) NewStationNode(tiploc *TiplocNode) *StationNode {
	// Nothing to set so ignore
	if tiploc == nil || tiploc.Crs == "" {
		return nil
	}

	// Validate the station
	id, err := strconv.ParseInt(tiploc.Crs, IdBase, 64)
	if err != nil {
		log.Printf("Error: crs \"%s\" invalid", tiploc.Crs)
		return nil
	}

	return &StationNode{
		id:     id, // Crs codes shouldn't clash with Tiplocs as crs are 3 chars & Tiplocs are 4-7.
		Crs:    tiploc.Crs,
		Name:   tiploc.Name, // Name from first tiploc
		tiploc: []*TiplocNode{tiploc},
		graph:  d,
	}
}

func (n StationNode) ID() int64 {
	return n.id
}

func (n StationNode) String() string {
	return n.Name
}

func (n StationNode) NodeType() int {
	return NodeStation
}

func (n *StationNode) addTiploc(tiploc *TiplocNode) {
	if tiploc != nil && tiploc.Crs == n.Crs {
		n.tiploc = append(n.tiploc, tiploc)
	}
}

func (n *StationNode) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var s []string
	for _, t := range n.tiploc {
		s = append(s, t.Tiploc)
	}
	return util.NewXmlBuilder(e, start).
		AddAttribute(xml.Name{Local: "id"}, strconv.FormatInt(n.id, IdBase)).
		AddAttribute(xml.Name{Local: "crs"}, n.Crs).
		AddAttribute(xml.Name{Local: "name"}, n.Name).
		AddCharData(strings.Join(s, ",")).
		Build()
}

func (n *StationNode) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		var err error
		switch attr.Name.Local {
		case "id":
			n.id, err = strconv.ParseInt(attr.Value, IdBase, 64)
		case "crs":
			n.Crs = attr.Value
		case "name":
			n.Name = attr.Value
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
			for _, t := range strings.Split(strings.ReplaceAll(string(tok), " ", ""), ",") {
				tn := n.graph.GetTiploc(t)
				if tn == nil {
					return errors.Errorf("tiploc \"%s\" not found for \"%s\"", t, n.Crs)
				}

				n.tiploc = append(n.tiploc, tn)
			}
		case xml.EndElement:
			return nil
		}
	}
}
