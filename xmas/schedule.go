package xmas

import (
  "fmt"
  "github.com/peter-mount/nre-feeds/darwind3"
  "sort"
  "time"
)

const (
  longToMinutes   = (-60.0 / 15.0) // Number of minutes per degree of longitude
  NorthPoleTiploc = "NPLEINT"      // North Pole tiploc
  NorthPoleCrs    = "XNP"
  NorthPoleName   = "North Pole International"
  TimeToNorthPole = 5 * time.Minute
)

func northPole(t time.Time) *Station {
  return &Station{
    Tiploc:    NorthPoleTiploc,
    Crs:       NorthPoleCrs,
    Name:      NorthPoleName,
    Longitude: 0,
    Latitude:  90,
    tm:        t,
  }
}

// Sorts the stations based on a time for longitude 0
func (x *XmasService) sortStations(t0 time.Time) []*Station {
  var stations []*Station

  for _, station := range x.stationMap {
    station.tm = t0.Add(time.Minute * time.Duration(station.Longitude*longToMinutes))
    stations = append(stations, station)
  }

  sort.SliceStable(stations, func(i, j int) bool {
    return stations[i].tm.Before(stations[j].tm)
  })

  // Now we know the time limits add the origin & destination
  var s []*Station
  s = append(s, northPole(stations[0].tm.Add(-TimeToNorthPole)))
  s = append(s, stations...)
  s = append(s, northPole(stations[len(stations)-1].tm.Add(TimeToNorthPole)))

  return s
}

// createSchedule creates a D3 Schedule
func (x *XmasService) createSchedule(stations []*Station) *darwind3.Schedule {
  t0 := stations[0].tm

  // Convert to proper locations
  var locations []*darwind3.Location
  for _, station := range stations {
    locations = append(locations, station.toLocation())
  }

  if locations == nil {
    return nil
  }

  s := &darwind3.Schedule{
    // 2019 12 12 87 80241
    RID:              t0.Format("20060102") + "640" + t0.Format("1504"), // 64 = @ & 0 to pad time to 5 digits
    UID:              "@0" + t0.Format("1504"),                          // Invalid UID as @ not a Letter
    TrainId:          "1X02",                                            // Our head code
    Toc:              "XM",                                              // Our dummy TOC
    Status:           "P",                                               // Permanent
    TrainCat:         "PP",                                              // Parcels
    PassengerService: false,
    Active:           true,
    CancelReason:     darwind3.DisruptionReason{},
    LateReason:       darwind3.DisruptionReason{},
    Locations:        locations,
    LastReport:       darwind3.CallingPoint{},
    Formation:        darwind3.ScheduleFormation{},
    Date:             time.Now(),
  }

  // Set SSD to the start date
  s.SSD.Set(t0)

  // Fix the origin
  origin := locations[0]            // Always first location
  origin.Type = "OR"                // Origin
  origin.Times.Pta = nil            // No arrival time
  origin.Times.Wta = nil            // No arrival time
  origin.Forecast.Arrival.ET = nil  // No expected arrival
  origin.Forecast.Arrival.WET = nil // No expected arrival

  // Fix the destination
  destination := locations[len(locations)-1] // Always last location
  destination.Type = "DT"                    // Destination
  destination.Times.Ptd = nil                // No departure time
  destination.Times.Wtd = nil                // No departure time
  destination.Forecast.Departure.ET = nil    // No expected departure
  destination.Forecast.Departure.WET = nil   // No expected departure

  // Now finalise the times
  s.UpdateTime()

  return s
}

func (x *XmasService) debugSchedule(s *darwind3.Schedule) {
  fmt.Printf("Schedule:\n RID %s\n UID %s\n TID %s\nOrig %-7s %s\nDest %-7s %s\n\n",
    s.RID,
    s.UID,
    s.TrainId,
    s.Origin.Tiploc, x.tiplocMap[s.Origin.Tiploc].Name,
    s.Destination.Tiploc, x.tiplocMap[s.Destination.Tiploc].Name,
  )

  fmt.Printf(
    "%4s %-7s %4s %5s %5s %8s %8s %3s %s\n",
    "Seq",
    "Tiploc",
    "Plat",
    "Pt Ar",
    "Pt Dp",
    "Wt Arr",
    "Wt Dep",
    "CRS",
    "Name",
  )

  for i, l := range s.Locations {
    t := x.tiplocMap[l.Tiploc]
    fmt.Printf(
      "%04d %-7s %4s %5s %5s %8s %8s %3s %s\n",
      i,
      l.Tiploc,
      l.Forecast.Platform.Platform,
      l.Times.Pta,
      l.Times.Ptd,
      l.Times.Wta,
      l.Times.Wtd,
      t.Crs,
      t.Name,
    )
  }
}
