package service

import (
	bolt "github.com/etcd-io/bbolt"
	"github.com/peter-mount/golib/rest"
	"github.com/peter-mount/nre-feeds/darwinref"
)

func (dr *DarwinRefService) TiplocHandler(r *rest.Rest) error {
	return dr.reference.View(func(tx *bolt.Tx) error {
		tpl := r.Var("id")

		if location, exists := dr.reference.GetTiploc(tx, tpl); exists {
			r.Status(200).
				JSON().
				Value(location)
		} else {
			r.Status(404)
		}

		return nil
	})
}

func (dr *DarwinRefService) TiplocsHandler(r *rest.Rest) error {

	tiplocs := make([]string, 0)

	if err := r.Body(&tiplocs); err != nil {
		return err
	}

	return dr.reference.View(func(tx *bolt.Tx) error {
		var ary []*darwinref.Location

		for _, tpl := range tiplocs {
			if location, exists := dr.reference.GetTiploc(tx, tpl); exists {
				ary = append(ary, location)
			}
		}

		r.Status(200).
			JSON().
			Value(ary)

		return nil
	})
}
