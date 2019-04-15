package service

import (
	"github.com/peter-mount/golib/rest"
	"github.com/peter-mount/nre-feeds/darwind3"
	d3client "github.com/peter-mount/nre-feeds/darwind3/client"
	"github.com/peter-mount/nre-feeds/darwinref"
	refclient "github.com/peter-mount/nre-feeds/darwinref/client"
	"github.com/peter-mount/nre-feeds/ldb"
	"github.com/peter-mount/nre-feeds/util"
	"sort"
	"time"
)

type stationResult struct {
	// The requested crs
	Crs string `json:"crs"`
	// The departures
	Services []ldb.Service `json:"departures"`
	// Details about this station
	Station []string `json:"station"`
	// Map of Tiploc's
	Tiplocs *darwinref.LocationMap `json:"tiploc"`
	// Map of Toc's
	Tocs *darwinref.TocMap `json:"toc"`
	// StationMessages
	Messages []*darwind3.StationMessage `json:"messages"`
	// Cancellation or Late Reasons
	Reasons *darwinref.ReasonMap `json:"reasons"`
	// Map of Via text by RID
	Via map[string]*darwinref.Via `json:"via"`
	// The date of this request
	Date time.Time `json:"date"`
	// The URL of this departure board
	Self string `json:"self"`
}

func (d *LDBService) stationHandler(r *rest.Rest) error {

	crs := r.Var("crs")

	station := d.ldb.GetStationCrs(crs)

	if station == nil {
		r.Status(404)
	} else {

		d3Client := &d3client.DarwinD3Client{Url: d.ldb.Darwin}
		refClient := &refclient.DarwinRefClient{Url: d.ldb.Reference}

		// We want everything for the next hour
		location, _ := time.LoadLocation("Europe/London")
		now := time.Now().In(location)
		from := util.WorkingTime_FromTime(now)
		to := util.WorkingTime_FromTime(now.Add(time.Hour))

		services := d.ldb.GetServices(station, from, to)

		res := &stationResult{
			Crs: crs,
			//Services: services,
			Tiplocs:  darwinref.NewLocationMap(),
			Tocs:     darwinref.NewTocMap(),
			Messages: station.GetMessages(d3Client),
			Reasons:  darwinref.NewReasonMap(),
			Date:     now,
			Self:     r.Self("/boards/" + crs),
		}

		b := stationBuilder{
			// Set of tiplocs
			tiplocs: make(map[string]interface{}),

			// Map of via texts
			vias: make(map[string]*darwinref.ViaResolveRequest),
		}

		// Station details
		if sl, _ := refClient.GetCrs(crs); sl != nil {
			for _, l := range sl.Tiploc {
				res.Station = append(res.Station, l.Tiploc)
				b.tiplocs[l.Tiploc] = nil
			}
		}

		// Tiplocs within the departures
		for _, se := range services {
			s := d.getService(&b, se)
			if s.RID != "" {
				res.Services = append(res.Services, s)

				if len(s.CallingPoints) > 0 {
					viaRequest := &darwinref.ViaResolveRequest{
						Crs:         station.Crs,
						Destination: s.CallingPoints[len(s.CallingPoints)-1].Tiploc,
					}
					b.vias[s.RID] = viaRequest

					for _, cp := range s.CallingPoints {
						b.tiplocs[cp.Tiploc] = nil
						viaRequest.Tiplocs = append(viaRequest.Tiplocs, cp.Tiploc)
					}
				}

				// The association tiplocs
				for _, assoc := range s.Associations {
					b.tiplocs[assoc.Tiploc] = nil
					if assoc.IsJoin() || assoc.IsSplit() {
						ar := assoc.Main.RID
						ai := assoc.Main.LocInd
						if ar == s.RID {
							ar = assoc.Assoc.RID
							ai = assoc.Assoc.LocInd
						}
						if ar != s.RID {
							as := d.ldb.GetSchedule(ar)
							if as != nil {
								assoc.Schedule = as
								refClient.AddToc(res.Tocs, as.Toc)

								if as.Origin != nil {
									b.tiplocs[as.Origin.Tiploc] = nil
								}

								if as.Destination != nil {
									b.tiplocs[as.Destination.Tiploc] = nil
								}

								if ai < (len(as.Locations) - 1) {
									viaRequest := &darwinref.ViaResolveRequest{
										Crs:         station.Crs,
										Destination: as.Locations[len(as.Locations)-1].Tiploc,
									}
									b.vias[ar] = viaRequest

									for _, l := range as.Locations[ai:] {
										b.tiplocs[l.Tiploc] = nil
										viaRequest.Tiplocs = append(viaRequest.Tiplocs, l.Tiploc)
									}
								}

								// Cancellation reason
								if as.CancelReason.Reason > 0 {
									if reason, _ := refClient.GetCancelledReason(as.CancelReason.Reason); reason != nil {
										res.Reasons.AddReason(reason)
									}

									if as.CancelReason.Tiploc != "" {
										b.tiplocs[as.CancelReason.Tiploc] = nil
									}
								}

								// Late reason
								if as.LateReason.Reason > 0 {
									if reason, _ := refClient.GetLateReason(as.LateReason.Reason); reason != nil {
										res.Reasons.AddReason(reason)
									}

									if as.LateReason.Tiploc != "" {
										b.tiplocs[as.LateReason.Tiploc] = nil
									}
								}

							}
						}
					}
				}

				// Toc running this service
				refClient.AddToc(res.Tocs, s.Toc)

				// Cancellation reason
				if s.CancelReason.Reason > 0 {
					if reason, _ := refClient.GetCancelledReason(s.CancelReason.Reason); reason != nil {
						res.Reasons.AddReason(reason)
					}

					if s.CancelReason.Tiploc != "" {
						b.tiplocs[s.CancelReason.Tiploc] = nil
					}
				}

				// Late reason
				if s.LateReason.Reason > 0 {
					if reason, _ := refClient.GetLateReason(s.LateReason.Reason); reason != nil {
						res.Reasons.AddReason(reason)
					}

					if s.LateReason.Tiploc != "" {
						b.tiplocs[s.LateReason.Tiploc] = nil
					}
				}

				// Set self to point to our service endpoint
				s.Self = r.Self("/service/" + s.RID)
			}
		}

		// Now resolve the tiplocs en-masse and resolve the toc's at the same time
		if locs, _ := refClient.GetTiplocsMapKeys(b.tiplocs); locs != nil {
			res.Tiplocs.AddAll(locs)

			for _, l := range locs {
				refClient.AddToc(res.Tocs, l.Toc)
			}
		}

		// Resolve via texts
		if len(b.vias) > 0 {
			if vias, _ := refClient.GetVias(b.vias); vias != nil {
				res.Via = vias
			}
		}

		// sort into time order
		sort.SliceStable(res.Services, func(i, j int) bool {
			return res.Services[i].Location.Compare(&res.Services[j].Location)
		})

		r.Status(200).
			JSON().
			Value(res)
	}

	return nil
}

type stationBuilder struct {
	tiplocs map[string]interface{}
	vias    map[string]*darwinref.ViaResolveRequest
}

func (d *LDBService) getService(b *stationBuilder, se ldb.ServiceEntry) ldb.Service {
	sched := d.ldb.GetSchedule(se.RID)
	if sched == nil {
		return ldb.Service{}
	}

	s := ldb.Service{}
	if !s.Update(sched, se.LocationIndex) {
		return ldb.Service{}
	}

	// Destination & location tiplocs
	b.tiplocs[s.Destination] = nil
	b.tiplocs[s.Location.Tiploc] = nil

	// The origin Location
	if s.Origin.Tiploc != "" {
		b.tiplocs[s.Origin.Tiploc] = nil
	}

	// The destination Location
	if s.Dest.Tiploc != "" {
		b.tiplocs[s.Dest.Tiploc] = nil
	}

	// Add CallingPoints tiplocs to map & via request
	s.Associations = sched.Associations

	s.CallingPoints = sched.GetCallingPoints(s.LocationIndex)

	s.LastReport = sched.GetLastReport()
	if s.LastReport.Tiploc != "" {
		b.tiplocs[s.LastReport.Tiploc] = nil
	}

	return s
}
