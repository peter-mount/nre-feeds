package service

import (
	"github.com/peter-mount/golib/rest"
	"github.com/peter-mount/nre-feeds/darwinref"
)

// viaResolveHandler resolves The via(s) for a set of schedules
func (dr *DarwinRefService) ViaResolveHandler(r *rest.Rest) error {

	// The query
	queries := make(map[string]*darwinref.ViaResolveRequest)

	// The response
	response := make(map[string]*darwinref.Via)

	// Run the queries
	if err := r.Body(&queries); err != nil {

		// Fail safe by returning 500 but still a {} object
		r.Status(500).Value(response)

	} else {

		for rid, query := range queries {
			if via := dr.reference.ResolveVia(query.Crs, query.Destination, query.Tiplocs); via != nil {
				via.SetSelf(r)
				response[rid] = via
			}
		}

		r.Status(200).
			JSON().
			Value(response)
	}

	return nil
}
