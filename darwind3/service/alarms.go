package service

import "github.com/peter-mount/go-kernel/rest"

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
