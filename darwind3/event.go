package darwind3

import (
  "log"
  "sync"
)

// The possible types of DarwinEvent
const (
  // A schedule was activated
  Event_Activated = iota
  // A schedule was deactivated
  Event_Deactivated
  // A schedule was updated
  Event_ScheduleUpdated
  // A location was updated
  Event_LocationUpdated
)

// An event notifying of something happening within DarwinD3
type DarwinEvent struct {
  // The type of the event
  Type        int
  // The RID of the train that caused this event
  RID         string
  // The affected Schedule or nil if none
  Schedule   *Schedule
  // The location of this event or nil
  Location   *Location
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

  d.mutex.Lock()
  defer d.mutex.Unlock()

  listeners := d.listeners[ eventType ]

  l := &darwinEventListener{
    sequence: d.listenerSeq,
    eventType: eventType,
    channel: make( chan *DarwinEvent, 1000 ),
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

  if ok {
    for _, l := range listeners {
      l.channel <- e
    }
  }
}
