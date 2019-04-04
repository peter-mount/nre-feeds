package darwind3

import (
	"encoding/binary"
	"github.com/etcd-io/bbolt"
	"log"
	"sync"
	"time"
)

// StationMessages is an in-memory with disk backup of all received StationMessage's
// This is periodically cleared down as messages expire
type StationMessages struct {
	d3       *DarwinD3
	mutex    sync.Mutex
	messages map[int]*StationMessage
	cacheDir string
	cache    string
}

func uint64key(id uint64) []byte {
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, id)
	return key
}

func (sm *StationMessages) ForEach(f func(*StationMessage) error) error {
	return sm.d3.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(messageBucket))
		return bucket.ForEach(func(k, v []byte) error {
			message := StationMessageFromBytes(v)
			if message == nil {
				bucket.Delete(k)
			} else {
				err := f(message)
				if err != nil {
					return err
				}
			}
			return nil
		})
	})
}

// Get returns the specified StationMessage or nil if none
func (sm *StationMessages) Get(id uint64) *StationMessage {
	var s *StationMessage

	_ = sm.d3.View(func(tx *bbolt.Tx) error {
		s = sm.get(tx, id)
		return nil
	})

	return s
}

func (sm *StationMessages) get(tx *bbolt.Tx, id uint64) *StationMessage {
	bucket := tx.Bucket([]byte(messageBucket))

	b := bucket.Get(uint64key(id))
	if b != nil {
		return StationMessageFromBytes(b)
	}
	return nil
}

// Put stores a StationMessage or deletes it if it has no applicable stations
func (sm *StationMessages) Put(s *StationMessage) error {
	// Check for the snapshot transaction being open. If so then use that
	if sm.d3.cache.tx != nil {
		return sm.put(sm.d3.cache.tx, s)
	}

	return sm.d3.Update(func(tx *bbolt.Tx) error {
		return sm.put(tx, s)
	})
}
func (sm *StationMessages) put(tx *bbolt.Tx, s *StationMessage) error {
	bucket := tx.Bucket([]byte(messageBucket))

	key := uint64key(s.ID)

	if len(s.Station) > 0 {
		b, err := s.Bytes()
		if err != nil {
			return err
		}
		return bucket.Put(key, b)
	} else {
		return bucket.Delete(key)
	}
}

// BroadcastStationMessages sends all StationMessage's to the event queue as if they have
// just been received.
func (d *DarwinD3) BroadcastStationMessages() {
	cnt := 0
	_ = d.Messages.ForEach(func(message *StationMessage) error {
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

				d.EventManager.PostEvent(&DarwinEvent{
					Type:                   Event_StationMessage,
					ExistingStationMessage: message,
					// Simulate a delete from Darwin as a message with no stations
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
