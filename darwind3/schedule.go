package darwind3

import (
  "darwintimetable"
  "encoding/xml"
  "fmt"
  "github.com/peter-mount/golib/codec"
  "sort"
)

// Train Schedule
type Schedule struct {
  RID               string
  UID               string
  TrainId           string
  SSD               darwintimetable.SSD
  Toc               string
  // Default P
  Status            string
  // Default OO
  TrainCat          string
  // Default true
  PassengerService  bool
  // Default true
  Active            bool
  // Default false
  Deleted           bool
  // Default false
  Charter           bool
  // Cancel reason
  CancelReason      DisruptionReason
  // The locations in this schedule
  Locations       []*Location
}

// Sort sorts the locations in a schedule by time order
func (s *Schedule) Sort() {
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
       len( a.Locations ) == len( b.Locations )

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
  codec.NewBinaryCodecFrom( b ).Read( sched )
  if sched.RID == "" {
    return nil
  }
  return sched
}

// Bytes returns the schedule as an encoded byte slice
func (s *Schedule) Bytes() ( []byte, error ) {
  encoder := codec.NewBinaryCodec()
  encoder.Write( s )
  if encoder.Error() != nil {
    return nil, encoder.Error()
  }
  return encoder.Bytes(), nil
}

func (t *Schedule) Write( c *codec.BinaryCodec ) {
  c.WriteString( t.RID ).
    WriteString( t.UID ).
    WriteString( t.TrainId )
  c.Write( &t.SSD )
  c.WriteString( t.Toc ).
    WriteString( t.Status ).
    WriteString( t.TrainCat ).
    WriteBool( t.PassengerService ).
    WriteBool( t.Active ).
    WriteBool( t.Deleted ).
    WriteBool( t.Charter ).
    Write( &t.CancelReason )

  c.WriteInt16( int16(len( t.Locations )) )
  for _, l := range t.Locations {
    c.Write( l )
  }
}

func (t *Schedule) Read( c *codec.BinaryCodec ) {
  c.ReadString( &t.RID ).
    ReadString( &t.UID ).
    ReadString( &t.TrainId )
  c.Read( &t.SSD )
  c.ReadString( &t.Toc ).
    ReadString( &t.Status ).
    ReadString( &t.TrainCat ).
    ReadBool( &t.PassengerService ).
    ReadBool( &t.Active ).
    ReadBool( &t.Deleted ).
    ReadBool( &t.Charter ).
    Read( &t.CancelReason )

  var n int16
  c.ReadInt16( &n )
  for i := 0; i < int(n); i++ {
    l := &Location{}
    c.Read( l )
    t.Locations = append( t.Locations, l )
  }
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

func (p *Schedule) String() string {
  s := fmt.Sprintf(
    "Schedule rid=%s uid=%s trainId=%s ssd=%s toc=%s status=%s trainCat=%s isPax=%v active=%v deleted=%v charter=%v cancelReason=%v locs=%d\n",
    p.RID,
    p.UID,
    p.TrainId,
    p.SSD.String(),
    p.Toc,
    p.Status,
    p.TrainCat,
    p.PassengerService,
    p.Active,
    p.Deleted,
    p.Charter,
    p.CancelReason,
    len( p.Locations ) )
  for i, l := range p.Locations {
    s += fmt.Sprintf( "%02d %s\n", i, l.String() )
  }
  return s
}
