package darwingraph

import (
	"log"
	"strconv"
	"strings"
	"time"
)

type lineSegment struct {
	start    int64          // Starting node id
	next     int64          // next node id
	vertices []int64        // node id's in sequence
	tpls     map[int64]bool // map to stop circular segments
}

type lineSegmentVisitor struct {
	d              *RailGraph                       // link to graph
	visitedStation map[int64]bool                   // station visited flag
	lineSegments   map[int64]map[int64]*lineSegment // map of all segments by origin
	workQueue      []*lineSegment                   // work queue of segments
	created        int                              // Number segments created
	completed      int                              // Number of segments completed
}

func (d *RailGraph) populateStations() {
	nodes := d.graph.Nodes()
	for nodes.Next() {
		n := nodes.Node().(RailNode)
		if n.NodeType() == NodeTiploc {
			t := n.(*TiplocNode)
			if t.Crs != "" {
				s := d.GetCrs(t.Crs)
				if s == nil {
					d.AddNode(d.NewStationNode(t))
				} else {
					s.addTiploc(t)
				}
			}
		}
	}
}

func (d *RailGraph) generateLineSegments() {
	visitor := &lineSegmentVisitor{
		d:              d,
		visitedStation: make(map[int64]bool),
		lineSegments:   make(map[int64]map[int64]*lineSegment),
	}

	nodes := d.graph.Nodes()
	log.Printf("Scanning %d stations as line starting points", nodes.Len())
	for nodes.Next() {
		n := nodes.Node().(RailNode)
		if n.NodeType() == NodeStation {
			visitor.visitStation(n.(*StationNode))
		}
	}
	log.Printf("%d initial line segments found", len(visitor.workQueue))

	log.Println("Beginning scan...")
	visitor.run()
	log.Printf("Completed scan, generated %d/%d segments", visitor.completed, visitor.created)

	visitor.insertIntoGraph()
}

func (lv *lineSegmentVisitor) run() {
	now := time.Now()
	notify := time.Second //* 10

	for len(lv.workQueue) > 0 {
		n := time.Now()
		if n.Sub(now) > notify {
			log.Printf("Work queue %d completed %d/%d",
				len(lv.workQueue),
				lv.completed,
				lv.created)
			now = n
		}
		ls := lv.workQueue[0]
		lv.workQueue = lv.workQueue[1:]
		lv.processLineSegment(ls)
	}
}

func (l *lineSegment) String() string {
	var v []string
	for _, n := range l.vertices {
		v = append(v, strconv.FormatInt(n, IdBase))
	}
	return "[" + strings.Join(v, ", ") + "]"
}

func (l *lineSegment) add(v int64) bool {
	if l.start == 0 {
		l.start = v
	}
	l.next = v
	l.vertices = append(l.vertices, v)

	// Return value is if we already have this tiploc - there's a couple of known loops
	_, exists := l.tpls[v]
	l.tpls[v] = true
	return exists
}

func (lv *lineSegmentVisitor) insertIntoGraph() {
	for origin, segs := range lv.lineSegments {
		f := lv.d.graph.Node(origin)
		if f == nil {
			log.Printf("Error: Line Seg origin absent %s", strconv.FormatInt(origin, IdBase))
		} else {
			for dest, seg := range segs {
				t := lv.d.graph.Node(dest)
				if t == nil {
					log.Printf("Error: Line Seg dest absent %s", strconv.FormatInt(dest, IdBase))
				} else {
					edge := &StationEdge{
						F: origin,
						T: dest,
						f: f.(RailNode),
						t: t.(RailNode),
					}
					for _, tpl := range seg.vertices {
						t := lv.d.graph.Node(tpl)
						if t == nil {
							log.Printf("Error: line seg tpl absent %s", strconv.FormatInt(tpl, IdBase))
						} else {
							edge.s = append(edge.s, t.(*TiplocNode))
						}
					}
					// First station so clone the tiplocs for station edges
					if lv.d.stations == nil {
						lv.d.initStationsGraph()
					}
					lv.d.stations.SetEdge(edge)
				}
			}
		}
	}
}

func (lv *lineSegmentVisitor) newLineSegment(a, b int64) *lineSegment {
	l := &lineSegment{
		tpls: make(map[int64]bool),
	}
	_ = l.add(a)
	_ = l.add(b)

	seg, exists := lv.lineSegments[a]
	if !exists {
		seg = make(map[int64]*lineSegment)
		lv.lineSegments[a] = seg
	}
	seg[b] = l

	lv.created++
	return l
}

func (lv *lineSegmentVisitor) visitStation(s *StationNode) {
	// Stop if we have visited this station
	if _, visited := lv.visitedStation[s.id]; visited {
		return
	}

	// mark as visited
	lv.visitedStation[s.id] = true

	// Visit each tiploc
	for _, t := range s.tiploc {
		lv.visitStationTiploc(t)
	}
}

func (lv *lineSegmentVisitor) visitStationTiploc(t *TiplocNode) {
	edges := lv.d.graph.From(t.id)
	for edges.Next() {
		edge := edges.Node().(*TiplocNode)
		// Only start a segment if we are not already doing it
		if !lv.isLineSegmentPresent(t.id, edge.id) {
			lv.visitLineSegment(lv.newLineSegment(t.id, edge.id))
		}
	}
}

func (lv *lineSegmentVisitor) isLineSegmentPresent(a, b int64) bool {
	segs, exists := lv.lineSegments[a]
	if exists {
		_, exists = segs[b]
	}
	return exists
}

func (lv *lineSegmentVisitor) terminateLineSegment(_ *lineSegment) {
	lv.completed++
}

func (lv *lineSegmentVisitor) visitLineSegment(l *lineSegment) {
	lv.workQueue = append(lv.workQueue, l)
}

func (lv *lineSegmentVisitor) processLineSegment(l *lineSegment) {

	nextNode := lv.d.graph.Node(l.next).(*TiplocNode)
	terminate := nextNode == nil || nextNode.Crs != ""
	nid := l.next
	var next []*TiplocNode
	if !terminate {
		edges := lv.d.graph.From(nid)
		for edges.Next() {
			next = append(next, edges.Node().(*TiplocNode))
		}

		switch len(next) {
		case 0:
			// No outbound edges
			terminate = true
		case 1:
			n := next[0]
			// Backtracking means we are terminating
			terminate = n.id == nid
			if !terminate {
				// Add this node's id to the segment & continue.
				// Terminate if the add indicates we already have this entry, i.e. a loop
				terminate = l.add(n.id)
			}
		default:
			// Stop this segment here
			terminate = true
			// Visit each node in the slice to start a new line
			for _, n := range next {
				lv.visitStationTiploc(n)
			}
		}
	}

	if terminate {
		// End of the line segment
		lv.terminateLineSegment(l)
	} else {
		// Visit the next entry
		lv.visitLineSegment(l)
	}

}
