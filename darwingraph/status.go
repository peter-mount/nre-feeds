package darwingraph

import (
	"log"
	"time"
)

func (d *DarwinGraph) Status() Status {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	s := d.graph.Status()

	// Percentage of nodes with a position
	ppc := 0.0
	if s.Nodes > 0 {
		ppc = float64(s.WithPosition) * 100.0 / float64(s.Nodes)
	}

	log.Printf(
		"Graph state at %s\nNodes %d Edges %d\nStations %d With CRS %d\nWith Position %d Without Position %d coverage %.2f%%",
		s.ExportDate.Format(time.RFC3339),
		s.Nodes,
		s.Edges,
		s.Stations,
		s.Crs,
		s.WithPosition,
		s.WithoutPosition,
		ppc)

	return s
}

func (d *TiplocGraph) Status() Status {
	s := Status{
		ExportDate: time.Now(),
		Edges:      d.graph.Edges().Len(),
	}

	nodes := d.graph.Nodes()
	s.Nodes = nodes.Len()

	for nodes.Next() {
		n := nodes.Node().(*TiplocNode)

		if n.Crs != "" {
			s.Crs++
		}

		if n.Station {
			s.Stations++
		}

		if isNullIsland(n.Lon) && isNullIsland(n.Lat) {
			s.WithoutPosition++
		} else {
			s.WithPosition++
		}
	}

	return s
}

type Status struct {
	ExportDate      time.Time `json:"exportDate" xml:"exportDate,attr"`           // Time of export
	Nodes           int       `json:"nodes" xml:"nodes,attr"`                     // Number of nodes, aka tiplocs
	Edges           int       `json:"edges" xml:"edges,attr"`                     // Number of edges
	Stations        int       `json:"stations" xml:"stations,attr"`               // Number of stations
	Crs             int       `json:"crs" xml:"crs,attr"`                         // Number of entries with CRS code
	WithPosition    int       `json:"withPosition" xml:"withPosition,attr"`       // Entries with a position
	WithoutPosition int       `json:"withoutPosition" xml:"withoutPosition,attr"` // Entries with a position
}
