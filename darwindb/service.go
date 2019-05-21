package darwindb

import (
	"database/sql"
	"encoding/json"
	"github.com/peter-mount/nre-feeds/darwind3"
)

// The json structure returned by darwin.getservice()
type ServiceDetail struct {
	Rid       int               `json:"rid"`
	Schedule  darwind3.Schedule `json:"schedule"`
	Tiploc    map[string]Tiploc `json:"tiploc"`
	Formation []Coach           `json:"formation"`
}

// Used to aggregate the formation coaches
type Coach struct {
	Number     string
	Class      string
	ToiletType string
}

func (s *DarwinDB) GetService(rid string) (ServiceDetail, error) {
	var service ServiceDetail
	var data sql.NullString

	err := s.getServiceStatement.QueryRow(rid).Scan(&data)

	if err == nil && data.Valid {
		err = json.Unmarshal([]byte(data.String), &service)
	}

	if err == nil {

		// Resolve the formation data
		for _, c := range service.Schedule.Formation.Formation.Coaches {
			service.Formation = append(service.Formation, Coach{
				Number:     c.CoachNumber,
				Class:      c.CoachClass,
				ToiletType: c.Toilet.Type,
			})
		}

	}

	return service, err
}
