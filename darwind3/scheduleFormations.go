package darwind3

import (
	"github.com/etcd-io/bbolt"
	"log"
	"time"
)

// Coach details for the schedule, containing class type & toilet details.
// defined in rttiPPTFormations_v2 Loading but shares the same object as Loading
type ScheduleFormation struct {
	RID       string    `json:"rid" xml:"rid,attr"`
	Formation Formation `json:"formation" xml:"formation"`
	Date      time.Time `json:"date,omitempty"`
}

type Formation struct {
	Fid     string           `json:"fid" xml:"fid,attr"`
	Src     string           `json:"src,omitempty" xml:"src,attr,omitempty"`
	SrcInst string           `json:"srcInst,omitempty" xml:"srcInst,attr,omitempty"`
	Coaches []CoachFormation `json:"coaches" xml:"coaches>coach"`
	Date    time.Time        `json:"date,omitempty"`
}

// The CoachData & CoachLoadingData/LoadingValue complexTypes in rttiPPTFormations_v1
// We share the object to keep things simple.
type CoachFormation struct {
	CoachNumber string `json:"coachNumber" xml:"coachNumber,attr"`
	CoachClass  string `json:"coachClass,omitempty" xml:"coachClass,attr,omitempty"`
	Toilet      Toilet `json:"toilet" xml:"toilet,omitempty"`
}

// Process processes an inbound loading element containing train formation data.
func (l *ScheduleFormation) Process(tx *Transaction) error {
	return tx.d3.UpdateBulkAware(func(dbtx *bbolt.Tx) error {
		return l.process(tx, dbtx)
	})
}

func (l *ScheduleFormation) process(tx *Transaction, dbtx *bbolt.Tx) error {
	// Retrieve the schedule to be updated
	sched := GetSchedule(dbtx, l.RID)

	// No schedule then try to fetch it from the timetable
	if sched == nil {
		sched = tx.ResolveSchedule(l.RID)
	}

	// If no schedule then warn as we need UID & SSD but don't have it in the
	// Loading message
	if sched == nil {
		log.Println("Unknown RID in Loading", l.RID)
		return nil
	}

	sched.Formation = *l
	sched.Formation.Date = tx.pport.TS
	sched.Date = tx.pport.TS

	if PutSchedule(dbtx, sched) {
		tx.d3.EventManager.PostEvent(&DarwinEvent{
			Type:     Event_ScheduleUpdated,
			RID:      sched.RID,
			Schedule: sched,
		})
	}
	return nil
}
