package darwingraph

import (
	"encoding/xml"
	"github.com/peter-mount/nre-feeds/util"
	"github.com/pkg/errors"
	"log"
	"strconv"
)

// StationNode is a node in StationGraph representing a station
type StationNode struct {
	id     int64         // Unique ID
	Crs    string        // CRS of this station
	Name   string        // Name of the station
	tiploc []*TiplocNode // Tiploc nodes for this entry
	graph  *StationGraph // LinkTiplocs to parent graph
}

func (d *StationGraph) NewStationNode(crs string, tiplocs []*TiplocNode) *StationNode {
	// Nothing to set so ignore
	if crs == "" || len(tiplocs) == 0 {
		return nil
	}

	// Validate the station
	id, err := strconv.ParseInt(crs, IdBase, 64)
	if err != nil {
		log.Printf("Error: crs \"%s\" invalid", crs)
		return nil
	}

	return &StationNode{
		id:     id,
		Crs:    crs,
		Name:   tiplocs[0].Name, // Name from first tiploc
		tiploc: tiplocs,
		graph:  d,
	}
}

func (n StationNode) ID() int64 {
	return n.id
}

func (n StationNode) String() string {
	return n.Name
}

func (n *StationNode) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	tiplocName := xml.Name{Local: "tiploc"}
	return util.NewXmlBuilder(e, start).
		AddAttribute(xml.Name{Local: "id"}, strconv.FormatInt(n.id, IdBase)).
		AddAttribute(xml.Name{Local: "crs"}, n.Crs).
		AddAttribute(xml.Name{Local: "name"}, n.Name).
		Run(func(builder *util.XmlBuilder) error {
			for _, t := range n.tiploc {
				builder.Append(tiplocName, t.Tiploc)
			}
			return nil
		}).
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
		case xml.StartElement:
			switch tok.Name.Local {
			case "tiploc":
				var t string
				err := decoder.DecodeElement(&t, &tok)
				if err != nil {
					return err
				}

				// Resolve the TiplocNode
				tn := n.graph.tgraph.GetNode(t)
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
