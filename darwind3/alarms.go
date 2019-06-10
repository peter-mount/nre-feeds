package darwind3

import (
	"encoding/json"
	"github.com/peter-mount/filecache"
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

func AlarmFromBytes(b []byte) interface{} {
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

func (d *DarwinD3) SetAlarm(a *Alarm) {
	log.Println("Set alarm", a)
	d.Alarms.Add(a.ID, a)
}

func (d *DarwinD3) GetAlarm(id string) *Alarm {
	v, _ := d.Alarms.Get(id)
	if v != nil {
		return v.Data().(*Alarm)
	}
	return nil
}

func (d *DarwinD3) GetAlarms() []*Alarm {
	var alarms []*Alarm

	d.Alarms.Foreach(func(key string, item *filecache.CacheItem) {
		if item != nil {
			alarms = append(alarms, item.Data().(*Alarm))
		}
	})

	return alarms
}

func (d *DarwinD3) DeleteAlarm(id string) {
	d.Alarms.DeleteFromMemoryAndDisk(id)
}

// Process processes an inbound Train Status update, merging it with an existing
// schedule in the database
func (p *RttiAlarm) Process(tx *Transaction) error {

	if p.Clear != "" {
		tx.d3.DeleteAlarm(p.Clear)

		tx.d3.EventManager.PostEvent(&DarwinEvent{
			Type:    Event_Alarm,
			AlarmId: p.Clear,
		})
	} else if p.Set.ID != "" {
		p.Set.Date = tx.pport.TS
		tx.d3.SetAlarm(&p.Set)

		tx.d3.EventManager.PostEvent(&DarwinEvent{
			Type:    Event_Alarm,
			AlarmId: p.Set.ID,
			Alarm:   &p.Set,
		})
	} else {
		log.Println("Unsupported Alarm", p)
	}

	return nil
}
