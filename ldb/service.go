package ldb

import (
  "darwind3"
  "darwintimetable"
  "github.com/peter-mount/golib/statistics"
  "time"
)

// A representation of a service at a location
type Service struct {
  // The RID of this service
  RID               string                      `json:"rid"`
  // The destination
  Destination       string                      `json:"destination"`
  // Service Start Date
  SSD               darwintimetable.SSD         `json:"ssd"`
  // The trainId (headcode)
  TrainId           string                      `json:"trainId"`
  // The operator of this service
  Toc               string                      `json:"toc"`
  // Is a passenger service
  PassengerService  bool                        `json:"passengerService,omitempty"`
  // Is a charter service
  Charter           bool                        `json:"charter,omitempty"`
  // The cancel reason
  CancelReason      darwind3.DisruptionReason   `json:"cancelReason"`
  // The "time" for this service
  Location         *darwind3.Location           `json:"location"`
  // The time this entry was set
  Date              time.Time                   `json:"date,omitempty" xml:"date,attr,omitempty"`
  // URL to the train detail page
  Self              string                      `json:"self,omitempty" xml:"self,attr,omitempty"`
}

// Compare two services by the times at a location
func (a *Service) Compare( b *Service ) bool {
  return b != nil && a.Location.Compare( b.Location )
}

func (s *Service) update( e *darwind3.DarwinEvent ) bool {
  sched := e.Schedule
  if sched != nil && sched.Date.After( s.Date ) {
    s.RID = e.RID
    s.Location = e.Location

    s.SSD = sched.SSD
    s.TrainId = sched.TrainId
    s.Toc = sched.Toc
    s.PassengerService = sched.PassengerService
    s.CancelReason = sched.CancelReason

    // Resolve the destination
    if s.Location.FalseDestination != "" {
      s.Destination = s.Location.FalseDestination
    } else if len( sched.Locations ) > 0 {
      // For now presume this is correct
      s.Destination = sched.Locations[ len( sched.Locations )-1 ].Tiploc
    } else {
      s.Destination = ""
    }

    // Copy the date/self of the underlying schedule
    s.Date = sched.Date
    s.Self = sched.Self

    return true
  }

  return false
}

// Adds a service to the station
func (s *Station) addService( e *darwind3.DarwinEvent ) {
  statistics.Incr( "ldb.update" )

  s.Update( func() error {

    // See if we already have this train
    if len( s.services ) > 0 {
      for i, l := range s.services {
        if e.RID == l.RID && e.Location.EqualInSchedule( l.Location ) {
          if s.services[ i ].update( e ) {
            s.update()
          }
          return nil
        }
      }
    }

    service := &Service{}
    if service.update( e ) {
      s.services = append( s.services, service )
      s.update()
    }

    return nil
  })
}
