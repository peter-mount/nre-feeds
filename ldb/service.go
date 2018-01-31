package ldb

import (
  "darwind3"
  "darwintimetable"
  "log"
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
  // Cancel running reason for this service. The reason applies to all locations
  // of this service which are marked as cancelled
  CancelReason      darwind3.DisruptionReason   `json:"cancelReason"`
  // Late running reason for this service. The reason applies to all locations
  // of this service which are not marked as cancelled
  LateReason        darwind3.DisruptionReason   `json:"lateReason"`
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

func (s *Service) update( e *darwind3.DarwinEvent, loc *darwind3.Location ) bool {
  sched := e.Schedule
  if sched != nil && sched.Date.After( s.Date ) {
    s.RID = e.RID
    s.Location = loc

    s.SSD = sched.SSD
    s.TrainId = sched.TrainId
    s.Toc = sched.Toc
    s.PassengerService = sched.PassengerService
    s.CancelReason = sched.CancelReason
    s.LateReason = sched.LateReason

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
func (s *Station) addService( e *darwind3.DarwinEvent, loc *darwind3.Location ) {
  // Only public stations can be updated. Pass to the channel so the worker thread
  // can read it
  if s.public && loc.Times.IsPublic() {
    s.addServiceChannel <- &stationAddService{ e: e, loc: loc }
  }
}

type stationAddService struct {
  e   *darwind3.DarwinEvent
  loc *darwind3.Location
}

// Adds a service to the station
func (s *Station) addServiceWorker() {
  for {
     e := <- s.addServiceChannel

     cl := len( s.addServiceChannel )
     if cl > 75 {
       log.Println("Station", s.Locations[0].Crs, cl )
     }

     s.Update( func() error {

       // See if we already have this train
       if l, exists := s.services[ e.e.RID ]; exists {
         if l.update( e.e, e.loc ) {
           s.update()
         }
         return nil
       }

       // A new service
       service := &Service{}
       if service.update( e.e, e.loc ) {
         s.services[ e.e.RID ] = service
         s.update()
       }

       return nil
     })
  }
}
