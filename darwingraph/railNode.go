package darwingraph

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"log"
	"strconv"
)

const (
	NodeTiploc  = iota // Node is for a Tiploc
	NodeStation        // Node is for a Station
)

type RailNode interface {
	graph.Node
	NodeType() int
	Clone() RailNode
}

// AddNode internal call to set a node in the graph
func (d *RailGraph) AddNode(n RailNode) {
	if n != nil {
		d.graph.AddNode(n)
	}
}

// Mirrors all tiplocs in the main graph.
// Needed as StationEdge cannot share the same graph as TiplocEdge when we have
// just 2 points in the edge.
func (d *RailGraph) initStationsGraph() {
	d.stations = simple.NewDirectedGraph()
	nodes := d.graph.Nodes()
	for nodes.Next() {
		n := nodes.Node().(RailNode)
		if n.NodeType() == NodeTiploc {
			d.stations.AddNode(n.Clone())
		}
	}
}

func (d *RailGraph) ComputeIfAbsent(id string, f func() RailNode) RailNode {
	tn := d.GetNode(id)

	if tn == nil {
		tn = f()
		if tn != nil {
			id64, _ := strconv.ParseInt(id, 36, 64)
			if id64 != tn.ID() {
				log.Printf("id missmatch type %d: %s != %s",
					tn.NodeType(),
					strconv.FormatInt(tn.ID(), IdBase),
					id)
			}
			d.AddNode(tn)
		}
	}

	return tn
}

// GetNode returns an existing TiplocNode or nil if it doesn't exist
func (d *RailGraph) GetNode(s string) RailNode {
	id, err := strconv.ParseInt(s, IdBase, 64)
	if err != nil {
		return nil
	}
	n := d.graph.Node(id)
	if n != nil {
		return n.(RailNode)
	}
	return nil
}

func (d *RailGraph) GetTiploc(tiploc string) *TiplocNode {
	n := d.GetNode(tiploc)
	if n != nil && n.NodeType() == NodeTiploc {
		return n.(*TiplocNode)
	}
	return nil
}

func (d *RailGraph) GetCrs(crs string) *StationNode {
	n := d.GetNode(crs)
	if n != nil && n.NodeType() == NodeStation {
		return n.(*StationNode)
	}
	return nil
}
