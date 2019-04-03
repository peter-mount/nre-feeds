package service

import (
	"github.com/peter-mount/golib/rest"
	"log"
)

func (d *DarwinKBService) StationHandler(r *rest.Rest) error {
	data, err := d.darwinkb.GetStation(r.Var("crs"))
	if err != nil {
		log.Println(err)
		return err
	}

	if data == nil {
		r.Status(404)
	} else {
		r.Status(200).
			JSON().
			Writer().
			Write(data)
	}

	return nil
}
