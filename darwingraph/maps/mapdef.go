package maps

import (
	"github.com/peter-mount/nre-feeds/darwingraph"
)

type MapDef struct {
	Name    string  // Name of map in -map parameter
	Handler MapTask // Function to plot map
}

func (m *MapService) initMapDefs() {
	m.mapDefs = []MapDef{
		{Name: "uk", Handler: m.ukMap},
		{Name: "se", Handler: m.seMap},
		{Name: "london", Handler: m.londonMap},
	}
}

// ukMap generate a map of the entire UK
func (m *MapService) ukMap(b *MapBuilder) {
	b.Size(600, 710).
		Center(-3.878, 54.413).
		Zoom(6).
		ForEachStationEdge(m.darwinGraph, func(b *MapBuilder, e *darwingraph.StationEdge) {
			b.AppendStationEdge(e)
		}).
		ForEachStationNode(m.darwinGraph, func(b *MapBuilder, n *darwingraph.StationNode) {
			b.AppendStation(n)
		})
}

// seMap generate a map of the South East
func (m *MapService) seMap(b *MapBuilder) {
	b.Size(1000, 710).
		Center(0.25, 51.25).
		Zoom(9).
		StationRadius(125).
		ForEachStationEdge(m.darwinGraph, func(b *MapBuilder, e *darwingraph.StationEdge) {
			b.AppendStationEdge(e)
		}).
		ForEachStationNode(m.darwinGraph, func(b *MapBuilder, n *darwingraph.StationNode) {
			b.AppendStation(n)
		})
}

// londonMap generate a map of London
func (m *MapService) londonMap(b *MapBuilder) {
	b.Size(1000, 710).
		Center(-.125, 51.5).
		Zoom(11).
		StationRadius(30).
		ForEachStationEdge(m.darwinGraph, func(b *MapBuilder, e *darwingraph.StationEdge) {
			b.AppendStationEdge(e)
		}).
		ForEachStationNode(m.darwinGraph, func(b *MapBuilder, n *darwingraph.StationNode) {
			b.AppendStation(n)
		})
}
