package service

import (
	bolt "github.com/etcd-io/bbolt"
	"github.com/peter-mount/go-kernel/v2/rest"
	"github.com/peter-mount/nre-feeds/darwinref"
	"sort"
	"strconv"
)

func (dr *DarwinRefService) ReasonCancelHandler(r *rest.Rest) error {
	return dr.reasonHandler(true, r)
}

func (dr *DarwinRefService) ReasonLateHandler(r *rest.Rest) error {
	return dr.reasonHandler(false, r)
}

func (dr *DarwinRefService) reasonHandler(cancelled bool, r *rest.Rest) error {
	id, err := strconv.Atoi(r.Var("id"))
	if err != nil {
		return err
	}

	return dr.reference.View(func(tx *bolt.Tx) error {
		var reason *darwinref.Reason
		var exists bool

		if cancelled {
			reason, exists = dr.reference.GetCancellationReason(tx, id)
		} else {
			reason, exists = dr.reference.GetLateReason(tx, id)
		}

		if exists {
			r.Status(200).
				JSON().
				Value(reason)
		} else {
			r.Status(404)
		}

		return nil
	})
}

func (dr *DarwinRefService) AllReasonCancelHandler(r *rest.Rest) error {
	return dr.allReasonHandler("DarwinCancelReason", "/reason/cancelled", r)
}

func (dr *DarwinRefService) AllReasonLateHandler(r *rest.Rest) error {
	return dr.allReasonHandler("DarwinLateReason", "/reason/late", r)
}

func (dr *DarwinRefService) allReasonHandler(bname string, prefix string, r *rest.Rest) error {
	return dr.reference.View(func(tx *bolt.Tx) error {
		resp := &darwinref.ReasonsResponse{}

		if err := tx.Bucket([]byte(bname)).ForEach(func(k, v []byte) error {
			reason := &darwinref.Reason{}
			if reason.FromBytes(v) {
				resp.Reasons = append(resp.Reasons, reason)
			}
			return nil
		}); err != nil {
			return err
		} else {

			// Sort result by code
			sort.SliceStable(resp.Reasons, func(a, b int) bool {
				ra := resp.Reasons[a]
				rb := resp.Reasons[b]
				return ra.Code < rb.Code
			})

			r.Status(200).
				JSON().
				Value(resp)
		}

		return nil
	})
}
