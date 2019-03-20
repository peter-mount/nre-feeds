package darwind3

import (
  "time"
)

// TimeTable updates
type TrackingID struct {
  // The berth details
  Berth             TDBerth        `json:"berth" xml:"berth"`
  // The incorrect TrainID
  IncorrectTrainID  string         `json:"incorrectTrainID" xml:"incorrectTrainID"`
  // The correct TrainID
  CorrectTrainID    string         `json:"correctTrainID" xml:"correctTrainID"`
  // Timestamp of this event
  Date              time.Time      `json:"date" xml:"-"`
}

type TDBerth struct {
  // Train Describer (TD) Area
  Area              string         `json:"area" xml:"area,attr"`
  // TD Berth
  Berth             string         `json:"berth" xml:",chardata"`
}

// All we do is send it out as a Event_TrackingID event.
func (p *TrackingID) Process( tx *Transaction ) error {

  p.Date = tx.pport.TS

  tx.d3.EventManager.PostEvent( &DarwinEvent{
    Type: Event_TrackingID,
    TrackingID: p,
  })

  return nil
}
