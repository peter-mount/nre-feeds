package xmas

import (
  "fmt"
  "github.com/peter-mount/nre-feeds/darwind3"
  "github.com/peter-mount/nre-feeds/util"
  "time"
)

// Station represents a location from the knowledge base and attached to a specific tiploc
type Station struct {
  Tiploc    string    `xml:"-"`         // Main Tiploc for the crs code
  Crs       string    `xml:"CrsCode"`   // Crs code
  Name      string    `xml:"Name"`      // Name used in schedule print only
  Longitude float64   `xml:"Longitude"` // Longitude of station
  Latitude  float64   `xml:"Latitude"`  // Latitude of station
  tm        time.Time `xml:"-"`         // Time of station within our schedule
}

func (s *Station) toLocation() *darwind3.Location {
  l := &darwind3.Location{
    Type:   "IP",     // IP = Intermediate Point
    Tiploc: s.Tiploc, // Our stop
    Length: 1,        // 1 coach = sleigh
    Delay:  0,        // never delayed
  }

  // Set the arrival & departure time as the same
  pt := &util.PublicTime{}
  pt.SetTime(s.tm)
  l.Times.Pta = pt

  pt = &util.PublicTime{}
  pt.SetTime(s.tm)
  l.Times.Ptd = pt

  wt := &util.WorkingTime{}
  wt.SetTime(s.tm)
  l.Times.Wta = wt
  l.Times.Wtd = wt

  // The forecast
  l.Forecast.Arrival.ET = wt
  l.Forecast.Arrival.WET = wt
  l.Forecast.Departure.ET = wt
  l.Forecast.Departure.WET = wt

  // The platform
  l.Forecast.Platform.Platform = "SKY"
  l.Forecast.Platform.Source = "NP"
  l.Forecast.Platform.Confirmed = true

  return l
}

func (s *Station) String() string {
  return fmt.Sprintf(
    "%-7s %3s %-32s %.8f %.8f",
    s.Tiploc,
    s.Crs,
    s.Name,
    s.Longitude,
    s.Latitude,
  )
}

func (x *XmasService) debugStations(s []*Station) {
  for i, station := range s {
    fmt.Printf("%4d %s\n", i, station)
  }
}
