package darwingraph

import (
	"fmt"
	"log"
	"sort"
	"strings"
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

	v := []string{fmt.Sprintf(
		"Graph state at %s\nNodes %d Edges %d\nStations %d With CRS %d\nWith Position %d Without Position %d coverage %.2f%%",
		s.ExportDate.Format(time.RFC3339),
		s.Nodes,
		s.Edges,
		s.Stations,
		s.Crs,
		s.WithPosition,
		s.WithoutPosition,
		ppc)}

	v = appendStatusSrc(v, "Name Sources", s.NameSrc)
	v = appendStatusSrc(v, "Geographic Sources", s.GeoSrc)

	log.Println(strings.Join(v, "\n"))

	return s
}

func appendStatusSrc(v []string, title string, src []StatusSrc) []string {
	if len(src) > 0 {
		v = append(v, "", title, "Source            Count", "=======================")
		for _, e := range src {
			v = append(v, fmt.Sprintf("%-16.16s %6d", e.Source, e.Count))
		}
	}
	return v
}
func (d *RailGraph) Status() Status {
	s := Status{
		ExportDate: time.Now(),
		Edges:      d.tiplocGraph.graph.Edges().Len(),
		Crs:        d.stationGraph.graph.Nodes().Len(),
	}

	locSrc := make(map[string]int)
	llSrc := make(map[string]int)

	nodes := d.tiplocGraph.graph.Nodes()
	s.Nodes = nodes.Len()

	for nodes.Next() {
		n := nodes.Node().(*TiplocNode)

		if n.Station {
			s.Stations++
		}

		if isNullIsland(n.Lon) && isNullIsland(n.Lat) {
			s.WithoutPosition++
		} else {
			s.WithPosition++
		}

		if n.LocSrc != "" {
			if _, exists := locSrc[n.LocSrc]; exists {
				locSrc[n.LocSrc] = locSrc[n.LocSrc] + 1
			} else {
				locSrc[n.LocSrc] = 1
			}
		}

		if n.LLSrc != "" {
			if _, exists := llSrc[n.LLSrc]; exists {
				llSrc[n.LLSrc] = llSrc[n.LLSrc] + 1
			} else {
				llSrc[n.LLSrc] = 1
			}
		}
	}

	s.NameSrc = appendStatusSrcMap(locSrc)
	s.GeoSrc = appendStatusSrcMap(llSrc)

	return s
}

func appendStatusSrcMap(src map[string]int) []StatusSrc {
	var ary []StatusSrc
	for k, v := range src {
		ary = append(ary, StatusSrc{Source: k, Count: v})
	}
	sort.Slice(ary, func(i, j int) bool {
		return ary[i].Source < ary[j].Source
	})
	return ary
}

type Status struct {
	ExportDate      time.Time   `json:"exportDate" xml:"exportDate,attr"`           // Time of export
	Nodes           int         `json:"nodes" xml:"nodes,attr"`                     // Number of nodes, aka tiplocs
	Edges           int         `json:"edges" xml:"edges,attr"`                     // Number of edges
	Stations        int         `json:"stations" xml:"stations,attr"`               // Number of stations
	Crs             int         `json:"crs" xml:"crs,attr"`                         // Number of entries with CRS code
	WithPosition    int         `json:"withPosition" xml:"withPosition,attr"`       // Entries with a position
	WithoutPosition int         `json:"withoutPosition" xml:"withoutPosition,attr"` // Entries with a position
	NameSrc         []StatusSrc `json:"nameSrc" xml:"nameSrc"`                      // List of sources for locations
	GeoSrc          []StatusSrc `json:"geoSrc" xml:"geoSrc"`                        // List of sources for LatLon
}

type StatusSrc struct {
	Source string `json:"source" xml:"source,attr"`
	Count  int    `json:"count" xml:",chardata"`
}
