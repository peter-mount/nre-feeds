package maps

import (
	"errors"
	"flag"
	"fmt"
	"github.com/peter-mount/go-kernel"
	"log"
	"strings"
)

// MapGenerator generates maps from various data
type MapGenerator struct {
	tileProviderName *string
	tileProvider     *tileProvider
}

func (m *MapGenerator) Name() string {
	return "MapGenerator"
}

func (m *MapGenerator) Init(k *kernel.Kernel) error {
	m.tileProviderName = flag.String("map-provider", "maplu", "The map provider, \"help\" to list available providers")

	return nil
}

func (m *MapGenerator) PostInit() error {
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

func (m *MapGenerator) Builder() *MapBuilder {
	return NewMapBuilder().
		SetTileProvider(m.tileProvider.Generator()).
		SetSize(1200, 900)
}
