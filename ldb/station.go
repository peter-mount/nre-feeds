package ldb

import (
	"encoding/json"
	"github.com/peter-mount/nre-feeds/darwind3"
	d3client "github.com/peter-mount/nre-feeds/darwind3/client"
	"github.com/peter-mount/nre-feeds/darwinref"
	"github.com/peter-mount/nre-feeds/util"
	"log"
	"sort"
)

// The holder for a station's departure boards
type Station struct {
	// The location details for this station
	Locations []*darwinref.Location
	Crs       string
	// The Services at this station
	Services map[string]*Service
	// This station is Public - i.e. has a CRS so can have departures
	Public bool
	// The Station message id's applicable to this station
	Messages []uint64
}

// Only valid for Public stations, initialise it
func (s *Station) init() {
	s.Services = make(map[string]*Service)
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
func (s *Station) GetServices(from *util.WorkingTime, to *util.WorkingTime) []*Service {

	var services []*Service

	// Get a copy the Services from the station, filtering as needed
	for _, service := range s.Services {
		if !service.Location.Forecast.Departed && service.Location.Time.Between(from, to) {
			services = append(services, service.Clone())
		}
	}

	// sort into time order
	sort.SliceStable(services, func(i, j int) bool {
		return services[i].Compare(services[j])
	})

	log.Printf("GetServices %s %v to %v returning %d/%d services", s.Crs, from, to, len(services), len(s.Services))
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

func (s *Station) addStationMessage(msg *darwind3.StationMessage) bool {
	for _, i := range s.Messages {
		if i == msg.ID {
			return false
		}
	}
	s.Messages = append(s.Messages, msg.ID)
	return true
}
