package darwingraph

import (
	"encoding/xml"
	"github.com/peter-mount/nre-feeds/util"
	"gonum.org/v1/gonum/graph/simple"
	"sort"
	"strconv"
)

// StationGraph sits on top of TiplocMap and provides a simplified graph of nodes
// which are stations.
//
// It's Edges are lines of Tiploc's which form a single path between Stations
type StationGraph struct {
	tgraph *TiplocGraph          // LinkTiplocs to underlying TiplocGraph
	graph  *simple.DirectedGraph // Underlying directed graph
}

// NewStationGraph creates a new StationGraph using data from a TiplocGraph.
// The two graph's will be linked by their nodes & edges
func NewStationGraph(tgraph *TiplocGraph) *StationGraph {
	d := &StationGraph{
		tgraph: tgraph,
		graph:  simple.NewDirectedGraph(),
	}
	return d
}

func (d *StationGraph) Populate() {
	// Populate nodes from the underlying TiplocGraph
	for crs, tiplocs := range d.tgraph.crs {
		d.graph.AddNode(d.NewStationNode(crs, d.tgraph.GetNodes(tiplocs...)))
	}
}

func (d *StationGraph) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	stationName := xml.Name{Local: "station"}

	// Get list of nodes sorted by tiploc
	var nodeAry []*StationNode
	nodes := d.graph.Nodes()
	for nodes.Next() {
		nodeAry = append(nodeAry, nodes.Node().(*StationNode))
	}
	sort.Slice(nodeAry, func(i, j int) bool {
		return nodeAry[i].Crs < nodeAry[j].Crs
	})

	return util.NewXmlBuilder(e, start).
		AddAttribute(xml.Name{Local: "stations"}, strconv.Itoa(len(nodeAry))).
		Run(func(builder *util.XmlBuilder) error {
			for _, n := range nodeAry {
				builder.Append(stationName, n)
			}
			return nil
		}).
		Build()
}

func (d *StationGraph) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	// We ignore attributes as they are just information in the generated xml file

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			switch tok.Name.Local {
			case "station":
				n := &StationNode{graph: d}
				err := decoder.DecodeElement(n, &tok)
				if err != nil {
					return err
				}
				d.graph.AddNode(n)
			}

		case xml.EndElement:
			return nil
		}
	}
}
