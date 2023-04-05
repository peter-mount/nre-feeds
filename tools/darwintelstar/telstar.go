package telstar

import (
	"fmt"
	"github.com/peter-mount/nre-feeds/ldb/client"
	"github.com/peter-mount/nre-feeds/ldb/service"
	"github.com/peter-mount/nre-feeds/util/telstar"
)

type Telstar struct {
	PageNumber *int    `kernel:"flag,page,Page number"`
	Crs        *string `kernel:"flag,crs,crs"`
}

func (t *Telstar) Start() error {
	/*if *t.PageNumber > 0 {
		return t.boards()
	}*/
	if len(*t.Crs) == 3 {
		return t.crs(*t.Crs)
	}
	return nil
}

func (t *Telstar) boards() error {
	pn := *t.PageNumber
	var crsa []byte
	for pn > 64 {
		crsa = append([]byte{byte(pn % 100)}, crsa...)
		pn = pn / 100
	}
	return t.crs(string(crsa))
}

func (t *Telstar) crs(crs string) error {

	crs = "LBG"

	cl := client.DarwinLDBClient{Url: "https://ldb.prod.a51.li"}
	result, err := cl.GetSchedule(crs)
	if err != nil {
		return err
	}

	response := telstar.NewResponse().
		Dynamic().
		PageNumber(*t.PageNumber).
		FrameId('a')

	f := newDepartureFrame(response, result)

	for _, m := range result.Messages {

		f.NewLine().
			Print(m.Message)

		f.Build()
		f = newDepartureFrame(response, result)
	}

	count := 0
	for _, departure := range result.Services {

		loc := departure.Location
		if !loc.IsDestination() {

			if count > 0 && (count%10) == 0 {
				f.Build()
				f = newDepartureFrame(response, result)
			}

			dest := departure.Dest
			destName := dest.Tiploc
			tpl, _ := result.Tiplocs.Get(dest.Tiploc)
			if tpl != nil {
				destName = tpl.Name
			}

			plat := ""
			if !(loc.Forecast.Platform.CISSuppressed || loc.Forecast.Platform.Suppressed) {
				plat = loc.Forecast.Platform.Platform
			}

			f.White().
				DoubleHeight().
				Printf("%-23.23s", destName).
				White().Printf("%2.2s", plat).
				White().Printf("%5.5s %5.5s", loc.Time.String()[:5], "On Time").
				NewLine()

			count++
		}
	}

	f.Route(10, *t.PageNumber).
		Build()

	s, err := response.Build()
	if err == nil {
		fmt.Print(s)
	}
	return err
}

func newDepartureFrame(response *telstar.Response, result *service.StationResult) *telstar.FrameBuilder {

	stationTiploc, _ := result.Tiplocs.Get(result.Station[0])

	return response.NewFrame().
		Header("DepartureBoards.mobi").
		NavMessage("[B][n][W]Data as of %s", result.Date.Format("2006 Jan 2 15:04:05")).
		Printf("[B][n][W][D]%s\r\n\n", stationTiploc.Name)
}
