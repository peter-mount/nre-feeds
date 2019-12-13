package xmas

import (
  "encoding/xml"
  "errors"
  "github.com/peter-mount/nre-feeds/darwinref/client"
  "io/ioutil"
  "log"
  "time"
)

// Used to unmarshal the Stations Knowledge Base feed
type StationList struct {
  Stations []*Station `xml:"Station"`
}

func (x *XmasService) loadStations() error {

  // Load the Stations KB xnl
  log.Println("Loading Knowledge base")
  stationXml, err := ioutil.ReadFile(*x.stationsXml)
  if err != nil {
    return err
  }

  log.Println("Parsing xml")
  stations := &StationList{}
  err = xml.Unmarshal(stationXml, stations)
  if err != nil {
    return err
  }

  log.Println("Importing station geometry")
  for _, station := range stations.Stations {
    // ASI & SPX are duplicates in KB for the Eurostar terminals at Ashford Intl & St Pancras Intl
    // Which we already have as AFK & STP
    if !(station.Crs != "ASI" && station.Crs != "SPX") {
      log.Println(station.Crs, station)
    }
    if station.Crs != "ASI" && station.Crs != "SPX" {
      x.stationMap[station.Crs] = station
    }
  }

  log.Println("Retrieving stations from reference feed")
  refClient := &client.DarwinRefClient{Url: *x.refUrl}
  locations, err := refClient.GetStations()
  if err != nil {
    return err
  }

  log.Println("Mapping Stations to tiploc")
  for _, loc := range locations {
    if station, exists := x.stationMap[loc.Crs]; exists {
      station.Tiploc = loc.Tiploc

      // Also add to lookup map
      x.tiplocMap[loc.Tiploc] = station
    }
  }

  // Count how many don't have a tiploc
  tc := 0
  for _, station := range x.stationMap {
    if station.Tiploc == "" {
      tc++
    }
  }
  if tc > 0 {
    log.Printf("%d stations found without tiplocs", tc)
    for _, station := range x.stationMap {
      if station.Tiploc == "" {
        log.Printf("Crs %3s Name %s\n", station.Crs, station.Name)
      }
    }
    return errors.New("Failed to map all stations to a tiploc")
  }

  log.Printf("Loaded %d visitable stations", len(x.stationMap))

  // Finally add the North pole to the lookup map, no need for an accurate time here
  x.tiplocMap[NorthPoleTiploc] = northPole(time.Now())

  return nil
}
