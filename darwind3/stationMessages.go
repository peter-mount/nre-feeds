package darwind3

import (
	"fmt"
	"github.com/peter-mount/filecache"
	"log"
	"strconv"
	"time"
)

// StationMessages is an in-memory with disk backup of all received StationMessage's
// This is periodically cleared down as messages expire
type StationMessages struct {
	d3       *DarwinD3
	cacheDir string
	cache    string
}

const (
	StationmessageResynchronisation = -1000
)

func uint64key(id int64) string {
	return fmt.Sprintf("%d", id)
}

func (sm *StationMessages) AddMotd(id int64, text string) {
	if id == 0 {
		return
	}

	if id > 0 {
		id = -id
	}

	message := &StationMessage{
		Motd:     true,
		Category: "System",
		ID:       id,
		Message:  text,
		Date:     time.Now(),
		Severity: 0,
		Suppress: false,
		Active:   true,
	}

	sm.d3.StationMessages.Add(uint64key(message.ID), message)

	sm.d3.EventManager.PostEvent(&DarwinEvent{
		Type:              Event_StationMessage,
		NewStationMessage: message,
	})
}

func (sm *StationMessages) RemoveMotd(id int64) {
	if id == 0 {
		return
	}

	if id > 0 {
		id = -id
	}

	message := sm.Get(id)

	sm.d3.StationMessages.DeleteFromMemoryAndDisk(uint64key(id))

	// Simulate a delete from Darwin as a message with no stations
	sm.d3.EventManager.PostEvent(&DarwinEvent{
		Type:                   Event_StationMessage,
		ExistingStationMessage: message,
		NewStationMessage:      &StationMessage{ID: id},
	})

}

func (sm *StationMessages) ForEach(f func(*StationMessage) error) {
	sm.d3.StationMessages.Foreach(func(key string, item *filecache.CacheItem) {
		if item != nil {
			_ = f(item.Data().(*StationMessage))
		}
	})
}

// Get returns the specified StationMessage or nil if none
func (sm *StationMessages) Get(id int64) *StationMessage {
	v, _ := sm.d3.StationMessages.Get(strconv.FormatUint(uint64(id), 10))
	if v != nil {
		return v.Data().(*StationMessage)
	}
	return nil
}

func (sm *StationMessages) put(s *StationMessage) {
	if len(s.Station) > 0 {
		sm.d3.StationMessages.Add(uint64key(s.ID), s)
	} else {
		sm.d3.StationMessages.DeleteFromMemoryAndDisk(uint64key(s.ID))
	}
}

// BroadcastStationMessages sends all StationMessage's to the event queue as if they have
// just been received.
func (d *DarwinD3) BroadcastStationMessages(e *DarwinEvent) {
	cnt := 0
	d.Messages.ForEach(func(message *StationMessage) error {
		d.EventManager.PostEvent(&DarwinEvent{
			Type:              Event_StationMessage,
			NewStationMessage: message,
		})

		cnt++

		return nil
	})

	log.Println("Broadcast", cnt, "StationMessage's")
}

// ExpireStationMessages expires any old (>24 hours) station messages
// Note: this is only to keep the DB size down, they should delete automatically now
/*
func (d *DarwinD3) ExpireStationMessages() {
  cutoff := time.Now().Add(-24 * time.Hour)
  cnt := 0

  _ = d.Update(func(tx *bbolt.Tx) error {
    bucket := tx.Bucket([]byte(messageBucket))

    return bucket.ForEach(func(k, v []byte) error {
      message := StationMessageFromBytes(v)
      if message == nil {
        // Damaged message
        cnt++
        _ = bucket.Delete(k)
      } else if message.Date.Before(cutoff) {
        // Expired message
        cnt++
        _ = bucket.Delete(k)

        // Simulate a delete from Darwin as a message with no stations
        d.EventManager.PostEvent(&DarwinEvent{
          Type:                   Event_StationMessage,
          ExistingStationMessage: message,
          NewStationMessage: &StationMessage{
            ID:       message.ID,
            Message:  message.Message,
            Category: message.Category,
            Severity: message.Severity,
            Suppress: message.Suppress,
            Date:     message.Date,
          },
        })
      }
      return nil
    })
  })

  if cnt > 0 {
    log.Println("Expired", cnt, "StationMessage's")
  }
}

*/
