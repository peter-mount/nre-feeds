package xmas

import (
	"errors"
	"flag"
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/nre-feeds/util"
	"time"
)

type XmasService struct {
	refUrl      *string             // Url of reference service
	stationsXml *string             // Stations KB xml file
	xmas        *bool               // Flag to set things for the real XMas run
	stationMap  map[string]*Station // Map of visitable stations
	tiplocMap   map[string]*Station // Map of stations by tiploc

}

func (x *XmasService) Name() string {
	return "XmasService"
}

func (x *XmasService) Init(k *kernel.Kernel) error {
	x.refUrl = flag.String("ref", "", "Url to reference microservice")
	x.stationsXml = flag.String("kb", "", "Stations knowledge base file")
	x.xmas = flag.Bool("xmas", false, "Set defaults for XMas")

	x.stationMap = make(map[string]*Station)
	x.tiplocMap = make(map[string]*Station)

	return nil
}

func (x *XmasService) Run() error {

	var t0 time.Time

	if *x.xmas {
		tYear := time.Now()
		t0 = time.Date(tYear.Year(), 12, 25, 0, 0, 0, 0, util.London())
	} else {
		t0 = time.Now()
	}

	err := x.loadStations()
	if err != nil {
		return err
	}

	stations := x.sortStations(t0)
	//x.debugStations(stations)

	sched := x.createSchedule(stations)
	if sched == nil {
		return errors.New("No schedule created")
	}

	x.debugSchedule(sched)

	return nil
}
