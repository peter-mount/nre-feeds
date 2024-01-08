package boards

import (
	"fmt"
	"github.com/peter-mount/nre-feeds/ldb/client"
	"github.com/peter-mount/nre-feeds/ldb/service"
	"github.com/peter-mount/nre-feeds/util/telstar"
	"strings"
)

// Departures handles serving realtime departure board frames
// This is called from Telstar when a relevant page is requested
type Departures struct {
	PageNumber *int    `kernel:"flag,page,Page number"`
	Crs        *string `kernel:"flag,crs,crs"`
}

func (t *Departures) Start() error {
	/*if *t.PageNumber > 0 {
		return t.boards()
	}*/
	if len(*t.Crs) == 3 {
		return t.crs(*t.Crs)
	}
	return nil
}

func (t *Departures) boards() error {
	pn := *t.PageNumber
	var crsa []byte
	for pn > 64 {
		crsa = append([]byte{byte(pn % 100)}, crsa...)
		pn = pn / 100
	}
	return t.crs(string(crsa))
}

func (t *Departures) crs(crs string) error {

	//crs = "LBG"

	cl := client.DarwinLDBClient{Url: "https://ldb.prod.a51.li"}
	result, err := cl.GetSchedule(crs)
	if err != nil {
		return err
	}

	response := telstar.NewResponse().
		Dynamic().
		PageNumber(*t.PageNumber).
		FrameId('a')

	f := t.newDepartureFrame(response, result)

	for _, m := range result.Messages {

		f.NewLine()

		s := m.Message
		if i := strings.Index(s, "More detail"); i > -1 {
			s = s[:i]
		}
		for i, l := range strings.Split(s, ".") {
			l = strings.TrimSpace(l)
			if l != "" {
				for _, s := range telstar.Split(l+".", 38) {
					if i == 0 {
						f.White()
					} else {
						f.Yellow()
					}
					f.Println(s)
				}
			}
			f.NewLine()
		}

		f.Build()
		f = t.newDepartureFrame(response, result)
	}

	count := 0
	for _, departure := range result.Services {

		loc := departure.Location
		if !loc.IsDestination() {

			dest := departure.Dest
			destName := dest.Tiploc
			if loc.FalseDestination != "" {
				destName = loc.FalseDestination
			}
			tpl, _ := result.Tiplocs.Get(destName)
			if tpl != nil {
				destName = tpl.Name
			}

			via := ""
			if v, exists := result.Via[departure.RID]; exists {
				via = v.Text
			}

			plat := ""
			if !(loc.Forecast.Platform.CISSuppressed || loc.Forecast.Platform.Suppressed) {
				plat = loc.Forecast.Platform.Platform
			}

			var reasonText []string
			if result.Reasons != nil {
				reason := result.Reasons.Late[departure.LateReason.Reason]
				if loc.Cancelled && departure.CancelReason.Reason != 0 {
					reason = result.Reasons.Cancelled[departure.CancelReason.Reason]
				}
				if reason != nil {
					reasonText = telstar.Split(reason.Text, 38)
				}
			}

			toc := departure.Toc
			if tocName, _ := result.Tocs.Get(toc); tocName != nil {
				toc = tocName.Name
			}

			lastName := departure.LastReport.Tiploc
			if lastName != "" {
				if lastName == departure.Location.Tiploc {
					// Don't show for this location
					lastName = ""
				} else {
					// resolve the name
					tpl, _ = result.Tiplocs.Get(lastName)
					if tpl != nil {
						lastName = tpl.Name
					}
				}
			}

			// Work out departure height
			height := 2

			// The separator row
			if count > 0 {
				height++
			}

			if via != "" {
				height++
			}

			height += len(reasonText)

			if toc != "" || loc.Length > 0 {
				height++
			}

			if lastName != "" {
				height++
			}

			// Page break if not enough room

			if (count + height) >= 20 {
				f.Build()
				f = t.newDepartureFrame(response, result)
				count = 0
			}

			// Now render the departure

			if count > 0 {
				f.MosaicBlue().SepGraphSolidMid().Print("\r")
			}

			if loc.Cancelled {
				f.Red()
			} else if loc.Forecast.Arrived {
				f.White()
			} else {
				f.Green()
			}
			f.DoubleHeight().
				Printf("%-28.28s", destName)

			if loc.Cancelled {
				f.Red().Flash().Print("Canceled")
			} else {
				f.Printf(" %2.2s", plat)

				if loc.Forecast.Arrived {
					f.White().Print("Arrvd")
				} else if loc.Forecast.Delayed {
					f.Red().Print("Delyd")
				} else {
					/*t := loc.Time.Time(time.Now())
					dt := t.Sub(time.Now()).Minutes()
					if dt > 0 && dt < 10 {
						f.Yellow().Printf("%1.1d min", int(dt))
					} else*/{
						f.Green().Printf("%5.5s", loc.Time.String()[:5])
					}
				}
			}
			f.NewLine().NewLine()

			if via != "" {
				f.Yellow().Println(via)
			}

			for _, s := range reasonText {
				f.Yellow().Println(s)
			}

			if toc != "" || loc.Length > 0 {
				if toc != "" {
					f.Yellow().Printf("%s service ", toc)
				}
				if loc.Length > 0 {
					f.Yellow().Printf("%d coaches", loc.Length)
				}
				f.NewLine()
			}

			if lastName != "" {
				s := "Last seen " + departure.LastReport.Time.String()[:5]
				/*if departure.LastReport.Departed {
					s = s + " departed"
				} else {
					s = s + " at"
				}*/
				s = s + " " + lastName
				f.Yellow().
					Printf("%-38.38s", s).
					NewLine()
			}

			// Move to next entry
			count = count + height
		}
	}

	f.Build()

	s, err := response.Build()
	if err == nil {
		fmt.Print(s)
	}
	return err
}

func (t *Departures) newDepartureFrame(response *telstar.Response, result *service.StationResult) *telstar.FrameBuilder {

	stationTiploc, _ := result.Tiplocs.Get(result.Station[0])

	f := response.NewFrame().
		Header("[C]DepartureBoards.mobi[W]").
		NavMessage("[B][n][W]Data as of %s", result.Date.Format("2006 Jan 2 15:04:05")).
		Printf("[B][n][W][D]%s\r\n\n", stationTiploc.Name)
	f.Route(1, *t.PageNumber).
		Route(10, *t.PageNumber)
	return f
}
