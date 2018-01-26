package ldb

import (
  "darwinref"
  "fmt"
  "sort"
  "sync"
)

// The holder for a station's departure boards
type Station struct {
  // The location details for this station
  Locations  []*darwinref.Location
  // The services at this station
  services   []*Service
  // This station is public - i.e. has a CRS so can have departures
  public        bool
  // Mutex for this station
  mutex        *sync.Mutex
  // If true then station needs persisting
  updated       bool
}

func (s *Station) String() string {
  return fmt.Sprintf( "CRS %s Services %d Updated %v", s.Locations[0].Crs, len( s.services ), s.updated )
}

// Perform an action on the station with an exclusive lock
func (s *Station) Update( f func() error ) error {
  s.mutex.Lock()
  defer s.mutex.Unlock()
  return f()
}

func (s *Station) update() {
  sort.SliceStable( s.services, func( i, j int ) bool {
    return s.services[ i ].Compare( s.services[ j ] )
  } )
  s.updated = true
}
