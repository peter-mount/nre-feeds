package service

import (
	"github.com/peter-mount/golib/rest"
	"log"
	"net/http"
	"strings"
)

func (s *DarwinDBService) getService(r *rest.Rest) error {
	rid := r.Var("rid")

	log.Println(rid)

	service, err := s.db.GetService(rid)

	if err != nil || service.Schedule.RID != rid {
		log.Println(rid, err)
		return err
	}

	// Fix Activity so it's space delimited
	for _, l := range service.Schedule.Locations {
		if l.Planned.ActivityType != "" {
			var a []string
			for i := 0; i < len(l.Planned.ActivityType); i += 2 {
				a = append(a, l.Planned.ActivityType[i:i+2])
			}
			l.Planned.ActivityType = strings.Join(a, " ")
		}
	}

	r.Status(http.StatusOK).
		JSON().
		Value(service)

	return nil
}
