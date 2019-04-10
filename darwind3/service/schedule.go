package service

import (
	"github.com/peter-mount/golib/rest"
)

func (d *DarwinD3Service) ScheduleHandler(r *rest.Rest) error {
	if sched := d.darwind3.GetSchedule(r.Var("rid")); sched != nil {
		sched.SetSelf(r)

		d.darwind3.UpdateAssociations(sched)

		r.Status(200).
			JSON().
			Value(sched)
	} else {
		r.Status(404)
	}

	return nil
}
