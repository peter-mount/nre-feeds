package xmas

import (
  "errors"
  "flag"
  "github.com/peter-mount/golib/kernel"
  "time"
)

type XmasService struct {
  refUrl      *string             // Url of reference service
  stationsXml *string             // Stations KB xml file
  stationMap  map[string]*Station // Map of visitable stations
  tiplocMap   map[string]*Station // Map of stations by tiploc

}

func (x *XmasService) Name() string {
  return "XmasService"
}

func (x *XmasService) Init(k *kernel.Kernel) error {
  x.refUrl = flag.String("ref", "", "Url to reference microservice")
  x.stationsXml = flag.String("kb", "", "Stations knowledge base file")

  x.stationMap = make(map[string]*Station)
  x.tiplocMap = make(map[string]*Station)

  return nil
}

func (x *XmasService) Run() error {

  err := x.loadStations()
  if err != nil {
    return err
  }

  stations := x.sortStations(time.Now())
  //x.debugStations(stations)

  sched := x.createSchedule(stations)
  if sched == nil {
    return errors.New("No schedule created")
  }

  x.debugSchedule(sched)

  return nil
}
