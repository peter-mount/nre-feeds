package darwindb

import (
	"database/sql"
	"encoding/json"
	"github.com/peter-mount/nre-feeds/darwind3"
	"github.com/peter-mount/nrod-cif/cif"
	"time"
)

type StationServices struct {
	Station  []string              `json:"station"`
	Services []StationService      `json:"services"`
	Reason   map[string]Reason     `json:"reason"`
	Tiploc   map[string]cif.Tiploc `json:"tiploc"`
}

type StationService struct {
	Rid              int64                      `json:"rid"`
	Location         darwind3.Location          `json:"location"`
	Destination      string                     `json:"destination"`
	Cancelled        bool                       `json:"cancelled"`
	CancelReason     darwind3.DisruptionReason  `json:"cancelReason"`
	DelayReason      darwind3.DisruptionReason  `json:"delayReason"`
	Uid              string                     `json:"uid"`
	Status           string                     `json:"status"`
	Headcode         string                     `json:"trainId"`
	PassengerService bool                       `json:"passengerService"`
	CharterService   bool                       `json:"charterService"`
	Toc              string                     `json:"toc"`
	Association      []darwind3.Association     `json:"association"`
	Formation        darwind3.ScheduleFormation `json:"formation"`
	Delay            time.Duration              `json:"delay"`
}

// Reason text lookup as returned by the proc
type Reason struct {
	Id     int    `json:"id"`
	Cancel string `json:"cancel"`
	Late   string `json:"late"`
}

func (s *DarwinDB) GetServices(crs string, ts time.Time) (StationServices, error) {

	var services StationServices
	var data sql.NullString

	err := s.statements[getStationServicesStatement].QueryRow(crs, ts).Scan(&data)

	if err == nil && data.Valid {
		err = json.Unmarshal([]byte(data.String), &services)

		if err == nil {
			for _, sve := range services.Services {
				// Calculate the delay
				sve.Delay = time.Duration(sve.Location.Delay) * time.Second
			}
		}
	}

	return services, err
}
