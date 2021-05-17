package darwingraph

import (
	"encoding/xml"
	"github.com/peter-mount/nre-feeds/util"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"log"
	"sort"
	"strconv"
)

const (
	IdBase = 36 // Base for persisting id's, 36 = lowercase tiploc by coincidence
)

type TiplocGraph struct {
	id    int64                 // Current ID
	ids   map[string]int64      // Map of tiploc to id
	crs   map[string][]string   // Map of CRS to tiplocs (1..n relationship)
	graph *simple.DirectedGraph // Underlying directed graph
}

func NewTiplocGraph() *TiplocGraph {
	return &TiplocGraph{
		ids:   make(map[string]int64),
		crs:   make(map[string][]string),
		graph: simple.NewDirectedGraph(),
	}
}

func (d *TiplocGraph) GetCrs(crs string) []string {
	return d.crs[crs]
}

// GetNode returns an existing TiplocNode or nil if it doesn't exist
func (d *TiplocGraph) GetNode(tiploc string) *TiplocNode {
	if id, exists := d.ids[tiploc]; exists {
		return d.graph.Node(id).(*TiplocNode)
	}
	return nil
}

func (d *TiplocGraph) NextID() int64 {
	id := d.id
	d.id = d.id + 1
	return id
}

// addCrs internal call to add tiploc to a crs
func (d *TiplocGraph) addCrs(crs, tiploc string) {
	if crs != "" {
		tpls := d.GetCrs(crs)
		if tpls == nil || len(tpls) == 0 {
			d.crs[crs] = []string{tiploc}
		} else {
			for _, tpl := range tpls {
				if tpl == tiploc {
					return
				}
			}
			d.crs[crs] = append(tpls, tiploc)
		}
	}
}

// setNode internal call to set a node in the graph
func (d *TiplocGraph) setNode(n *TiplocNode) {
	d.graph.AddNode(n)
	d.ids[n.Tiploc] = n.id
	d.addCrs(n.Crs, n.Tiploc)
}

func (d *TiplocGraph) ComputeIfAbsent(tiploc string, f func() *TiplocNode) *TiplocNode {
	tn := d.GetNode(tiploc)

	if tn == nil {
		tn = f()
		if tn != nil {
			tn.Tiploc = tiploc
			// Gen ID from tiploc so base36 works here
			tn.id, _ = strconv.ParseInt(tiploc, 36, 64)
			if tn.id == 0 {
				log.Printf("id=0 for tpl \"%s\"", tn.Tiploc)
			}
			d.setNode(tn)
		}
	}

	return tn
}

// Link links two tiplocs together
// Returns the new TiplocEdge or nil if one already exists
func (d *TiplocGraph) Link(a, b string) *TiplocEdge {
	aT := d.GetNode(a)
	bT := d.GetNode(b)

	if aT != nil && bT != nil && aT.id != bT.id && !d.graph.HasEdgeBetween(aT.id, bT.id) {
		edge := &TiplocEdge{f: aT, t: bT}
		d.graph.SetEdge(edge)
		return edge
	}

	return nil
}

func (d *TiplocGraph) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	nodeName := xml.Name{Local: "node"}
	edgeName := xml.Name{Local: "edge"}

	// Get list of nodes sorted by tiploc
	var nodeAry []*TiplocNode
	nodes := d.graph.Nodes()
	for nodes.Next() {
		nodeAry = append(nodeAry, nodes.Node().(*TiplocNode))
	}
	sort.Slice(nodeAry, func(i, j int) bool {
		return nodeAry[i].Tiploc < nodeAry[j].Tiploc
	})

	// Get list of edges sorted by tiploc
	var edgeAry []graph.Edge
	edges := d.graph.Edges()
	for edges.Next() {
		edgeAry = append(edgeAry, edges.Edge())
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

	return util.NewXmlBuilder(e, start).
		AddAttribute(xml.Name{Local: "id"}, strconv.FormatInt(d.id, IdBase)).
		AddAttribute(xml.Name{Local: "nodes"}, strconv.FormatInt(int64(len(d.ids)), 10)).
		Run(func(builder *util.XmlBuilder) error {
			for _, n := range nodeAry {
				builder.Append(nodeName, n)
			}
			return nil
		}).
		Run(func(builder *util.XmlBuilder) error {
			for _, e := range edgeAry {
				builder.Append(edgeName, e)
			}
			return nil
		}).
		Build()
}

func (d *TiplocGraph) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		var err error
		switch attr.Name.Local {
		case "id":
			d.id, err = strconv.ParseInt(attr.Value, IdBase, 64)
			if err != nil {
				return err
			}
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
			case "node":
				n := &TiplocNode{}
				err := decoder.DecodeElement(n, &tok)
				if err != nil {
					return err
				}
				d.setNode(n)

			case "edge":
				e := &TiplocEdge{}
				err := decoder.DecodeElement(e, &tok)
				if err != nil {
					return err
				}
				e.f = d.graph.Node(e.F).(*TiplocNode)
				e.t = d.graph.Node(e.T).(*TiplocNode)
				log.Printf("edge f=\"%s\" t=\"%s\" f=%v t=%v", e.F, e.T, e.f, e.t)
				if e.f != nil && e.t != nil {
					d.graph.SetEdge(e)
				}
			}

		case xml.EndElement:
			return nil
		}
	}
}
