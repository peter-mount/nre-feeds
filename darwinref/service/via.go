package service

import (
	bolt "github.com/etcd-io/bbolt"
	"github.com/peter-mount/go-kernel/rest"
	"log"
)

// viaHandler returns the unique instance of a via entry
func (dr *DarwinRefService) ViaHandler(r *rest.Rest) error {
	return dr.reference.View(func(tx *bolt.Tx) error {
		log.Printf("via '%s' '%s' '%s' '%s'", r.Var("at"), r.Var("dest"), r.Var("loc1"), r.Var("loc2"))

		if via, exists := dr.reference.GetVia(tx, r.Var("at"), r.Var("dest"), r.Var("loc1"), r.Var("loc2")); exists {
			r.Status(200).
				JSON().
				Value(via)
		} else {
			r.Status(404)
		}

		return nil
	})
}
