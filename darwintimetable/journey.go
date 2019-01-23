// Reference timetable
package darwintimetable

import (
  bolt "github.com/etcd-io/bbolt"
  "encoding/json"
  "encoding/xml"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/nre-feeds/util"
  "time"
)

type Journey struct {
  XMLName         xml.Name      `json:"-" xml:"Journey"`
  // RTTI unique Train ID
  RID             string        `json:"rid" xml:"rid,attr"`
  // Train UID
  UID             string        `json:"uid" xml:"uid,attr"`
  // Train ID (Headcode)
  TrainID         string        `json:"trainId" xml:"trainId"`
  // Scheduled Start Date
  SSD             util.SSD      `json:"ssd" xml:"ssd,attr"`
  // ATOC Code
  Toc             string        `json:"toc" xml:"toc,attr"`
  // Type of service, i.e. Train/Bus/Ship.
  Status          string        `json:"status,omitempty" xml:"status,attr,omitempty"`
  // Category of service.
  TrainCat        string        `json:"trainCat" xml:"trainCat,attr"`
  // True if Darwin classifies the train category as a passenger service.
  Passenger       bool          `json:"isPassengerSvc" xml:"isPassengerSvc,attr"`
  // Service has been deleted and should not be used/displayed.
  Deleted         bool          `json:"deleted,omitempty" xml:"deleted,attr,omitempty"`
  // Indicates if this service is a charter service.
  Charter         bool          `json:"isCharter,omitempty" xml:"isCharter,attr,omitempty"`
  // True if this is a Q Train (runs as required) that has not yet been activated.
  // Note that a Q Train that has been activated before the XML Timetable file
  // has been built will not have this attribute set true.
  QTrain          bool          `json:"qtrain,omitempty" xml:"qtrain,attr,omitempty"`
  // The schedule
  Schedule      []*Location     `json:"locations" xml:location`
  CancelReason    int           `json:"cancelReason" xml:"cancelReason,attr"`
  // Associations
  //Associations  []*Association  `xml:"-"`
  // Date entry was inserted into the database
  Date        time.Time `json:"date" xml:"date,attr"`
  // URL to this entity
  Self        string    `json:"self" xml:"self,attr,omitempty"`
}

type cancelReason struct {
  text string `xml:",chardata"`
}

func (a *Journey) Equals( b *Journey ) bool {
  if b == nil {
    return false
  }
  return a.RID == b.RID &&
    a.UID == b.UID &&
    a.TrainID == b.TrainID &&
    a.SSD == b.SSD &&
    a.Toc == b.Toc &&
    a.TrainCat == b.TrainCat &&
    a.Passenger == b.Passenger &&
    a.CancelReason == b.CancelReason
}

func (t *Journey) SetSelf( r *rest.Rest ) {
  t.Self = r.Self( r.Context() + "/journey/" + t.RID )
}

// GetJourney returns details of a Journey
func (r *DarwinTimetable) GetJourney( tx *bolt.Tx, rid string ) ( *Journey, bool ) {
  loc, exists := r.GetJourneyBucket( tx.Bucket( []byte("DarwinJourney") ), rid )
  return loc, exists
}

// Internal method that uses the shared writable transaction
func (r *DarwinTimetable) getJourney( rid string ) ( *Journey, bool ) {
  loc, exists := r.GetJourneyBucket( r.journeys, rid )
  return loc, exists
}

func (r *DarwinTimetable) GetJourneyBucket( bucket *bolt.Bucket, rid string ) ( *Journey, bool ) {
  b := bucket.Get( []byte( rid ) )

  if b != nil {
    var journey *Journey = &Journey{}
    err := json.Unmarshal( b, journey )
    if err != nil {
      return nil, false
    }
    return journey, true
  }

  return nil, false
}

func (r *DarwinTimetable) addJourney( journey *Journey ) ( error, bool ) {
  // Update only if it does not exist or is different
  if old, exists := r.getJourney( journey.RID ); !exists || !journey.Equals( old ) {
    journey.Date = time.Now()

    b, err := json.Marshal( journey )
    if err != nil {
      return err, false
    }

    err = r.journeys.Put( []byte( journey.RID ), b )
    if err != nil {
      return err, false
    }

    return nil, true
  }

  return nil, false
}
