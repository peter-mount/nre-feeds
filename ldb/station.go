package ldb

import (
  "darwinref"
  "fmt"
  "sync"
)

// The holder for a station's departure boards
type Station struct {
  // The location details for this station
  Locations        []*darwinref.Location
  // The services at this station
  services            map[string]*Service
  // This station is public - i.e. has a CRS so can have departures
  public              bool
  // Mutex for this station
  mutex              *sync.Mutex
  // If true then station needs persisting
  updated             bool
  // Update channel
  addServiceChannel   chan *stationAddService
}

func (s *Station) String() string {
  return fmt.Sprintf( "CRS %s Services %d Updated %v", s.Locations[0].Crs, len( s.services ), s.updated )
}

// Only valid for public stations, initialise it
func (s *Station) init() {
  s.services = make( map[string]*Service )
  s.mutex = &sync.Mutex{}

  // The addService channel & worker
  s.addServiceChannel = make( chan *stationAddService, 100 )
  go s.addServiceWorker()
}

// Perform an action on the station with an exclusive lock
func (s *Station) Update( f func() error ) error {
  s.mutex.Lock()
  defer s.mutex.Unlock()
  return f()
}

func (s *Station) update() {
  s.updated = true
}
