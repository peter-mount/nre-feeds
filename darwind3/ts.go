package darwind3

import (
  "encoding/xml"
  "github.com/peter-mount/nre-feeds/util"
)

// Train Status. Update to the "real time" forecast data for a service.
type TS struct {
  XMLName           xml.Name  `json:"-" xml:"TS"`
  // RTTI unique Train Identifier
  RID               string    `json:"rid" xml:"rid,attr"`
  // Train UID
  UID               string    `json:"uid" xml:"uid,attr"`
  // Scheduled Start Date
  SSD               util.SSD  `json:"ssd" xml:"ssd,attr"`
  // Indicates whether a train that divides is working with portions in
  // reverse to their normal formation. The value applies to the whole train.
  // Darwin will not validate that a divide association actually exists for
  // this service.
  ReverseFormation  bool      `json:"isReverseFormation,omitempty" xml:"isReverseFormation,attr,omitempty"`
  //Late running reason for this service.
  // The reason applies to all locations of this service.
  LateReason        DisruptionReason  `xml:"LateReason"`
  // The locations in this update
  Locations       []*Location
}

func (s *TS) UnmarshalXML( decoder *xml.Decoder, start xml.StartElement ) error {
  for _, attr := range start.Attr {
    switch attr.Name.Local {
      case "rid":
        s.RID = attr.Value

      case "uid":
        s.UID = attr.Value

      case "ssd":
        s.SSD.Parse( attr.Value )

      case "isReverseFormation":
        s.ReverseFormation = attr.Value == "true"
    }
  }

  for {
    token, err := decoder.Token()
    if err != nil {
      return err
    }

    switch tok := token.(type) {
      case xml.StartElement:
        switch tok.Name.Local {
          case "Location":
            // TS for Train Status.
            // If you see this in any output it means that we have received a
            // forecast with no original schedule entry
            loc := &Location{ Type: "TS" }
            if err := decoder.DecodeElement( loc, &tok ); err != nil {
              return err
            }
            s.Locations = append( s.Locations, loc )

          case "LateReason":
            if err := decoder.DecodeElement( &s.LateReason, &tok ); err != nil {
              return err
            }

          default:
            if err := decoder.Skip(); err != nil {
              return err
            }
        }

      case xml.EndElement:
        return nil
    }
  }
}
