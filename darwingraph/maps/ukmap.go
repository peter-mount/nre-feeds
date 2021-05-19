package maps

import (
	"flag"
	"github.com/peter-mount/go-kernel"
	"github.com/peter-mount/nre-feeds/darwingraph"
	"io"
	"log"
	"os"
)

// UKMap generates a map covering the UK of the lines in a RailMap
type UKMap struct {
	darwinGraph  *darwingraph.DarwinGraph // Graph component
	mapGenerator *MapGenerator            // Map generator
	plotMap      *string                  // File to write
}

func (m *UKMap) Name() string {
	return "UKMap"
}

func (m *UKMap) Init(k *kernel.Kernel) error {
	m.plotMap = flag.String("map-uk", "", "Plot rail map on the UK")

	svce, err := k.AddService(&darwingraph.DarwinGraph{})
	if err != nil {
		return err
	}
	m.darwinGraph = (svce).(*darwingraph.DarwinGraph)

	svce, err = k.AddService(&MapGenerator{})
	if err != nil {
		return err
	}
	m.mapGenerator = (svce).(*MapGenerator)

	return nil
}

func (m *UKMap) Start() error {
	if *m.plotMap != "" {
		log.Printf("Generating UK Rail map %s", *m.plotMap)
		f, err := os.Create(*m.plotMap)
		if err != nil {
			return err
		}
		defer f.Close()
		err = m.Plot(f, 400, 300)
		if err != nil {
			return err
		}
		log.Printf("Generated %s", *m.plotMap)
	}

	return nil
}

func (m *UKMap) Plot(w io.Writer, width, height int) error {

	b := m.mapGenerator.Builder().
		Size(1200, 900).
		Zoom(6)

	m.darwinGraph.ForEachStationEdge(func(edge *darwingraph.StationEdge) {
		b.AppendStationEdge(edge)
	})

	m.darwinGraph.ForEachStationNode(func(node *darwingraph.StationNode) {
		b.AppendStation(node)
	})

	return b.Render(w)

	/*ctx := sm.NewContext()
	ctx.Size(width, height)
	ctx.AddObject(
		sm.NewMarker(
			s2.LatLngFromDegrees(52.514536, 13.350151),
			color.RGBA{0xff, 0, 0, 0xff},
			16.0,
		),
	)

	img, err := ctx.Render()
	if err != nil {
		return err
	}

	return png.Encode(w, img)*/
}
