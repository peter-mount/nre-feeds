package darwind3

import (
  "encoding/xml"
  "time"
)

// The Pport element
type Pport struct {
  XMLName     xml.Name  `json:"-" xml:"Pport"`
  TS          time.Time `json:"ts" xml:"ts,attr"`
  Version     string    `json:"version" xml:"version,attr"`
  Actions   []Processor
}

func (s *Pport) UnmarshalXML( decoder *xml.Decoder, start xml.StartElement ) error {

  for _, attr := range start.Attr {
    switch attr.Name.Local {
      case "ts":
        if t, err := time.Parse( time.RFC3339Nano, attr.Value ); err != nil {
          return err
        } else {
          s.TS = t
        }

      case "version":
        s.Version = attr.Value
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
          case "uR":
            elem = &UR{}

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
func (p *Pport) Process( d3 *DarwinD3 ) error {

  if len( p.Actions ) > 0 {
    for _, s := range p.Actions {
      // Use a write transaction for each action
      if err := d3.ProcessUpdate( p, func( tx *Transaction ) error {
        return s.Process( tx )
      }); err != nil {
        return err
      }
    }
  }

  return nil
}
