// Reference timetable
package darwintimetable

import (
  bolt "github.com/coreos/bbolt"
  "encoding/xml"
  "github.com/peter-mount/golib/codec"
  "github.com/peter-mount/golib/rest"
  "time"
  "util"
)

type Journey struct {
  XMLName         xml.Name      `json:"-" xml:"Journey"`
  RID             string        `json:"rid" xml:"rid,attr"`
  UID             string        `json:"uid" xml:"uid,attr"`
  TrainID         string        `json:"trainId" xml:"trainId"`
  SSD             util.SSD      `json:"ssd" xml:"ssd,attr"`
  Toc             string        `json:"toc" xml:"toc,attr"`
  TrainCat        string        `json:"trainCat" xml:"trainCat,attr"`
  Passenger       bool          `json:"isPassengerSvc" xml:"isPassengerSvc,attr"`
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

func (t *Journey) Write( c *codec.BinaryCodec ) {
  c.WriteString( t.RID ).
    WriteString( t.UID ).
    WriteString( t.TrainID ).
    Write( &t.SSD ).
    WriteString( t.Toc ).
    WriteString( t.TrainCat ).
    WriteBool( t.Passenger ).
    WriteInt( t.CancelReason ).
    WriteTime( t.Date )

  c.WriteInt16( int16( len( t.Schedule ) ) )
  for _, l := range t.Schedule {
    c.Write( l )
  }
}

func (t *Journey) Read( c *codec.BinaryCodec ) {
  c.ReadString( &t.RID ).
    ReadString( &t.UID ).
    ReadString( &t.TrainID ).
    Read( &t.SSD ).
    ReadString( &t.Toc ).
    ReadString( &t.TrainCat ).
    ReadBool( &t.Passenger ).
    ReadInt( &t.CancelReason ).
    ReadTime( &t.Date )

  var lc int16
  c.ReadInt16( &lc )
  for i := 0; i < int(lc); i++ {
    l := &Location{}
    c.Read( l )
    t.Schedule = append( t.Schedule, l )
  }
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

func (t *Journey) fromBytes( b []byte ) bool {
  if b != nil {
    codec.NewBinaryCodecFrom( b ).Read( t )
  }
  return t.RID != ""
}

func (r *DarwinTimetable) GetJourneyBucket( bucket *bolt.Bucket, rid string ) ( *Journey, bool ) {
  b := bucket.Get( []byte( rid ) )

  if b != nil {
    var journey *Journey = &Journey{}
    if journey.fromBytes( b ) {
      return journey, true
    }
  }

  return nil, false
}

func (r *DarwinTimetable) addJourney( journey *Journey ) ( error, bool ) {
  // Update only if it does not exist or is different
  if old, exists := r.getJourney( journey.RID ); !exists || !journey.Equals( old ) {
    journey.Date = time.Now()
    codec := codec.NewBinaryCodec()
    codec.Write( journey )
    if codec.Error() != nil {
      return codec.Error(), false
    }

    if err := r.journeys.Put( []byte( journey.RID ), codec.Bytes() ); err != nil {
      return err, false
    }

    return nil, true
  }

  return nil, false
}
