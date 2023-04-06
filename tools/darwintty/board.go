package darwintty

import (
	"github.com/peter-mount/nre-feeds/ldb/service"
)

type Board struct {
	Crs        string      `json:"crs" xml:"crs,attr"`
	Name       string      `json:"name" xml:"name,attr"`
	Departures []Departure `json:"departures" xml:"departures"`
}

type Departure struct {
	Destination string   `json:"destination" xml:"destination,attr"`
	Via         string   `json:"via,omitempty" xml:"via,omitempty"`
	Plat        string   `json:"plat,omitempty" xml:"plat,attr,omitempty"`
	Depart      string   `json:"depart" xml:"depart,attr"`
	Expected    string   `json:"expected" xml:"expected,attr"`
	Delay       int      `json:"delay,omitempty" xml:"delay,omitempty"`
	Length      int      `json:"length,omitempty" xml:"length,attr,omitempty"`
	Cancelled   bool     `json:"cancelled,omitempty" xml:"cancelled,attr,omitempty"`
	Reason      string   `json:"reason" xml:"reason,omitempty"`
	Toc         string   `json:"toc" xml:"toc,attr"`
	LastReport  Location `json:"lastReport" xml:"lastReport"`
}

type Location struct {
	Location string `json:"location" xml:"location,attr"`
	Time     string `json:"time" xml:"time,attr"`
	At       bool   `json:"at,omitempty" xml:"at,attr,omitempty"`
	Departed bool   `json:"departed,omitempty" xml:"departed,attr,omitempty"`
}

func NewBoard(result *service.StationResult) *Board {
	b := &Board{
		Crs:  result.Crs,
		Name: GetTiploc(result, result.Station[0]),
	}

	for _, departure := range result.Services {
		loc := departure.Location
		if !loc.IsDestination() && !loc.IsSetDownOnly() {
			d := Departure{
				Destination: GetDestName(result, departure),
				Depart:      loc.Time.String()[:5],
				Expected:    loc.Forecast.Time.String()[:5],
				Cancelled:   loc.Cancelled,
				Toc:         departure.Toc,
				Length:      loc.Length,
				Delay:       loc.Delay,
				LastReport: Location{
					Location: GetTiploc(result, departure.LastReport.Tiploc),
					Time:     departure.LastReport.Time.String()[:5],
					At:       departure.LastReport.At,
					Departed: departure.LastReport.Departed,
				},
			}

			if result.Reasons != nil {
				reason := result.Reasons.Late[departure.LateReason.Reason]
				if loc.Cancelled && departure.CancelReason.Reason != 0 {
					reason = result.Reasons.Cancelled[departure.CancelReason.Reason]
				}
				if reason != nil {
					d.Reason = reason.Text
				}
			}

			if !(loc.Forecast.Platform.CISSuppressed || loc.Forecast.Platform.Suppressed) {
				d.Plat = loc.Forecast.Platform.Platform
			}

			if tocName, _ := result.Tocs.Get(d.Toc); tocName != nil {
				d.Toc = tocName.Name
			}

			b.Departures = append(b.Departures, d)
		}
	}

	return b
}
