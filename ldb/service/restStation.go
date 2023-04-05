package service

import (
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/darwind3"
	d3client "github.com/peter-mount/nre-feeds/darwind3/client"
	"github.com/peter-mount/nre-feeds/darwinref"
	refclient "github.com/peter-mount/nre-feeds/darwinref/client"
	"github.com/peter-mount/nre-feeds/ldb"
	"github.com/peter-mount/nre-feeds/util"
	"sort"
	"strconv"
	"time"
)

type StationResult struct {
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
}

// boardFilter handles an individual request
type boardFilter struct {
	d             *LDBService
	length        int                                     // limit to this number of services if >0
	terminated    bool                                    // if true then don't include services terminating at this location
	callAt        bool                                    // If present then filter services that only arrive at a specific station
	callAtTiplocs []string                                // tiplocs to filter when callAt is true
	d3Client      *d3client.DarwinD3Client                // client to d3 service
	refClient     *refclient.DarwinRefClient              // client to ref service
	station       *ldb.Station                            // The station details
	services      []ldb.ServiceEntry                      // The available services
	res           *StationResult                          // The final result
	tiplocs       map[string]interface{}                  // tiplocs in the response
	vias          map[string]*darwinref.ViaResolveRequest // vias
	now           time.Time                               // The time the request was made
}

func (d *LDBService) createBoardFilter(r *rest.Rest, crs string, station *ldb.Station) *boardFilter {
	d3Client := &d3client.DarwinD3Client{Url: d.ldb.Darwin}
	refClient := &refclient.DarwinRefClient{Url: d.ldb.Reference}

	location, _ := time.LoadLocation("Europe/London")
	now := time.Now().In(location)

	bf := &boardFilter{
		d:         d,
		d3Client:  d3Client,
		refClient: refClient,
		now:       now,
		station:   station,
		tiplocs:   make(map[string]interface{}),
		vias:      make(map[string]*darwinref.ViaResolveRequest),
		res: &StationResult{
			Crs:      crs,
			Tiplocs:  darwinref.NewLocationMap(),
			Tocs:     darwinref.NewTocMap(),
			Messages: station.GetMessages(d3Client),
			Reasons:  darwinref.NewReasonMap(),
			Date:     now,
		},
	}

	// Parse the query parameters. Don't use gorillamux to do this as they are optional
	url := r.Request().URL
	if url != nil {

		for k, v := range url.Query() {
			switch k {
			// Limit number of returned services - the hard limit will always be in place
			case "len":
				for _, s := range v {
					l, err := strconv.Atoi(s)
					if err == nil && l >= 0 {
						bf.length = l
					}
				}

				// Filter out terminating services
			case "term":
				bf.terminated = false
				if v != nil {
					for _, e := range v {
						bf.terminated = e == "false"
					}
				}

				// Filter out services which do not stop at a station
			case "callsAt":
				for _, dest := range v {
					destination := d.ldb.GetStationCrs(dest)
					if destination != nil {
						for _, loc := range destination.Locations {
							bf.callAtTiplocs = append(bf.callAtTiplocs, loc.Tiploc)
						}
					}
					bf.callAt = len(bf.callAtTiplocs) > 0
				}

			default:
				// Ignore unknown parameters
			}
		}

	}

	return bf
}

// Is the tiploc one for this station
func (bf *boardFilter) atStation(tpl string) bool {
	for _, s := range bf.res.Station {
		if s == tpl {
			return true
		}
	}
	return false
}

// Does the service call at a specific station
func (bf *boardFilter) callsAt(callingPoints []darwind3.CallingPoint, tpls []string) bool {
	for _, cp := range callingPoints {
		for _, tpl := range tpls {
			if tpl == cp.Tiploc {
				return true
			}
		}
	}

	return false
}

// Add a tiploc to the result so that it will be included in the tiploc map
func (bf *boardFilter) addTiploc(tiploc string) {
	if tiploc != "" {
		bf.tiplocs[tiploc] = nil
	}
}

// Add a ViaResolveRequest to the response
func (bf *boardFilter) addVia(rid, dest string) *darwinref.ViaResolveRequest {
	viaRequest := &darwinref.ViaResolveRequest{
		Crs:         bf.station.Crs,
		Destination: dest,
	}
	bf.vias[rid] = viaRequest
	return viaRequest
}

// Process calling points so that we generate the appropriate via and include their tiplocs
func (bf *boardFilter) processCallingPoints(s ldb.Service) {
	if len(s.CallingPoints) > 0 {
		viaRequest := bf.addVia(s.RID, s.CallingPoints[len(s.CallingPoints)-1].Tiploc)

		for _, cp := range s.CallingPoints {
			bf.addTiploc(cp.Tiploc)
			viaRequest.AppendTiploc(cp.Tiploc)
		}
	}
}

// Process any associations, pulling in their schedules
func (bf *boardFilter) processAssociations(s ldb.Service) {
	for _, assoc := range s.Associations {
		assoc.AddTiplocs(bf.tiplocs)

		//if assoc.IsJoin() || assoc.IsSplit() {
		ar := assoc.Main.RID
		ai := assoc.Main.LocInd
		if ar == s.RID {
			ar = assoc.Assoc.RID
			ai = assoc.Assoc.LocInd
		}

		// Resolve the schedule if a split, join or if NP only if previous service & we are not yet running
		//if ar != s.RID {
		if assoc.Category != "NP" || (s.LastReport.Tiploc == "" && assoc.Assoc.RID == s.RID) {
			as := bf.d.ldb.GetSchedule(ar)
			if as != nil {
				assoc.Schedule = as
				as.AddTiplocs(bf.tiplocs)

				as.LastReport = as.GetLastReport()

				bf.processToc(as.Toc)

				if ai < (len(as.Locations) - 1) {
					if as.Origin != nil {
						bf.addTiploc(as.Destination.Tiploc)
					}

					destination := as.Locations[len(as.Locations)-1].Tiploc
					if as.Destination != nil {
						destination = as.Destination.Tiploc
					}
					viaRequest := bf.addVia(ar, destination)

					for _, l := range as.Locations[ai:] {
						bf.addTiploc(l.Tiploc)
						viaRequest.AppendTiploc(l.Tiploc)
					}
				}

				bf.processReason(as.CancelReason, true)
				bf.processReason(as.LateReason, false)

			}
		}

	}
}

func (bf *boardFilter) processToc(toc string) {
	bf.refClient.AddToc(bf.res.Tocs, toc)
}

func (bf *boardFilter) processReason(r darwind3.DisruptionReason, cancelled bool) {
	if r.Reason > 0 {
		var reason *darwinref.Reason

		if cancelled {
			reason, _ = bf.refClient.GetCancelledReason(r.Reason)
		} else {
			reason, _ = bf.refClient.GetLateReason(r.Reason)
		}

		if reason != nil {
			bf.res.Reasons.AddReason(reason)
		}

		bf.addTiploc(r.Tiploc)
	}
}

func (bf *boardFilter) resolve() {
	// Now resolve the tiplocs en-masse and resolve the toc's at the same time
	if locs, _ := bf.refClient.GetTiplocsMapKeys(bf.tiplocs); locs != nil {
		bf.res.Tiplocs.AddAll(locs)

		for _, l := range locs {
			bf.refClient.AddToc(bf.res.Tocs, l.Toc)
		}
	}

	// Resolve via texts
	if len(bf.vias) > 0 {
		if vias, _ := bf.refClient.GetVias(bf.vias); vias != nil {
			bf.res.Via = vias
		}
	}

	// sort into time order
	sort.SliceStable(bf.res.Services, func(i, j int) bool {
		return bf.res.Services[i].Location.Times.Compare(&bf.res.Services[j].Location.Times)
	})

	// Truncate
	if bf.length > 0 && len(bf.res.Services) > bf.length {
		bf.res.Services = bf.res.Services[:bf.length]
	}
}

// acceptService returns true if the service is to be accepted, false if it's to be ignored
func (bf *boardFilter) acceptService(service ldb.Service) bool {
	// Original requirement, must have an RID
	if service.RID == "" {
		return false
	}

	// remove terminating services
	if bf.terminated && bf.atStation(service.Destination) {
		return false
	}

	if bf.callAt && !bf.callsAt(service.CallingPoints, bf.callAtTiplocs) {
		return false
	}

	return true
}

func (bf *boardFilter) getService(se ldb.ServiceEntry) ldb.Service {
	sched := bf.d.ldb.GetSchedule(se.RID)
	if sched == nil {
		return ldb.Service{}
	}

	sched.UpdateTime()

	s := ldb.Service{}
	if !s.Update(sched, se.LocationIndex) {
		return ldb.Service{}
	}

	s.Associations = sched.Associations

	s.CallingPoints = sched.GetCallingPoints(s.LocationIndex)

	s.LastReport = sched.GetLastReport()

	bf.addTiploc(s.Dest.Tiploc)
	bf.addTiploc(s.Destination)
	bf.addTiploc(s.LastReport.Tiploc)
	bf.addTiploc(s.Location.Tiploc)
	bf.addTiploc(s.Origin.Tiploc)
	bf.addTiploc(s.Terminates.Tiploc)

	return s
}

func (d *LDBService) stationHandler(r *rest.Rest) error {

	crs := r.Var("crs")

	station := d.ldb.GetStationCrs(crs)

	if station == nil {
		r.Status(404)
	} else {
		filter := d.createBoardFilter(r, crs, station)

		// The services
		from := util.WorkingTime_FromTime(filter.now)
		to := util.WorkingTime_FromTime(filter.now.Add(time.Hour))
		filter.services = d.ldb.GetServices(station, from, to)

		// Station details
		if sl, _ := filter.refClient.GetCrs(crs); sl != nil {
			for _, l := range sl.Tiploc {
				filter.res.Station = append(filter.res.Station, l.Tiploc)
				filter.addTiploc(l.Tiploc)
			}
		}

		for _, se := range filter.services {
			// Resolve the service
			s := filter.getService(se)

			if filter.acceptService(s) {
				filter.res.Services = append(filter.res.Services, s)

				filter.processCallingPoints(s)
				filter.processAssociations(s)
				filter.processToc(s.Toc)
				filter.processReason(s.CancelReason, true)
				filter.processReason(s.LateReason, false)
			}
		}

		filter.resolve()

		r.Status(200).
			JSON().
			Value(filter.res)
	}

	return nil
}
