package index

import (
	"github.com/peter-mount/go-kernel/v2/log"
	"github.com/peter-mount/nre-feeds/darwinref"
	refClient "github.com/peter-mount/nre-feeds/darwinref/client"
	strings2 "github.com/peter-mount/nre-feeds/util/strings"
	"sort"
	"strings"
)

// Index creates the station index pages.
// This is called once in a while to keep the index in sync with new rail stations.
type Index struct {
	Api       *Api `kernel:"inject"`
	BaseFrame *int `kernel:"flag,index-base,Base frame number"`
	refClient refClient.DarwinRefClient
}

func (i *Index) Start() error {
	if *i.BaseFrame > 0 {
		i.refClient.Url = "https://ref.prod.a51.li"
		return i.rebuildIndex()
	}
	return nil
}

func (i *Index) rebuildIndex() error {

	if err := i.Api.Login(); err != nil {
		return err
	}

	log.Println("Retrieving current station listings")
	stations, err := i.refClient.GetStations()
	if err != nil {
		return err
	}

	// Filter out non-stations and ensure we are unique by CRS
	// as some CRS codes can at times get duplicated
	var a []*darwinref.Location
	crs := make(map[string]*darwinref.Location)
	for _, s := range stations {
		if s.IsPublic() && s.Station {
			if _, exists := crs[s.Crs]; !exists {
				a = append(a, s)
				crs[s.Crs] = s
			}
		}
	}
	stations = a

	// Sort by name
	sort.SliceStable(stations, func(i, j int) bool {
		a := strings.ToLower(stations[i].Name)
		b := strings.ToLower(stations[j].Name)
		return a < b
	})

	// Create map keyed by first letter
	m := make(map[string][]*darwinref.Location)
	for _, s := range stations {
		key := strings.ToLower(s.Name)[0:1]
		a := m[key]
		a = append(a, s)
		m[key] = a
	}

	// Get the letters
	var letters []string
	for k, _ := range m {
		letters = append(letters, k)
	}
	strings2.SortLower(letters)

	log.Println(letters)
	return nil
}

func (i *Index) generateIndexPage() error {
	return nil
}
