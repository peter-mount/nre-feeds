package darwind3

import (
  "encoding/json"
  "encoding/xml"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/nre-feeds/util"
  "sort"
  "sync"
  "time"
)

// Train Schedule
type Schedule struct {
  RID               string                `json:"rid"`
  UID               string                `json:"uid"`
  TrainId           string                `json:"trainId"`
  SSD               util.SSD              `json:"ssd"`
  // The Train Operating Company
  Toc               string                `json:"toc"`
  // Default P
  Status            string                `json:"status"`
  // Default OO
  TrainCat          string                `json:"trainCat"`
  // Default true
  PassengerService  bool                  `json:"passengerService,omitempty"`
  // Default true
  Active            bool                  `json:"active,omitempty"`
  // Default false
  Deleted           bool                  `json:"deleted,omitempty"`
  // Default false
  Charter           bool                  `json:"charter,omitempty"`
  // Cancel running reason for this service. The reason applies to all locations
  // of this service which are marked as cancelled
  CancelReason      DisruptionReason      `json:"cancelReason"`
  // Late running reason for this service. The reason applies to all locations
  // of this service which are not marked as cancelled
  LateReason        DisruptionReason      `json:"lateReason"`
  // The locations in this schedule
  Locations      []*Location             `json:"locations"`
  // The origin of this service
  Origin           *Location              `json:"originLocation"`
  // The destination of this service
  Destintion       *Location              `json:"destinationLocation"`
  // Usually this is the date we insert into the db but here we use the TS time
  // as returned from darwin
  Date              time.Time             `json:"date,omitempty" xml:"date,attr,omitempty"`
  // URL to this entity
  Self              string                `json:"self,omitempty" xml:"self,attr,omitempty"`
  mutex             sync.RWMutex
}

// Update runs a function within a Write lock within the schedule
func (s *Schedule) Update( f func() error ) error {
  s.mutex.Lock()
  defer s.mutex.Unlock()
  return f()
}

// View runs a function within a Read lock within the schedule
func (s *Schedule) View( f func() error ) error {
  s.mutex.RLock()
  defer s.mutex.RUnlock()
  return f()
}

func (a *Schedule) Clone() *Schedule {
  var b *Schedule

  a.Update( func() error {
    b = &Schedule{
      RID: a.RID,
      UID: a.UID,
      TrainId: a.TrainId,
      SSD: a.SSD,
      Toc: a.Toc,
      Status: a.Status,
      TrainCat: a.TrainCat,
      PassengerService: a.PassengerService,
      Active: a.Active,
      Deleted: a.Deleted,
      Charter: a.Charter,
      CancelReason: a.CancelReason,
      LateReason: a.LateReason,
    }

    for _, l := range a.Locations {
      b.Locations = append( b.Locations, l.Clone() )
    }

    return nil
  })

  b.UpdateTime()
  return b
}

func (s *Schedule) SetSelf( r *rest.Rest ) {
  s.Self = r.Self( r.Context() + "/schedule/" + s.RID )
}

// Sort sorts the locations in a schedule by time order
func (s *Schedule) Sort() {
  s.UpdateTime()
  sort.SliceStable( s.Locations, func( i, j int ) bool {
    return s.Locations[ i ].Compare( s.Locations[ j ] )
  } )
}

func (a *Schedule) Equals( b *Schedule ) bool {
  r := b != nil &&
       a.RID == b.RID &&
       a.UID == b.UID &&
       a.TrainId == b.TrainId &&
       a.SSD.Equals( &b.SSD ) &&
       a.Toc == b.Toc &&
       a.Status == b.Status &&
       a.TrainCat == b.TrainCat &&
       a.PassengerService == b.PassengerService &&
       a.Active == b.Active &&
       a.Deleted == b.Deleted &&
       a.Charter == b.Charter &&
       a.CancelReason.Equals( &b.CancelReason ) &&
       len( a.Locations ) == len( b.Locations ) &&
       a.Date == b.Date

  if r {
    // This works as we've already confirmed the length
    for i, l := range a.Locations {
      if !l.Equals( b.Locations[i] ) {
        return false
      }
    }
  }

  return r
}

// ScheduleFromBytes returns a schedule based on a slice or nil if none
func ScheduleFromBytes( b []byte ) *Schedule {
  if b == nil {
    return nil
  }

  sched := &Schedule{}
  err := json.Unmarshal( b, sched )
  if err != nil || sched.RID == "" {
    return nil
  }
  sched.UpdateTime()
  return sched
}

func (s *Schedule) UpdateTime() {
  s.Origin = nil
  s.Destintion = nil

  for _, l := range s.Locations {
    l.UpdateTime()

    // Set origin or destination.
    // Note: OR & DT override OPOR & OPDT
    switch l.Type {
      case "OR":
        if s.Origin == nil || s.Origin.Type == "OPOR" {
          s.Origin = l
        }

      case "OPOR":
        if s.Origin == nil || s.Origin.Type != "OR" {
          s.Origin = l
        }

      case "DT":
        if s.Destintion == nil || s.Destintion.Type == "OPDT" {
          s.Destintion = l
        }

      case "OPDT":
        if s.Destintion == nil || s.Destintion.Type != "DT" {
          s.Destintion = l
        }
    }
  }

}

// Bytes returns the schedule as an encoded byte slice
func (s *Schedule) Bytes() ( []byte, error ) {
  b, err := json.Marshal( s )
  return b, err
}

// Defaults sets the default values for a schedule
func (s *Schedule) Defaults() {
  s.Status = "P"
  s.TrainCat = "OO"
  s.PassengerService = true
  s.Active = true
}

func (s *Schedule) UnmarshalXML( decoder *xml.Decoder, start xml.StartElement ) error {
  s.Defaults()

  for _, attr := range start.Attr {
    switch attr.Name.Local {
      case "rid":
        s.RID = attr.Value

      case "uid":
        s.UID = attr.Value

      case "trainId":
        s.TrainId = attr.Value

      case "ssd":
        s.SSD.Parse( attr.Value )

      case "toc":
        s.Toc = attr.Value

      case "status":
        s.Status = attr.Value

      case "isPassengerSvc":
        s.PassengerService = attr.Value == "true"

      case "isActive":
        s.Active = attr.Value == "true"

      case "deleted":
        s.Deleted = attr.Value == "true"

      case "isCharter":
        s.Charter = attr.Value == "true"
    }
  }

  for {
    token, err := decoder.Token()
    if err != nil {
      return err
    }

    switch tok := token.(type) {
      case xml.StartElement:
        var elem *Location
        switch tok.Name.Local {
          case "OR":
            elem = &Location{ Type: "OR" }

          case "OPOR":
            elem = &Location{ Type: "OPOR" }

          case "IP":
            elem = &Location{ Type: "IP" }

          case "OPIP":
            elem = &Location{ Type: "OPIP" }

          case "PP":
            elem = &Location{ Type: "PP" }

          case "DT":
            elem = &Location{ Type: "DT" }

          case "OPDT":
            elem = &Location{ Type: "OPDT" }

          case "cancelReason":
            if err := decoder.DecodeElement( &s.CancelReason, &tok ); err != nil {
              return err
            }

          default:
            if err := decoder.Skip(); err != nil {
              return err
            }
        }

        if elem != nil {
          if err := decoder.DecodeElement( elem, &tok ); err != nil {
            return err
          }
          s.Locations = append( s.Locations, elem )
        }

      case xml.EndElement:
        s.Sort()
        return nil
    }
  }
}
