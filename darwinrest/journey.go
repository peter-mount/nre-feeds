package darwinrest

import (
	"encoding/xml"
	bolt "github.com/etcd-io/bbolt"
	"github.com/peter-mount/go-kernel/rest"
	"github.com/peter-mount/nre-feeds/darwinref"
	"github.com/peter-mount/nre-feeds/darwintimetable"
)

type result struct {
	XMLName   xml.Name                 `json:"-" xml:"result"`
	RID       string                   `json:"rid" xml:"rid,attr"`
	Journey   *darwintimetable.Journey `json:"journey" xml:"journey"`
	Locations *darwinref.LocationMap   `json:"locations" xml:"locations>LocationRef"`
	Tocs      *darwinref.TocMap        `json:"tocs" xml:"tocs>TocRef"`
	Self      string                   `json:"self" xml:"self,attr"`
}

// JourneyHandler returns a Journey from the timetable and any reference data
func (rs *DarwinRest) JourneyHandler(r *rest.Rest) error {
	res := &result{RID: r.Var("rid")}

	if err := rs.TT.View(func(tx *bolt.Tx) error {
		if journey, exists := rs.TT.GetJourney(tx, res.RID); exists {
			res.Journey = journey
		}
		return nil
	}); err != nil {
		return err
	}

	if res.Journey == nil {
		r.Status(404)
		return nil
	}

	res.Locations = darwinref.NewLocationMap()
	res.Tocs = darwinref.NewTocMap()
	if err := rs.Ref.View(func(tx *bolt.Tx) error {

		var tpls []string
		for _, l := range res.Journey.Schedule {
			tpls = append(tpls, l.Tiploc)
		}
		res.Locations.AddTiplocs(rs.Ref, tx, tpls)

		var tocs []string
		res.Locations.ForEach(func(l *darwinref.Location) {
			if l.Toc != "" {
				tocs = append(tocs, l.Toc)
			}
		})
		res.Tocs.AddTocs(rs.Ref, tx, tocs)

		return nil
	}); err != nil {
		return err
	}

	res.Self = r.Self(r.Context() + "/journey/" + res.RID)
	r.Status(200).Value(res)
	return nil
}
