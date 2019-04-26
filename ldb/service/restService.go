package service

import (
	"github.com/peter-mount/golib/rest"
	"github.com/peter-mount/nre-feeds/darwind3"
	d3client "github.com/peter-mount/nre-feeds/darwind3/client"
	"github.com/peter-mount/nre-feeds/darwinref"
	refclient "github.com/peter-mount/nre-feeds/darwinref/client"
	"time"
)

type serviceResult struct {
	// The service RID
	RID string `json:"rid"`
	// Origin
	Origin *darwind3.Location `json:"origin"`
	// Destination
	Destination *darwind3.Location `json:"destination"`
	// The service
	Service *darwind3.Schedule `json:"service"`
	// Map of Tiploc's
	Tiplocs *darwinref.LocationMap `json:"tiploc"`
	// Map of Toc's
	Tocs *darwinref.TocMap `json:"toc"`
	// Cancellation or Late Reasons
	Reasons *darwinref.ReasonMap `json:"reasons"`
	// Map of Via text by RID
	Via map[string]*darwinref.Via `json:"via"`
	// The date of this request
	Date time.Time `json:"date"`
	// The URL of this departure board
	Self string `json:"self"`
}

// serviceHandler proxies the service from d3 but fills in the required
// details of tiplocs, toc etc
func (d *LDBService) serviceHandler(r *rest.Rest) error {

	rid := r.Var("rid")

	d3Client := &d3client.DarwinD3Client{Url: d.ldb.Darwin}
	refClient := &refclient.DarwinRefClient{Url: d.ldb.Reference}

	service, err := d3Client.GetSchedule(rid)
	if err != nil {
		return err
	}

	if service == nil {
		r.Status(404)
		return nil
	}

	res := &serviceResult{
		RID:     rid,
		Service: service,
		Tiplocs: darwinref.NewLocationMap(),
		Tocs:    darwinref.NewTocMap(),
		Reasons: darwinref.NewReasonMap(),
		Date:    time.Now(),
		Self:    r.Self("/service/" + rid),
	}

	// resolve the origin & destination
	if service.Origin != nil {
		res.Origin = service.Origin
	}
	if service.Destination != nil {
		res.Destination = service.Destination
	}

	if res.Origin == nil || res.Destination == nil {
		// Just incase, if we don't have an origin/destination then use the first & last locations in the schedule
		if len(service.Locations) > 0 {
			if res.Origin == nil {
				res.Origin = service.Locations[0]
			}

			if res.Destination == nil {
				res.Destination = service.Locations[len(service.Locations)-1]
			}
		}
	}

	// Set of tiplocs
	tiplocs := make(map[string]interface{})
	service.AddTiplocs(tiplocs)

	// Toc running this service
	refClient.AddToc(res.Tocs, service.Toc)

	// Tiploc in a cancel or late reason
	if service.CancelReason.Reason > 0 {
		if reason, _ := refClient.GetCancelledReason(service.CancelReason.Reason); reason != nil {
			res.Reasons.AddReason(reason)
		}

		if service.CancelReason.Tiploc != "" {
			tiplocs[service.CancelReason.Tiploc] = nil
		}
	}

	if service.LateReason.Reason > 0 {
		if reason, _ := refClient.GetLateReason(service.LateReason.Reason); reason != nil {
			res.Reasons.AddReason(reason)
		}

		if service.LateReason.Tiploc != "" {
			tiplocs[service.LateReason.Tiploc] = nil
		}
	}

	// Resolve the via text. For the service this is for the origin only
	vias := make(map[string]*darwinref.ViaResolveRequest)

	if len(service.Locations) > 2 {
		// We need the crs of the origin from the resolved tiploc map
		loc, _ := refClient.GetTiploc(service.Origin.Tiploc)
		if loc.Crs != "" {
			viaRequest := &darwinref.ViaResolveRequest{
				Crs:         loc.Crs,
				Destination: service.Locations[len(service.Locations)-1].Tiploc,
			}
			for _, loc := range service.Locations[1:] {
				viaRequest.Tiplocs = append(viaRequest.Tiplocs, loc.Tiploc)
			}
			vias[rid] = viaRequest
		}
	}

	// The association tiplocs
	for _, assoc := range service.Associations {
		assoc.AddTiplocs(tiplocs)

		if assoc.IsJoin() || assoc.IsSplit() {
			ar := assoc.Main.RID
			ai := assoc.Main.LocInd
			if ar == service.RID {
				ar = assoc.Assoc.RID
				ai = assoc.Assoc.LocInd
			}
			if ar != service.RID {
				as := d.ldb.GetSchedule(ar)
				if as != nil {
					assoc.Schedule = as
					refClient.AddToc(res.Tocs, as.Toc)

					if ai < (len(as.Locations) - 1) {
						loc, _ := refClient.GetTiploc(as.Origin.Tiploc)
						if loc.Crs != "" {
							viaRequest := &darwinref.ViaResolveRequest{
								Crs:         loc.Crs,
								Destination: as.Locations[len(as.Locations)-1].Tiploc,
							}

							for _, l := range as.Locations[ai:] {
								tiplocs[l.Tiploc] = nil
								viaRequest.Tiplocs = append(viaRequest.Tiplocs, l.Tiploc)
							}

							vias[ar] = viaRequest
						}
					}

					// Cancellation reason
					if as.CancelReason.Reason > 0 {
						if reason, _ := refClient.GetCancelledReason(as.CancelReason.Reason); reason != nil {
							res.Reasons.AddReason(reason)
						}

						if as.CancelReason.Tiploc != "" {
							tiplocs[as.CancelReason.Tiploc] = nil
						}
					}

					// Late reason
					if as.LateReason.Reason > 0 {
						if reason, _ := refClient.GetLateReason(as.LateReason.Reason); reason != nil {
							res.Reasons.AddReason(reason)
						}

						if as.LateReason.Tiploc != "" {
							tiplocs[as.LateReason.Tiploc] = nil
						}
					}

				}
			}
		}
	}

	if len(vias) > 0 {
		resolvedVias, _ := refClient.GetVias(vias)
		if resolvedVias != nil {
			res.Via = resolvedVias
			for _, v := range res.Via {
				tiplocs[v.Dest] = nil
				if v.Loc1 != "" {
					tiplocs[v.Loc1] = nil
				}
				if v.Loc2 != "" {
					tiplocs[v.Loc2] = nil
				}
			}
		}
	}

	// Now resolve the tiplocs en-masse and resolve the toc's at the same time
	if locs, _ := refClient.GetTiplocsMapKeys(tiplocs); locs != nil {
		res.Tiplocs.AddAll(locs)

		for _, l := range locs {
			refClient.AddToc(res.Tocs, l.Toc)
		}
	}

	r.Status(200).
		JSON().
		Value(res)

	return nil
}
