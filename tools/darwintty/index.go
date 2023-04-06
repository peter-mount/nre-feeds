package darwintty

import (
	"fmt"
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/darwinref"
	"github.com/peter-mount/nre-feeds/tools/darwintty/render"
	"sort"
	"strings"
)

func (s *Server) refreshIndex() error {
	fmt.Println("Refreshing station index")
	stations, err := s.refClient.GetStations()
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

	fmt.Printf("Refreshed %d stations\n", len(stations))

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.stations = m
	s.stationCount = len(stations)

	return nil
}

func (s *Server) getStationCount() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.stationCount
}

func (s *Server) getStations(k string) []*darwinref.Location {
	if k == "" {
		return nil
	}
	k = strings.ToLower(k)[0:1]

	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.stations[k]
}

func (s *Server) getStationIndex() []string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var a []string
	for k, _ := range s.stations {
		a = append(a, k)
	}

	sort.SliceStable(a, func(i, j int) bool {
		return a[i] < a[j]
	})

	return a
}

func (s *Server) index(r *rest.Rest) error {

	prefix := strings.ToLower(r.Var("prefix"))
	if len(prefix) != 1 || prefix[0] < 'a' || prefix[0] > 'z' {
		prefix = ""
	}

	crsPrefix := "http://" + r.Request().Host + "/"
	crsLen := len(crsPrefix) + 3 // +3 for crs

	urlPrefix := "http://" + r.Request().Host + "/index/"
	prefixLen := len(urlPrefix) + 1 // +1 for letter

	b := render.New()
	if prefix == "" {
		b = b.Println("CRS Index").
			NewLine().
			Print(topLeft).
			Repeat(horiz, 3).
			Print(bottomLower).
			Repeat(horiz, prefixLen+2).
			Print(topRight).
			NewLine()

		for _, k := range s.getStationIndex() {
			b = b.Print(vertical).
				Printf(" %s ", k).
				Print(vertical).
				Print(" ").
				Link(urlPrefix + k).
				Print(" ").
				Print(vertical).
				NewLine()
		}

		b = b.Print(bottomLeft).
			Repeat(horiz, 3).
			Print(topUpper).
			Repeat(horiz, prefixLen+2).
			Print(bottomRight).
			NewLine()

	} else {

		index := s.getStations(prefix)
		max := 0
		for _, l := range index {
			max = Max(max, len(l.Name))
		}

		b = b.Printf("CRS Index for %s", strings.ToUpper(prefix)).
			NewLine().
			NewLine().
			Print(topLeft).
			Repeat(horiz, max+2).
			Print(bottomLower).
			Repeat(horiz, 5).
			Print(bottomLower).
			Repeat(horiz, crsLen+2).
			Print(topRight).
			NewLine()

		f1 := fmt.Sprintf(" %%-%d.%ds ", max, max)
		for _, l := range index {
			b = b.Print(vertical).
				Printf(f1, l.Name).
				Print(vertical).
				Printf(" %-3.3s ", l.Crs).
				Print(vertical).
				Print(" ").
				Link(crsPrefix + strings.ToLower(l.Crs)).
				Print(" ").
				Print(vertical).
				NewLine()
		}

		b = b.Print(bottomLeft).
			Repeat(horiz, max+2).
			Print(topUpper).
			Repeat(horiz, 5).
			Print(topUpper).
			Repeat(horiz, crsLen+2).
			Print(bottomRight).
			NewLine()

	}

	return s.respond(r, b)
}
