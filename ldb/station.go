package ldb

import (
  "darwind3"
  "darwinref"
  "sync"
)

// The holder for a station's departure boards
type Station struct {
  // The location details for this station
  Locations            []*darwinref.Location
  Crs                     string
  // The services at this station
  services                map[string]*Service
  // This station is public - i.e. has a CRS so can have departures
  public                  bool
  // The Station message id's applicable to this station
  messages              []int
  // Mutex for this station
  mutex                  *sync.Mutex
  // Update channel
  addServiceChannel     chan *stationAddService
  removeServiceChannel  chan string
  // Pointer to Stations object
  ldb                    *LDB
}

// Only valid for public stations, initialise it
func (s *Station) init() {
  s.services = make( map[string]*Service )
  s.mutex = &sync.Mutex{}

  // The addService channel & worker
  s.addServiceChannel = make( chan *stationAddService, 100 )
  go s.addServiceWorker()

  // The removeService channel & worker
  s.removeServiceChannel = make( chan string, 100 )
  go s.removeServiceWorker()

}

// Perform an action on the station with an exclusive lock
func (s *Station) Update( f func() error ) error {
  s.mutex.Lock()
  defer s.mutex.Unlock()
  return f()
}

func (s *Station) update() {
  s.ldb.Darwin.EventManager.PostEvent( &darwind3.DarwinEvent{
    Type: darwind3.Event_BoardUpdate,
    Crs: s.Crs,
  })
}
