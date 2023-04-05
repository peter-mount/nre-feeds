package service

import "github.com/peter-mount/go-kernel/v2/rest"

func (d *DarwinD3Service) AlarmHandler(r *rest.Rest) error {
	r.Status(200).
		JSON().
		Value(d.darwind3.GetAlarm(r.Var("id")))

	return nil
}

func (d *DarwinD3Service) AlarmsHandler(r *rest.Rest) error {
	r.Status(200).
		JSON().
		Value(d.darwind3.GetAlarms())

	return nil
}
