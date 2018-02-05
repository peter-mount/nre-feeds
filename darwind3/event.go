package darwind3

import (
  "log"
  "sync"
)

// The possible types of DarwinEvent
const (
  // A schedule was deactivated
  Event_Deactivated = iota
  // A schedule was updated
  Event_ScheduleUpdated
  // A new StationMessage
  Event_StationMessage
  // A station's departure boards have been updated (LDB only)
  Event_BoardUpdate
)

// An event notifying of something happening within DarwinD3
type DarwinEvent struct {
  // The type of the event
  Type        int
  // The RID of the train that caused this event
  RID         string
  // The affected Schedule or nil if none
  Schedule   *Schedule
  // The CRS code of the station in this event (LDB only)
  Crs         string
  // The StationMessage that's been updated
  NewStationMessage        *StationMessage
  // The existing message before the update (or nil)
  ExistingStationMessage   *StationMessage
}

// The core of the eventing system
type DarwinEventManager struct {
  mutex        *sync.Mutex
  listeners     map[int][]*darwinEventListener
  listenerSeq   int
}

type darwinEventListener struct {
  sequence       int
  eventType      int
  channel   chan *DarwinEvent
}

// NewDarwinEventManager creates a new DarwinEventManager
func NewDarwinEventManager() *DarwinEventManager {
  d := &DarwinEventManager{}
  d.mutex = &sync.Mutex{}
  d.listeners = make( map[int][]*darwinEventListener )
  return d
}

// ListenToEvents will run a function which will reveive DarwinEvent's for the
// specified type until it exists.
func (d *DarwinEventManager) ListenToEvents( eventType int, f func( chan *DarwinEvent ) ) {
  d.ListenToEventsCapacity( eventType, 1000, f )
}

func (d *DarwinEventManager) ListenToEventsCapacity( eventType int, capacity int, f func( chan *DarwinEvent ) ) {

  d.mutex.Lock()
  defer d.mutex.Unlock()

  listeners := d.listeners[ eventType ]

  l := &darwinEventListener{
    sequence: d.listenerSeq,
    eventType: eventType,
    channel: make( chan *DarwinEvent, capacity ),
  }
  d.listenerSeq++
  listeners = append( listeners, l )

  d.listeners[ eventType ] = listeners

  // Launch the listener
  go func() {
    // Ensure we deregister if the listener panics
    defer func() {
      if err := recover(); err != nil {
        log.Println( err )
      }
      d.deregisterEventListener( l )
    }()

    f( l.channel )
  }()

}

// DeregisterEventListener removes a channel from receiving events
func (d *DarwinEventManager) deregisterEventListener( l *darwinEventListener ) {

  d.mutex.Lock()
  defer d.mutex.Unlock()

  if listeners, ok := d.listeners[ l.eventType ]; ok {
    var arr []*darwinEventListener
    for _, oc := range listeners {
      if oc.sequence != l.sequence {
        arr = append( arr, oc )
      }
    }
    d.listeners[ l.eventType ] = arr
  }
}

// PostEvent posts a DarwinEvent to all listeners listening for that specific type
func (d *DarwinEventManager) PostEvent( e *DarwinEvent ) {

  d.mutex.Lock()
  listeners, ok := d.listeners[ e.Type ]
  d.mutex.Unlock()

  if ok  {

    for _, l := range listeners {

      lc := len( l.channel )
      if lc >= 750 {
        log.Println( "evt", l.sequence, lc )
      }
      l.channel <- e
    }
  }
}
