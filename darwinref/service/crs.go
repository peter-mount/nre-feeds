package service

import (
	"encoding/json"
	bolt "github.com/etcd-io/bbolt"
	"github.com/peter-mount/go-kernel/rest"
	"github.com/peter-mount/nre-feeds/darwinref"
)

func (dr *DarwinRefService) CrsHandler(r *rest.Rest) error {
	return dr.reference.View(func(tx *bolt.Tx) error {
		crs := r.Var("id")

		if locations, exists := dr.reference.GetCrs(tx, crs); exists {
			resp := &darwinref.CrsResponse{}
			r.Status(200).
				JSON().
				Value(resp)

			resp.Crs = crs
			resp.Tiploc = locations
		} else {
			r.Status(404)
		}

		return nil
	})
}

func (dr *DarwinRefService) AllCrsHandler(r *rest.Rest) error {
	var t []*darwinref.Location

	if err := dr.reference.View(func(tx *bolt.Tx) error {
		crsBucket := tx.Bucket([]byte("DarwinCrs"))
		tiplocBucket := tx.Bucket([]byte("DarwinTiploc"))

		return crsBucket.ForEach(func(k, v []byte) error {
			var tpls []string

			err := json.Unmarshal(v, &tpls)
			if err != nil {
				return err
			}

			for _, tpl := range tpls {
				if loc, exists := dr.reference.GetTiplocBucket(tiplocBucket, tpl); exists {
					t = append(t, loc)
				}
			}

			r.Status(200).
				JSON().
				Value(t)

			return nil
		})
	}); err != nil {
		r.Status(500).
			JSON().
			Value(err)
	}

	r.Status(200).
		JSON().
		Value(t)
	return nil
}
