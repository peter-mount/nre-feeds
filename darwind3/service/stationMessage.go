package service

import (
	"github.com/peter-mount/go-kernel/rest"
	"github.com/peter-mount/nre-feeds/darwind3"
	"strconv"
)

// StationMessageHandler implements the /live/message/{id} rest endpoint
func (d *DarwinD3Service) StationMessageHandler(r *rest.Rest) error {
	if id, err := strconv.ParseInt(r.Var("id"), 10, 64); err != nil {
		r.Status(404)
	} else if msg := d.darwind3.Messages.Get(int64(id)); msg != nil {
		r.Status(200).
			JSON().
			Value(msg)
	} else {
		r.Status(404)
	}

	return nil
}

// BroadcastStationMessagesHandler allows us to re-broadcast all messages
func (d *DarwinD3Service) BroadcastStationMessagesHandler(r *rest.Rest) error {
	d.darwind3.BroadcastStationMessages(nil)
	r.Status(200).
		JSON().
		Value("OK")

	return nil
}

// CrsMessageHandler Returns all messages for a CRS
func (d *DarwinD3Service) AllMessageHandler(r *rest.Rest) error {
	var messages []*darwind3.StationMessage

	d.darwind3.Messages.ForEach(func(s *darwind3.StationMessage) error {
		messages = append(messages, s)
		return nil
	})

	r.Status(200).
		JSON().
		Value(messages)

	return nil
}

// CrsMessageHandler Returns all messages for a CRS
func (d *DarwinD3Service) CrsMessageHandler(r *rest.Rest) error {
	crs := r.Var("crs")

	var messages []*darwind3.StationMessage

	d.darwind3.Messages.ForEach(func(s *darwind3.StationMessage) error {
		for _, c := range s.Station {
			if c == crs {
				messages = append(messages, s)
				break
			}
		}
		return nil
	})

	r.Status(200).
		JSON().
		Value(messages)

	return nil
}
