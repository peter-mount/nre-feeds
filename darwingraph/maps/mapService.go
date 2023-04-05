package maps

import (
	"errors"
	"flag"
	"fmt"
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/nre-feeds/darwingraph"
	"log"
	"os"
	"strings"
)

// MapService generates a map covering the UK of the lines in a RailMap
type MapService struct {
	darwinGraph      *darwingraph.DarwinGraph // Graph component
	mapName          *string                  // Map to plot
	plotMap          *string                  // File to write
	tileProviderName *string
	tileProvider     *tileProvider
	mapDefs          []MapDef // Map definitions
}

func (m *MapService) Name() string {
	return "MapService"
}

func (m *MapService) Init(k *kernel.Kernel) error {
	m.mapName = flag.String("map", "", "Map to plot, help for available maps")
	m.plotMap = flag.String("map-file", "", "File to contain generated map")
	m.tileProviderName = flag.String("map-provider", "maplu", "The map provider, \"help\" to list available providers")

	svce, err := k.AddService(&darwingraph.DarwinGraph{})
	if err != nil {
		return err
	}
	m.darwinGraph = (svce).(*darwingraph.DarwinGraph)

	m.initMapDefs()

	return nil
}

func (m *MapService) PostInit() error {
	if *m.tileProviderName != "" {
		if *m.tileProviderName == "help" {
			s := []string{"Available tile providers:"}
			for _, p := range tileProviders {
				s = append(s, fmt.Sprintf("%24.24s %s", p.Name, p.Title))
			}
			fmt.Println(strings.Join(s, "\n"))
			return errors.New("Abort")
		}

		for _, p := range tileProviders {
			if p.Name == *m.tileProviderName {
				// Must set local var first then set m.tileProvider otherwise we get the wrong pointer
				// and always use the last entry in tileProviders
				n := p
				m.tileProvider = &n
			}
		}
	}

	if m.tileProvider == nil {
		m.tileProvider = &tileProviders[0]
	}

	log.Printf("Using %s map provider: %s", m.tileProvider.Name, m.tileProvider.Title)
	return nil
}

func (m *MapService) Start() error {
	if *m.plotMap != "" {
		var def *MapDef
		for _, d := range m.mapDefs {
			if d.Name == *m.mapName {
				// Make copy to get pointer otherwise we'll always get the last entry
				t := d
				def = &t
			}
		}

		if def == nil {
			var s []string
			for _, d := range m.mapDefs {
				s = append(s, d.Name)
			}
			fmt.Printf("Available maps: %s\n", strings.Join(s, ", "))
			return errors.New("-map is required")
		}

		log.Printf("Generating %s to %s", def.Name, *m.plotMap)

		f, err := os.Create(*m.plotMap)
		if err != nil {
			return err
		}
		defer f.Close()

		err = NewMapBuilder().
			TileProvider(m.tileProvider.Generator()).
			Size(1200, 900).
			Zoom(7).
			Run(def.Handler).
			Render(f)
		if err != nil {
			return err
		}

		log.Printf("Generated %s to %s", def.Name, *m.plotMap)
	}

	return nil
}
