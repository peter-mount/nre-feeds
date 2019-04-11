package ldb

import (
	"encoding/json"
	"github.com/etcd-io/bbolt"
	"github.com/peter-mount/nre-feeds/darwind3"
	d3client "github.com/peter-mount/nre-feeds/darwind3/client"
	"github.com/peter-mount/nre-feeds/darwinref"
	"github.com/peter-mount/nre-feeds/util"
	"sort"
)

// The holder for a station's departure boards
type Station struct {
	// The location details for this station
	Locations []*darwinref.Location
	Crs       string
	// This station is Public - i.e. has a CRS so can have departures
	Public bool
	// The Station message id's applicable to this station
	Messages []int64
}

func (l *LDB) updateStation(s *Station) {
	l.EventManager.PostEvent(&darwind3.DarwinEvent{
		Type: darwind3.Event_BoardUpdate,
		Crs:  s.Crs,
	})
}

// Bytes returns the message as an encoded byte slice
func (s *Station) Bytes() ([]byte, error) {
	b, err := json.Marshal(s)
	return b, err
}

// ScheduleFromBytes returns a schedule based on a slice or nil if none
func StationFromBytes(b []byte) *Station {
	if b == nil {
		return nil
	}

	station := &Station{}
	err := json.Unmarshal(b, station)
	if err != nil {
		return nil
	}
	return station
}

// GetServices returns all Services that have not yet departed that are within
// the specified time range.
// If from is before to then it's resumed the time range crosses midnight.
func (d *LDB) GetServices(s *Station, from *util.WorkingTime, to *util.WorkingTime) []*Service {
	var services []*Service

	_ = d.View(func(tx *bbolt.Tx) error {
		services = s.getServices(tx, from, to)
		return nil
	})

	return services
}

// GetMessages returns all station Messages for this Station.
func (s *Station) GetMessages(client *d3client.DarwinD3Client) []*darwind3.StationMessage {
	var messages []*darwind3.StationMessage

	for _, id := range s.Messages {
		if sm, _ := client.GetStationMessage(id); sm != nil {
			messages = append(messages, sm)
		}
	}

	return messages
}

func (s *Station) addStationMessage(msg *darwind3.StationMessage) {
	found := false
	for idx, i := range s.Messages {
		if i == msg.ID {
			s.Messages[idx] = msg.ID
			found = true
		}
	}

	if !found {
		s.Messages = append(s.Messages, msg.ID)
	}

	sort.SliceStable(s.Messages, func(i, j int) bool {
		return s.Messages[i] < s.Messages[j]
	})
}
