package darwingraph

import "gonum.org/v1/gonum/graph"

const (
	EdgeTiploc  = iota // Edge is for a Tiploc
	EdgeStation        // Edge is for a line segment
)

type RailEdge interface {
	graph.Edge
	EdgeType() int
}

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
