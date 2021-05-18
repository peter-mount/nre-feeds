package darwingraph

import (
	"gonum.org/v1/gonum/graph/path"
	"log"
	"strings"
)

func (d *DarwinGraph) Neighbours(tiploc string) []string {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.graph.tiplocGraph.Neighbours(tiploc)
}

func (d *TiplocGraph) Neighbours(tiploc string) []string {
	var t []string
	for _, n := range d.NeighbouringNodes(d.GetNode(tiploc)) {
		t = append(t, n.Tiploc)
	}
	return t
}

func (d *TiplocGraph) NeighbouringNodes(n *TiplocNode) []*TiplocNode {
	if n == nil {
		return nil
	}

	var t []*TiplocNode
	nodes := d.graph.From(n.ID())
	for nodes.Next() {
		t = append(t, nodes.Node().(*TiplocNode))
	}
	return t
}

func (d *DarwinGraph) test() {
	from := d.GetTiplocNode("MSTONEE")
	to := d.GetTiplocNode("MSTONEW")

	log.Printf("Searching %s to %s", from.Tiploc, to.Tiploc)
	pth := path.DijkstraAllFrom(from, d.graph.tiplocGraph.graph)
	log.Println("AllTo")
	ladders, weight := pth.AllTo(to.ID())
	log.Printf("AllTo ladders %d weight %f", len(ladders), weight)

	for i, l := range ladders {
		var s []string
		for _, t := range l {
			s = append(s, t.(*TiplocNode).Name)
		}
		log.Printf("%04d %04d %s", i, len(l), strings.Join(s, "\n          "))
	}
}
