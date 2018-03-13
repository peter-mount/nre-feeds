package ldb

import (
  "bytes"
  "darwind3"
  "encoding/json"
  "time"
  "util"
)

// A representation of a service at a location
type Service struct {
  // The RID of this service
  RID               string                      `json:"rid"`
  // The destination
  Destination       string                      `json:"destination"`
  // Service Start Date
  SSD               util.SSD                    `json:"ssd"`
  // The trainId (headcode)
  TrainId           string                      `json:"trainId"`
  // The operator of this service
  Toc               string                      `json:"toc"`
  // Is a passenger service
  PassengerService  bool                        `json:"passengerService,omitempty"`
  // Is a charter service
  Charter           bool                        `json:"charter,omitempty"`
  // Cancel running reason for this service. The reason applies to all locations
  // of this service which are marked as cancelled
  CancelReason      darwind3.DisruptionReason   `json:"cancelReason"`
  // Late running reason for this service. The reason applies to all locations
  // of this service which are not marked as cancelled
  LateReason        darwind3.DisruptionReason   `json:"lateReason"`
  // The "time" for this service
  Location         *darwind3.Location           `json:"location"`
  // The calling points from this location
  CallingPoints  []*darwind3.CallingPoint       `json:"calling"`
  // The last report
  LastReport       *darwind3.CallingPoint       `json:"lastReport,omitempty"`
  // The latest schedule entry used for this service
  schedule         *darwind3.Schedule           `json:"-"`
  // The index within the schedule of this location
  locationIndex     int                         `json:"-"`
  // The time this entry was set
  Date              time.Time                   `json:"date,omitempty" xml:"date,attr,omitempty"`
  // URL to the train detail page
  Self              string                      `json:"self,omitempty" xml:"self,attr,omitempty"`
}

// Compare two services by the times at a location
func (a *Service) Compare( b *Service ) bool {
  return b != nil && a.Location.Compare( b.Location )
}

// Timestamp returns the time.Time of this service based on the SSD and Location's Time.
// TODO this does not currently handle midnight correctly
func (s *Service) Timestamp() time.Time {
  return s.SSD.Time().Add( time.Duration( s.Location.Forecast.Time.Get() ) * time.Second )
}

func (s *Service) update( sched *darwind3.Schedule, idx int ) bool {

  if sched != nil && //sched.Date.After( s.Date ) &&
      ( s.RID == "" || s.RID == sched.RID ) &&
      idx >=0 && idx < len( sched.Locations ) {

    // Copy of our meta data
    s.schedule = sched
    s.locationIndex = idx

    // Clear calling points so we'll update again later when needed
    s.CallingPoints = nil

    s.RID = sched.RID

    // Clone the location
    s.Location = sched.Locations[ idx ].Clone()

    s.SSD = sched.SSD
    s.TrainId = sched.TrainId
    s.Toc = sched.Toc
    s.PassengerService = sched.PassengerService
    s.CancelReason = sched.CancelReason
    s.LateReason = sched.LateReason

    // Resolve the destination
    if s.Location.FalseDestination != "" {
      s.Destination = s.Location.FalseDestination
    } else if len( sched.Locations ) > 0 {
      // For now presume this is correct
      s.Destination = sched.Locations[ len( sched.Locations )-1 ].Tiploc
    } else {
      s.Destination = ""
    }

    // Copy the date/self of the underlying schedule
    s.Date = sched.Date
    s.Self = sched.Self

    return true
  }

  return false
}

// Clone returns a copy of this Service
func (a *Service) Clone() *Service {
  return &Service{
    RID: a.RID,
    Destination: a.Destination,
    SSD: a.SSD,
    TrainId: a.TrainId,
    Toc: a.Toc,
    PassengerService: a.PassengerService,
    Charter: a.Charter,
    CancelReason: a.CancelReason,
    LateReason: a.LateReason,
    Location: a.Location.Clone(),
    schedule: a.schedule,
    locationIndex: a.locationIndex,
    Date: a.Date,
    Self: a.Self,
  }
}

func (t *Service) append( b *bytes.Buffer, c bool, f string, v interface{} ) bool {
  // Any null, "" or false ignore
  if vb, err := json.Marshal( v );
    err == nil &&
    !( len(vb) == 2 && vb[0] == '"' && vb[1] == '"' ) &&
    !( len(vb) == 4 && vb[0] == 'n' && vb[1] == 'u' && vb[2] == 'l' && vb[3] == 'l') &&
    !( len(vb) == 5 && vb[0] == 'f' && vb[1] == 'a' && vb[2] == 'l' && vb[3] == 's' && vb[4] == 'e') {
    if c {
      b.WriteByte( ',' )
    }

    b.WriteByte( '"' )
    b.WriteString( f )
    b.WriteByte( '"' )
    b.WriteByte( ':' )
    b.Write( vb )
    return true
  }

  return c
}

func (t *Service) MarshalJSON() ( []byte, error ) {
  var b bytes.Buffer

  b.WriteByte( '{' )
  c := t.append( &b, false, "rid", t.RID )
  c = t.append( &b, c, "destination", t.Destination )
  c = t.append( &b, c, "ssd", &t.SSD )
  c = t.append( &b, c, "trainId", t.TrainId )
  c = t.append( &b, c, "toc", t.Toc )
  c = t.append( &b, c, "passengerService", &t.PassengerService )
  c = t.append( &b, c, "charter", &t.Charter )

  if t.CancelReason.Reason > 0 {
    c = t.append( &b, c, "cancelReason", &t.CancelReason )
  }

  if t.LateReason.Reason > 0 {
    c = t.append( &b, c, "lateReason", &t.LateReason )
  }

  c = t.append( &b, c, "location", &t.Location )

  if len( t.CallingPoints ) > 0 {
    c = t.append( &b, c, "calling", t.CallingPoints )
  }

  if t.LastReport != nil {
    c = t.append( &b, c, "lastReport", t.LastReport )
  }

  c = t.append( &b, c, "date", t.Date )
  c = t.append( &b, c, "self", t.Self )

  b.WriteByte( '}' )
  return b.Bytes(), nil
}
