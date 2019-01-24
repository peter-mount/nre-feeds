package darwind3

import (
  "encoding/xml"
  "time"
)

// Defines the expected Train order at a platform
type TrainOrder struct {
  Order     int       `json:"order" xml:"order,attr"`
  // The platform number where the train order applies
  Platform  string    `json:"plat,omitempty" xml:"plat,attr,omitempty"`
  // This is the TS time from Darwin so we keep a copy of when this struct
  // was sent to us
  Date              time.Time             `json:"date,omitempty"`
}

// The trainOrder as received from darwin
type trainOrderWrapper struct {
  // The tiploc where the train order applies
  Tiploc    string
  // The CRS code of the station where the train order applies
  CRS       string
  // The platform number where the train order applies
  Platform  string
  // The Train orders to set
  Set      *trainOrderData
  // Clear the current train order
  Clear     bool
}

func (s *trainOrderWrapper) UnmarshalXML( decoder *xml.Decoder, start xml.StartElement ) error {

  for _, attr := range start.Attr {
    switch attr.Name.Local {
      case "tiploc":
        s.Tiploc = attr.Value

      case "crs":
        s.CRS = attr.Value

      case "platform":
        s.Platform = attr.Value
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
          case "clear":
            s.Clear = true
            if err := decoder.Skip(); err != nil {
              return err
            }

          case "set":
            s.Set = &trainOrderData{}
            if err := decoder.DecodeElement( s.Set, &tok ); err != nil {
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
