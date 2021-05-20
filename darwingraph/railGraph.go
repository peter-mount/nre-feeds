package darwingraph

import (
	"encoding/xml"
	"github.com/peter-mount/nre-feeds/util"
	"gonum.org/v1/gonum/graph/simple"
	"sort"
	"strconv"
)

const (
	IdBase = 36 // Base for persisting id's, 36 = lowercase tiploc by coincidence
)

// RailGraph is a wrapper around a TiplocGraph & a StationGraph
type RailGraph struct {
	graph    *simple.DirectedGraph // Underlying directed graph
	stations *simple.DirectedGraph // Copy of graph with just Tiplocs, used for Station Edges
}

func NewRailGraph() *RailGraph {
	return &RailGraph{
		graph:    simple.NewDirectedGraph(),
		stations: simple.NewDirectedGraph(),
	}
}

func (d *RailGraph) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	nodeName := xml.Name{Local: "node"}
	edgeName := xml.Name{Local: "edge"}
	stationName := xml.Name{Local: "station"}
	lineName := xml.Name{Local: "line"}

	// Get list of nodes sorted by tiploc
	var nodeAry []*TiplocNode
	var stnAry []*StationNode
	nodes := d.graph.Nodes()
	for nodes.Next() {
		n := nodes.Node().(RailNode)
		switch n.NodeType() {
		case NodeTiploc:
			nodeAry = append(nodeAry, n.(*TiplocNode))
		case NodeStation:
			stnAry = append(stnAry, n.(*StationNode))
		}
	}

	sort.Slice(nodeAry, func(i, j int) bool {
		return nodeAry[i].Tiploc < nodeAry[j].Tiploc
	})

	sort.Slice(stnAry, func(i, j int) bool {
		return stnAry[i].Crs < stnAry[j].Crs
	})

	// Get list of edges sorted by tiploc
	var edgeAry []*TiplocEdge
	edges := d.graph.Edges()
	for edges.Next() {
		edge := edges.Edge().(*TiplocEdge)
		edgeAry = append(edgeAry, edge)
	}

	sort.Slice(edgeAry, func(i, j int) bool {
		af := edgeAry[i].From().(*TiplocNode).Tiploc
		bf := edgeAry[j].From().(*TiplocNode).Tiploc
		if af == bf {
			af = edgeAry[i].To().(*TiplocNode).Tiploc
			bf = edgeAry[j].To().(*TiplocNode).Tiploc
		}
		return af < bf
	})

	var lineAry []*StationEdge
	edges = d.stations.Edges()
	for edges.Next() {
		edge := edges.Edge().(*StationEdge)
		lineAry = append(lineAry, edge)
	}

	sort.Slice(lineAry, func(i, j int) bool {
		af := lineAry[i].From().(*TiplocNode).Tiploc
		bf := lineAry[j].From().(*TiplocNode).Tiploc
		if af == bf {
			af = lineAry[i].To().(*TiplocNode).Tiploc
			bf = lineAry[j].To().(*TiplocNode).Tiploc
		}
		return af < bf
	})

	return util.NewXmlBuilder(e, start).
		Run(func(builder *util.XmlBuilder) error {
			for _, n := range nodeAry {
				builder.Append(nodeName, n)
			}
			return nil
		}).
		Run(func(builder *util.XmlBuilder) error {
			for _, n := range stnAry {
				builder.Append(stationName, n)
			}
			return nil
		}).
		Run(func(builder *util.XmlBuilder) error {
			for _, e := range edgeAry {
				builder.Append(edgeName, e)
			}
			return nil
		}).
		Run(func(builder *util.XmlBuilder) error {
			for _, n := range lineAry {
				builder.Append(lineName, n)
			}
			return nil
		}).
		Build()
}

func (d *RailGraph) UnmarshalXML(decoder *xml.Decoder, _ xml.StartElement) error {
	d.graph = simple.NewDirectedGraph()
	d.stations = nil

	// We ignore attributes as they are just information in the generated xml file

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			switch tok.Name.Local {
			case "node":
				n := &TiplocNode{}
				err := decoder.DecodeElement(n, &tok)
				if err != nil {
					return err
				}
				n.id, _ = strconv.ParseInt(n.Tiploc, IdBase, 64)
				d.AddNode(n)

			case "edge":
				e := &TiplocEdge{}
				err := decoder.DecodeElement(e, &tok)
				if err != nil {
					return err
				}
				e.f = d.graph.Node(e.F).(*TiplocNode)
				e.t = d.graph.Node(e.T).(*TiplocNode)
				if e.f != nil && e.t != nil {
					d.graph.SetEdge(e)
				}

			case "station":
				// First station so clone the tiplocs for station edges
				if d.stations == nil {
					d.initStationsGraph()
				}
				n := &StationNode{graph: d}
				err := decoder.DecodeElement(n, &tok)
				if err != nil {
					return err
				}
				n.id, _ = strconv.ParseInt(n.Crs, IdBase, 64)
				d.AddNode(n)

			case "line":
				// First station so clone the tiplocs for station edges
				if d.stations == nil {
					d.initStationsGraph()
				}
				e := &StationEdge{}
				err := decoder.DecodeElement(e, &tok)
				if err != nil {
					return err
				}
				e.f = d.graph.Node(e.F).(*TiplocNode)
				e.t = d.graph.Node(e.T).(*TiplocNode)
				for _, tpl := range e.ss {
					n := d.graph.Node(tpl).(*TiplocNode)
					if n != nil {
						e.s = append(e.s, n)
					}
				}
				if e.f != nil && e.t != nil {
					d.stations.SetEdge(e)
				}
			}

		case xml.EndElement:

			return nil
		}
	}
}
