package darwind3

import (
  "encoding/xml"
)

// Update Response
type UR struct {
  XMLName             xml.Name            `json:"-" xml:"uR"`
  UpdateOrigin        string              `xml:"updateOrigin,attr,omitempty"`
  RequestSource       string              `xml:"requestSource,attr,omitempty"`
  RequestId           string              `xml:"requestId,attr,omitempty"`
  Actions          []Processor
}

func (s *UR) UnmarshalXML( decoder *xml.Decoder, start xml.StartElement ) error {

  for _, attr := range start.Attr {
    switch attr.Name.Local {
      case "updateOrigin":
        s.UpdateOrigin = attr.Value
      case "requestSource":
        s.RequestSource = attr.Value
      case "requestId":
        s.RequestId = attr.Value
    }
  }

  for {
    token, err := decoder.Token()
    if err != nil {
      return err
    }

    switch tok := token.(type) {
      case xml.StartElement:
        var elem Processor
        switch tok.Name.Local {
          case "schedule":
            elem = &Schedule{}

          case "deactivated":
            elem = &DeactivatedSchedule{}

          case "TS":
            elem = &TS{}

          default:
            if err := decoder.Skip(); err != nil {
              return err
            }
        }

        if elem != nil {
          if err := decoder.DecodeElement( elem, &tok ); err != nil {
            return err
          }
          s.Actions = append( s.Actions, elem )
        }

      case xml.EndElement:
        return nil
    }
  }
}

// Process this message
func (p *UR) Process( d3 *DarwinD3, r *Pport ) error {

  if len( p.Actions ) > 0 {
    for _, s := range p.Actions {
      if err:= s.Process( d3, r ); err != nil {
        return err
      }
    }
  }

  return nil
}
