package darwingraph

import (
	"encoding/xml"
	"github.com/peter-mount/nre-feeds/util"
	"time"
)

// RailGraph is a wrapper around a TiplocGraph & a StationGraph
type RailGraph struct {
	tiplocGraph  *TiplocGraph  // TiplocGraph of tiplocs
	stationGraph *StationGraph // StationGraph of stations
}

func NewRailGraph() *RailGraph {
	tg := NewTiplocGraph()
	return &RailGraph{
		tiplocGraph:  tg,
		stationGraph: NewStationGraph(tg),
	}
}

func (d *RailGraph) GetTiploc(tiploc string) *TiplocNode {
	return d.tiplocGraph.GetNode(tiploc)
}

func (d *RailGraph) ComputeTiplocIfAbsent(tiploc string, f func() *TiplocNode) *TiplocNode {
	return d.tiplocGraph.ComputeIfAbsent(tiploc, f)
}

// GetTiplocsForCrs returns the tiplocs associated with a CRS code or nil if none
func (d *RailGraph) GetTiplocsForCrs(crs string) []string {
	return d.tiplocGraph.GetCrs(crs)
}

// LinkTiplocs links two tiplocs together
// Returns the new TiplocEdge or nil if one already exists
func (d *RailGraph) LinkTiplocs(a, b string) *TiplocEdge {
	return d.tiplocGraph.Link(a, b)
}

func (d *RailGraph) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return util.NewXmlBuilder(e, start).
		AddAttribute(xml.Name{Local: "generated"}, time.Now().Format(time.RFC3339)).
		Append(xml.Name{Local: "tiplocs"}, d.tiplocGraph).
		Append(xml.Name{Local: "stations"}, d.stationGraph).
		Build()
}

func (d *RailGraph) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	// We ignore attributes as they are just information in the generated xml file

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			switch tok.Name.Local {
			case "tiplocs":
				err := decoder.DecodeElement(d.tiplocGraph, &tok)
				if err != nil {
					return err
				}
			case "stations":
				err := decoder.DecodeElement(d.stationGraph, &tok)
				if err != nil {
					return err
				}
			}

		case xml.EndElement:
			return nil
		}
	}
}
