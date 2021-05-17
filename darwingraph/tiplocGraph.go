package darwingraph

import (
	"encoding/xml"
	"github.com/peter-mount/nre-feeds/util"
	"gonum.org/v1/gonum/graph/simple"
	"log"
	"strconv"
)

type TiplocGraph struct {
	id    int64                   // Current ID
	ids   map[string]int64        // Map of tiploc to id
	crs   map[string][]string     // Map of CRS to tiplocs (1..n relationship)
	graph *simple.UndirectedGraph // Underlying graph
}

func NewTiplocGraph() *TiplocGraph {
	return &TiplocGraph{
		ids:   make(map[string]int64),
		crs:   make(map[string][]string),
		graph: simple.NewUndirectedGraph(),
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
			tn.id, _ = strconv.ParseInt(tiploc, 36, 64)
			if tn.id == 0 {
				log.Printf("id=0 for tpl \"%s\"", tn.Tiploc)
			}
			d.setNode(tn)
		}
	}

	return tn
}

func (d *TiplocGraph) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return util.NewXmlBuilder(e, start).
		AddAttribute(xml.Name{Local: "id"}, strconv.FormatInt(d.id, 36)).
		AddAttribute(xml.Name{Local: "nodes"}, strconv.FormatInt(int64(len(d.ids)), 10)).
		Run(func(builder *util.XmlBuilder) error {
			nodeName := xml.Name{Local: "node"}
			for _, v := range d.ids {
				builder.Append(nodeName, d.graph.Node(v))
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
			d.id, err = strconv.ParseInt(attr.Value, 36, 64)
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
			}

		case xml.EndElement:
			return nil
		}
	}
}
