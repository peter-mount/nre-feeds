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
  l.Times.Pta = publicTime(s.tm)
  l.Times.Ptd = publicTime(s.tm)

  l.Times.Wta = workingTime(s.tm)
  l.Times.Wtd = workingTime(s.tm)

  // The forecast
  l.Forecast.Arrival.ET = workingTime(s.tm)
  l.Forecast.Arrival.WET = workingTime(s.tm)
  l.Forecast.Departure.ET = workingTime(s.tm)
  l.Forecast.Departure.WET = workingTime(s.tm)

  // The platform
  l.Forecast.Platform.Platform = "SKY"
  l.Forecast.Platform.Source = "NP"
  l.Forecast.Platform.Confirmed = true

  return l
}

// Public timetable cannot have anything scheduled at 00:00 so add 1 minute if after midnight
func handleMidnight(t time.Time) time.Time {
  if t.Hour() == 0 {
    t = t.Add(time.Minute)
  }
  return t
}

func publicTime(t time.Time) *util.PublicTime {
  pt := &util.PublicTime{}
  pt.SetTime(handleMidnight(t))
  return pt
}

func workingTime(t time.Time) *util.WorkingTime {
  wt := &util.WorkingTime{}
  wt.SetTime(handleMidnight(t))
  return wt
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
