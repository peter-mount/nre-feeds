package darwintty

import (
	"fmt"
	"github.com/peter-mount/nre-feeds/darwind3"
	"github.com/peter-mount/nre-feeds/ldb/service"
	"github.com/peter-mount/nre-feeds/util"
	"sort"
	"strings"
)

type Board struct {
	Crs        string       `json:"crs" xml:"crs,attr"`
	Name       string       `json:"name" xml:"name,attr"`
	Departures []*Departure `json:"departures" xml:"departures"`
}

type Departure struct {
	Destination string `json:"destination" xml:"destination,attr"`
	Via         string `json:"via,omitempty" xml:"via,omitempty"`
	Plat        string `json:"plat,omitempty" xml:"plat,attr,omitempty"`
	Depart      string `json:"depart" xml:"depart,attr"`
	depart      util.WorkingTime
	Expected    string `json:"expected" xml:"expected,attr"`
	expected    util.WorkingTime
	Delay       int                       `json:"delay,omitempty" xml:"delay,omitempty"`
	Length      int                       `json:"length,omitempty" xml:"length,attr,omitempty"`
	Cancelled   bool                      `json:"cancelled,omitempty" xml:"cancelled,attr,omitempty"`
	Reason      string                    `json:"reason" xml:"reason,omitempty"`
	Toc         string                    `json:"toc" xml:"toc,attr"`
	LastReport  Location                  `json:"lastReport" xml:"lastReport"`
	Formation   []darwind3.CoachFormation `json:"formation,omitempty" xml:"formation,omitempty"`
}

func (d *Departure) TocName() string {
	if d.Toc != "" {
		return d.Toc + " service"
	}
	return ""
}

func (d *Departure) Coaches() string {
	if d.Length > 0 {
		return fmt.Sprintf("Formed of %d coaches", d.Length)
	}
	return ""
}

func (d *Departure) FormationString() string {
	var a []string
	for _, coach := range d.Formation {
		a = append(a, coach.CoachNumber)
	}
	if len(a) == 0 {
		return ""
	}
	return "Coaches: " + strings.Join(a, ", ")
}

func (d *Departure) ToiletStatus() string {
	var t []string
	var acc []string
	for _, coach := range d.Formation {
		switch coach.Toilet.Type {
		case "Standard":
			t = append(t, coach.CoachNumber)
		case "Accessible":
			acc = append(acc, coach.CoachNumber)
		}
	}
	if len(t)+len(acc) == 0 {
		return ""
	}

	var a []string
	if len(t) > 0 {
		a = append(a, "Standard: "+plural(t))
	}
	if len(acc) > 0 {
		a = append(a, "Accessible: "+plural(acc))
	}

	return "Toilets are available. " + strings.Join(a, ", ")
}

func plural(s []string) string {
	switch len(s) {
	case 0:
		return ""
	case 1:
		return s[0]
	default:
		l := len(s)
		return strings.Join(s[:l-1], ", ") + " & " + s[l-1]
	}
}

type Location struct {
	Location string `json:"location" xml:"location,attr"`
	Time     string `json:"time" xml:"time,attr"`
	At       bool   `json:"at,omitempty" xml:"at,attr,omitempty"`
	Departed bool   `json:"departed,omitempty" xml:"departed,attr,omitempty"`
}

func (l Location) String() string {
	if l.Location == "" {
		return ""
	}

	s := ""
	switch {
	case l.Departed:
		s = " departing"
	case l.At:
		s = " at"
	}
	return fmt.Sprintf("Last seen%s %s at %s", s, l.Location, l.Time)
}

type BoardVisitor interface {
	VisitBoard(*Board) error
	VisitDeparture(*Board, *Departure) error
}

func (b *Board) SortByExpected() {
	sort.SliceStable(b.Departures, b.sortByExpected)
}

func (b *Board) sortByExpected(i, j int) bool {
	return b.Departures[i].expected.Before(&b.Departures[j].expected)
}

func (b *Board) SortByPlatform() {
	sort.SliceStable(b.Departures, b.sortByPlatform)
}

func (b *Board) sortByPlatform(i, j int) bool {
	pi := strings.ToLower(b.Departures[i].Plat)
	pj := strings.ToLower(b.Departures[j].Plat)
	if pi < pj {
		return true
	}
	if pi == pj {
		return b.sortByExpected(i, j)
	}
	return false
}

// FilterByPlatformExists filters departures to include only those with platforms set
func (b *Board) FilterByPlatformExists() {
	var a []*Departure
	for _, d := range b.Departures {
		if d.Plat != "" {
			a = append(a, d)
		}
	}
	b.Departures = a
}

// FilterByPlatform filters departures to include only the specified platform.
// Note, Platforms like 1A, 1B or 1C will be distinct to 1
func (b *Board) FilterByPlatform(plat string) {
	plat = strings.ToLower(plat)
	var a []*Departure
	for _, d := range b.Departures {
		if strings.ToLower(d.Plat) == plat {
			a = append(a, d)
		}
	}
	b.Departures = a
}

func NewBoard(result *service.StationResult) *Board {
	b := &Board{
		Crs:  result.Crs,
		Name: GetTiploc(result, result.Station[0]),
	}

	for _, departure := range result.Services {
		loc := departure.Location
		if !loc.IsDestination() && !loc.IsSetDownOnly() {
			d := &Departure{
				Destination: GetDestName(result, departure),
				depart:      loc.Time,
				Depart:      loc.Time.String()[:5],
				expected:    loc.Forecast.Time,
				Expected:    loc.Forecast.Time.String()[:5],
				Cancelled:   loc.Cancelled,
				Toc:         departure.Toc,
				Length:      loc.Length,
				Delay:       loc.Delay,
				Formation:   departure.Formation.Formation.Coaches,
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
