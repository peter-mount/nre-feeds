package service

import (
	bolt "github.com/etcd-io/bbolt"
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/darwinref"
)

func (dr *DarwinRefService) TocHandler(r *rest.Rest) error {
	return dr.reference.View(func(tx *bolt.Tx) error {
		id := r.Var("id")

		if toc, exists := dr.reference.GetToc(tx, id); exists {
			r.Status(200).
				JSON().
				Value(toc)
		} else {
			r.Status(404)
		}

		return nil
	})
}

func (dr *DarwinRefService) AllTocsHandler(r *rest.Rest) error {
	return dr.reference.View(func(tx *bolt.Tx) error {
		resp := &darwinref.TocsResponse{}

		if err := tx.Bucket([]byte("DarwinToc")).ForEach(func(k, v []byte) error {
			toc := &darwinref.Toc{}
			if toc.FromBytes(v) {
				resp.Toc = append(resp.Toc, toc)
			}
			return nil
		}); err != nil {
			return err
		} else {
			r.Status(200).
				JSON().
				Value(resp)
		}

		return nil
	})
}
