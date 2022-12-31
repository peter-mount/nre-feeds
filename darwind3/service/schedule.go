package service

import (
	"github.com/peter-mount/go-kernel/v2/rest"
)

func (d *DarwinD3Service) ScheduleHandler(r *rest.Rest) error {
	if sched := d.darwind3.GetSchedule(r.Var("rid")); sched != nil {
		d.darwind3.UpdateAssociations(sched)

		r.Status(200).
			JSON().
			Value(sched)
	} else {
		r.Status(404)
	}

	return nil
}
