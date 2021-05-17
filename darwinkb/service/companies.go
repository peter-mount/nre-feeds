package service

import (
	"github.com/peter-mount/go-kernel/rest"
	"log"
)

func (d *DarwinKBService) CompaniesHandler(r *rest.Rest) error {
	data, err := d.darwinkb.GetCompanies()
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

func (d *DarwinKBService) CompanyHandler(r *rest.Rest) error {
	data, err := d.darwinkb.GetCompany(r.Var("id"))
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
