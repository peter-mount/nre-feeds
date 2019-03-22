package darwind3

import (
  "time"
)

// TimeTable updates
type TimeTableId struct {
  // Unique Timetable identifier
  TimeTableId         string              `json:"timeTableId" xml:",chardata"`
  // Timetable filename
  TTFile              string              `json:"ttfile,omitempty" xml:"ttfile,attr,omitempty"`
  // Reference filename
  TTRefFile           string              `json:"ttreffile,omitempty" xml:"ttreffile,attr,omitempty"`
  // Timestamp of this event
  Date                time.Time           `json:"date" xml:"-"`
}

// All we do is send it out as a Event_TimeTableUpdate event.
func (p *TimeTableId) Process( tx *Transaction ) error {

  p.Date = tx.pport.TS

  // Remove " " as they have to send it that way as the attributes are mandatory
  // in the v16 schemas
  if p.TTFile == " " {
    p.TTFile = ""
  }

  if p.TTRefFile == " " {
    p.TTRefFile = ""
  }

  // Send the event
  tx.d3.EventManager.PostEvent( &DarwinEvent{
    Type: Event_TimeTableUpdate,
    TimeTableId: p,
  })

  return nil
}
