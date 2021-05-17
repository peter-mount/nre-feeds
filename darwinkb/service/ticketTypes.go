package service

import (
	"github.com/peter-mount/go-kernel/rest"
	"log"
)

func (d *DarwinKBService) TicketTypesHandler(r *rest.Rest) error {
	data, err := d.darwinkb.GetTicketTypes()
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

func (d *DarwinKBService) TicketIdsHandler(r *rest.Rest) error {
	data, err := d.darwinkb.GetTicketIDs()
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

func (d *DarwinKBService) TicketTypeHandler(r *rest.Rest) error {
	data, err := d.darwinkb.GetTicketType(r.Var("id"))
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
