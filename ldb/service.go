package ldb

import (
	"encoding/json"
	"github.com/peter-mount/nre-feeds/darwind3"
	"github.com/peter-mount/nre-feeds/util"
	"time"
)

// A representation of a service at a location
type Service struct {
	// The RID of this service
	RID string `json:"rid"`
	// The destination - use this and not Dest.Tiploc as this can be overridden
	// if Location.FalseDestination is set
	Destination string `json:"destination"`
	// Origin Location of this service
	Origin darwind3.Location `json:"origin"`
	// Destination Location of this service
	Dest darwind3.Location `json:"dest"`
	// Service Start Date
	SSD util.SSD `json:"ssd"`
	// The trainId (headcode)
	TrainId string `json:"trainId"`
	// The operator of this service
	Toc string `json:"toc"`
	// Is a passenger service
	PassengerService bool `json:"passengerService,omitempty"`
	// Is a charter service
	Charter bool `json:"charter,omitempty"`
	// Cancel running reason for this service. The reason applies to all locations
	// of this service which are marked as cancelled
	CancelReason darwind3.DisruptionReason `json:"cancelReason"`
	// Late running reason for this service. The reason applies to all locations
	// of this service which are not marked as cancelled
	LateReason darwind3.DisruptionReason `json:"lateReason"`
	// The "time" for this service
	Location darwind3.Location `json:"location"`
	// The calling points from this location
	CallingPoints []darwind3.CallingPoint `json:"calling"`
	// The last report
	LastReport darwind3.CallingPoint `json:"lastReport,omitempty"`
	// The associations
	Associations []*darwind3.Association `json:"association"`
	// The latest schedule entry used for this service
	schedule *darwind3.Schedule `json:"schedule"`
	// The index within the schedule of this location
	LocationIndex int `json:"locind"`
	// The time this entry was set
	Date time.Time `json:"date,omitempty" xml:"date,attr,omitempty"`
	// URL to the train detail page
	Self string `json:"self,omitempty" xml:"self,attr,omitempty"`
}

// Entry in the DB, just the essential data
type ServiceEntry struct {
	RID           string
	LocationIndex int
	Departed      bool
	Date          time.Time
}

// Bytes returns the message as an encoded byte slice
func (s *ServiceEntry) Bytes() ([]byte, error) {
	b, err := json.Marshal(s)
	return b, err
}

// ScheduleFromBytes returns a schedule based on a slice or nil if none
func ServiceEntryFromBytes(b []byte) ServiceEntry {
	service := ServiceEntry{}
	if b != nil {
		err := json.Unmarshal(b, &service)
		if err != nil {
			service.RID = ""
		}
	}
	return service
}

// Timestamp returns the time.Time of this service based on the SSD and Location's Time.
// TODO this does not currently handle midnight correctly
func (s *Service) Timestamp() time.Time {
	return s.SSD.Time().Add(time.Duration(s.Location.Forecast.Time.Get()) * time.Second)
}

func (s *ServiceEntry) update(sched *darwind3.Schedule, idx int) bool {
	if sched == nil {
		return false
	}

	if !s.Date.IsZero() && !sched.Date.IsZero() && !s.Date.Before(sched.Date) {
		return false
	}

	if (s.RID == "" || s.RID == sched.RID) && idx >= 0 && idx < len(sched.Locations) {
		s.RID = sched.RID
		s.LocationIndex = idx
		s.Date = sched.Date

		return true
	}

	return false
}

func (s *Service) Update(sched *darwind3.Schedule, idx int) bool {

	if sched == nil {
		return false
	}

	// Check the schedule has been updated
	if !s.Date.IsZero() && !sched.Date.IsZero() && !s.Date.Before(sched.Date) {
		return false
	}

	if (s.RID == "" || s.RID == sched.RID) && idx >= 0 && idx < len(sched.Locations) {

		// Copy of our meta data
		s.schedule = sched
		s.LocationIndex = idx

		// Clear calling points so we'll update again later when needed
		s.CallingPoints = nil

		s.RID = sched.RID

		// Clone the location
		if sched.Locations != nil {
			s.Location = *sched.Locations[idx]
			s.Location.UpdateTime()
		}

		s.SSD = sched.SSD
		s.TrainId = sched.TrainId
		s.Toc = sched.Toc
		s.PassengerService = sched.PassengerService
		s.CancelReason = sched.CancelReason
		s.LateReason = sched.LateReason

		// The origin/destination Locations
		if sched.Origin != nil {
			s.Origin = *sched.Origin
		}
		if sched.Destination != nil {
			s.Dest = *sched.Destination
		}

		// Resolve the destination
		if s.Location.FalseDestination != "" {
			s.Destination = s.Location.FalseDestination
		} else if sched.Destination != nil {
			s.Destination = sched.Destination.Tiploc
		}
		if s.Destination == "" && len(sched.Locations) > 0 {
			// Use last location if no destination
			s.Destination = sched.Locations[len(sched.Locations)-1].Tiploc
		}

		// Copy associations
		s.Associations = sched.Associations

		// Copy the date/self of the underlying schedule
		s.Date = sched.Date
		s.Self = sched.Self

		return true
	}

	return false
}

// Clone returns a copy of this Service
func (a *Service) Clone() *Service {
	return &Service{
		RID:              a.RID,
		Destination:      a.Destination,
		Origin:           a.Origin,
		Dest:             a.Dest,
		SSD:              a.SSD,
		TrainId:          a.TrainId,
		Toc:              a.Toc,
		PassengerService: a.PassengerService,
		Charter:          a.Charter,
		CancelReason:     a.CancelReason,
		LateReason:       a.LateReason,
		Location:         a.Location,
		Associations:     a.Associations,
		schedule:         a.schedule,
		LocationIndex:    a.LocationIndex,
		Date:             a.Date,
		Self:             a.Self,
	}
}
