package service

import "github.com/peter-mount/go-kernel/rest"

type badgeResponse struct {
	SchemaVersion int    `json:"schemaVersion"`
	Label         string `json:"label"`
	Message       string `json:"message"`
	Colour        string `json:"color"`
}

func (d *DarwinD3Service) StatusHandler(r *rest.Rest) error {
	status, colour := d.darwind3.GetStatus()
	resp := badgeResponse{
		1,
		"DarwinV16",
		status,
		colour,
	}
	r.Status(200).
		JSON().
		Value(resp)
	return nil
}
