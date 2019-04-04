package darwind3

import (
	"encoding/json"
	"github.com/etcd-io/bbolt"
	"log"
	"time"
)

type Alarm struct {
	ID     string    `json:"id" xml:"id,attr"`
	TDArea string    `json:"tdArea,omitempty" xml:"tdAreaFail"`
	Tyrell string    `json:"tyrell,omitempty" xml:"tyrellFeedFail"`
	Date   time.Time `json:"date" xml:"-"`
}

type RttiAlarm struct {
	Set   Alarm  `xml:"set"`
	Clear string `xml:"clear"`
}

func (a *Alarm) Bytes() ([]byte, error) {
	b, err := json.Marshal(a)
	return b, err
}
func AlarmFromBytes(b []byte) *Alarm {
	if b == nil {
		return nil
	}

	a := &Alarm{}
	err := json.Unmarshal(b, a)
	if err != nil || a.ID == "" {
		return nil
	}
	return a
}

func (d *DarwinD3) SetAlarm(a *Alarm) error {
	if d.cache.tx != nil {
		return d.setAlarm(d.cache.tx, a)
	}
	return d.Update(func(tx *bbolt.Tx) error {
		return d.setAlarm(d.cache.tx, a)
	})
}

func (d *DarwinD3) setAlarm(tx *bbolt.Tx, a *Alarm) error {
	b, err := a.Bytes()
	if err != nil {
		return err
	}
	return tx.Bucket([]byte(alarmBucket)).Put([]byte(a.ID), b)
}

func (d *DarwinD3) GetAlarm(id string) *Alarm {
	var a *Alarm

	_ = d.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(alarmBucket)).Get([]byte(a.ID))

		if b != nil {
			a = AlarmFromBytes(b)
		}

		return nil
	})

	return a
}

func (d *DarwinD3) GetAlarms() []*Alarm {
	var alarms []*Alarm

	_ = d.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(alarmBucket)).
			ForEach(func(k, v []byte) error {
				alarms = append(alarms, AlarmFromBytes(v))
				return nil
			})
	})

	return alarms
}

func (d *DarwinD3) DeleteAlarm(id string) error {
	if d.cache.tx != nil {
		return d.deleteAlarm(d.cache.tx, id)
	}
	return d.Update(func(tx *bbolt.Tx) error {
		return d.deleteAlarm(d.cache.tx, id)
	})
}

func (d *DarwinD3) deleteAlarm(tx *bbolt.Tx, id string) error {
	return tx.Bucket([]byte(alarmBucket)).Delete([]byte(id))
}

// Process processes an inbound Train Status update, merging it with an existing
// schedule in the database
func (p *RttiAlarm) Process(tx *Transaction) error {

	if p.Clear != "" {
		err := tx.d3.DeleteAlarm(p.Clear)

		tx.d3.EventManager.PostEvent(&DarwinEvent{
			Type:    Event_Alarm,
			AlarmId: p.Set.ID,
		})

		return err
	}

	if p.Set.ID != "" {
		p.Set.Date = tx.pport.TS
		err := tx.d3.SetAlarm(&p.Set)

		tx.d3.EventManager.PostEvent(&DarwinEvent{
			Type:    Event_Alarm,
			AlarmId: p.Set.ID,
			Alarm:   &p.Set,
		})

		return err
	}

	log.Println("Unsupported Alarm", p)

	return nil
}

func (d *DarwinD3) ExpireAlarms() {
	cutoff := time.Now().Add(-24 * 3 * time.Hour)
	cnt := 0

	_ = d.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(alarmBucket))

		return bucket.ForEach(func(k, v []byte) error {
			alarm := AlarmFromBytes(v)
			if alarm == nil || alarm.ID == "" {
				// Damaged alarm
				cnt++
				_ = bucket.Delete(k)
			} else if alarm.Date.Before(cutoff) {
				// Expired alarm
				cnt++
				_ = bucket.Delete(k)

				d.EventManager.PostEvent(&DarwinEvent{
					Type:    Event_Alarm,
					AlarmId: alarm.ID,
				})
			}
			return nil
		})
	})

	if cnt > 0 {
		log.Println("Expired", cnt, "Alarms's")
	}
}
