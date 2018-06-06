package ldb

import (
  "github.com/peter-mount/nre-feeds/darwind3"
  d3client "github.com/peter-mount/nre-feeds/darwind3/client"
  "github.com/peter-mount/nre-feeds/darwinref"
  "github.com/peter-mount/nre-feeds/util"
  "sort"
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
  // Pointer to Stations object
  ldb                    *LDB
}

// Only valid for public stations, initialise it
func (s *Station) init() {
  s.services = make( map[string]*Service )
  s.mutex = &sync.Mutex{}
}

// Perform an action on the station with an exclusive lock
func (s *Station) Update( f func() error ) error {
  s.mutex.Lock()
  defer s.mutex.Unlock()
  return f()
}

func (s *Station) update() {
  s.ldb.EventManager.PostEvent( &darwind3.DarwinEvent{
    Type: darwind3.Event_BoardUpdate,
    Crs: s.Crs,
  })
}

// GetServices returns all services that have not yet departed that are within
// the specified time range.
// If from is before to then it's resumed the time range crosses midnight.
func (s *Station) GetServices( from *util.WorkingTime, to *util.WorkingTime ) []*Service {

  var services []*Service

  // Get a copy the services from the station within the lock, filtering as needed
  s.Update( func() error {
    for _,service := range s.services {
      if !service.Location.Forecast.Departed && service.Location.Time.Between( from, to ) {
        services = append( services, service.Clone() )
      }
    }
    return nil
  } );

  // sort into time order
  sort.SliceStable( services, func( i, j int ) bool {
    return services[ i ].Compare( services[ j ] )
  } )

  return services
}

// GetMessages returns all station messages for this Station.
func (s *Station) GetMessages( client *d3client.DarwinD3Client ) []*darwind3.StationMessage {

  // Get a copy of the current id's within the lock
  var ids []int
  s.Update( func() error {
    for _, id := range s.messages {
      ids = append( ids, id )
    }
    return nil
  } );

  // Now resolve them outside the lock as this is a rest call
  var messages []*darwind3.StationMessage

  for _, id := range ids {
    if sm, _ := client.GetStationMessage( id ); sm != nil {
      messages = append( messages, sm )
    }
  }

  return messages
}
