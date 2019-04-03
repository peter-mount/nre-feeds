package darwind3

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

// StationMessages is an in-memory with disk backup of all received StationMessage's
// This is periodically cleared down as messages expire
type StationMessages struct {
	mutex    sync.Mutex
	messages map[int]*StationMessage
	cacheDir string
	cache    string
}

func NewStationMessages(cacheDir string) *StationMessages {
	s := &StationMessages{}
	s.messages = make(map[int]*StationMessage)
	s.cacheDir = cacheDir
	s.cache = s.cacheDir + "/stationMessages.dat"

	err := s.Load()
	if err != nil {
		log.Println(err)
	}

	return s
}

func (sm *StationMessages) Update(f func() error) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	return f()
}

func (sm *StationMessages) ForEach(f func(*StationMessage) error) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for _, m := range sm.messages {
		err := f(m)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get returns the specified StationMessage or nil if none
func (sm *StationMessages) Get(id int) *StationMessage {
	var s *StationMessage

	_ = sm.Update(func() error {
		s = sm.messages[id]
		return nil
	})

	return s
}

// Put stores a StationMessage or deletes it if it has no applicable stations
func (sm *StationMessages) Put(s *StationMessage) error {
	_ = sm.Update(func() error {
		if len(s.Station) > 0 {
			sm.messages[s.ID] = s
		} else {
			delete(sm.messages, s.ID)
		}

		return nil
	})

	return sm.Persist()
}

// BroadcastStationMessages sends all StationMessage's to the event queue as if they have
// just been received.
func (d *DarwinD3) BroadcastStationMessages() {
	_ = d.Messages.Update(func() error {
		if len(d.Messages.messages) > 0 {
			for _, s := range d.Messages.messages {
				d.EventManager.PostEvent(&DarwinEvent{
					Type:              Event_StationMessage,
					NewStationMessage: s,
				})
			}

			log.Println("Broadcast", len(d.Messages.messages), "StationMessage's")
		}
		return nil
	})
}

// ExpireStationMessages expires any old (>6 hours) station messages
func (d *DarwinD3) ExpireStationMessages() {
	_ = d.Messages.Update(func() error {
		cutoff := time.Now().Add(-6 * time.Hour)
		cnt := 0

		for id, s := range d.Messages.messages {
			if s.Date.Before(cutoff) {
				cnt++
				delete(d.Messages.messages, id)

				d.EventManager.PostEvent(&DarwinEvent{
					Type:                   Event_StationMessage,
					ExistingStationMessage: s,
					// Simulate a delete from Darwin as a message with no stations
					NewStationMessage: &StationMessage{
						ID:       s.ID,
						Message:  s.Message,
						Category: s.Category,
						Severity: s.Severity,
						Suppress: s.Suppress,
						Date:     s.Date,
					},
				})
			}
		}

		if cnt > 0 {
			log.Println("Expired", cnt, "StationMessage's")
		}

		return nil
	})

	_ = d.Messages.Persist()
}

// Load reloads the station messages from disk
func (sm *StationMessages) Load() error {
	return sm.Update(func() error {
		log.Println("Loading StationMessages")

		err := os.MkdirAll(sm.cacheDir, 0777)
		if err != nil {
			return err
		}

		buf, err := ioutil.ReadFile(sm.cache)
		if err != nil {
			return err
		}

		var messages []*StationMessage
		err = json.Unmarshal(buf, &messages)
		if err != nil {
			return err
		}

		for _, m := range messages {
			sm.messages[m.ID] = m
		}

		log.Println("Loaded", len(sm.messages), "StationMessage's")
		return nil
	})
}

// Persist stores all StationMessage's to disk
func (sm *StationMessages) Persist() error {
	return sm.Update(func() error {
		err := os.MkdirAll(sm.cacheDir, 0777)
		if err != nil {
			return err
		}

		var messages []*StationMessage
		for _, m := range sm.messages {
			messages = append(messages, m)
		}

		b, err := json.Marshal(&messages)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(sm.cache, b, 0655)
	})
}
