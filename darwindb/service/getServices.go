package service

import (
	"github.com/peter-mount/golib/rest"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (s *DarwinDBService) getStationServices(r *rest.Rest) error {
	crs := r.Var("crs")
	dateStr := r.Var("date")
	hourStr := r.Var("hour")

	log.Println(crs, dateStr, hourStr)

	crs = strings.ToUpper(crs)

	var date time.Time
	var hour int
	if dateStr == "" || hourStr == "" {
		r.Status(http.StatusNotFound)
		return nil
	} else {
		t, err := time.Parse("20060102", dateStr)
		if err != nil {
			r.Status(http.StatusNotFound)
			return nil
		}
		date = t

		hour, err = strconv.Atoi(hourStr)
		if err != nil || hour < 0 || hour > 23 {
			r.Status(http.StatusNotFound)
			return nil
		}
	}

	// Get the services
	services, err := s.db.GetServices(crs, date.Add(time.Hour*time.Duration(hour)))
	if err != nil {
		log.Println(err)
		r.Status(http.StatusInternalServerError)
		return nil
	}

	r.Status(http.StatusOK).
		JSON().
		Value(services)

	return nil
}
