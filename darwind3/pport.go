package darwind3

import (
  "encoding/xml"
  "fmt"
  "time"
)

// The Pport element
type Pport struct {
  XMLName     xml.Name  `json:"-" xml:"Pport"`
  TS          time.Time `json:"ts" xml:"ts,attr"`
  Version     string    `json:"version" xml:"version,attr"`
  UR         *UR        `json:"uR" xml:"uR"`
}

// Process this message
func (p *Pport) Process( d3 *DarwinD3 ) error {

  if p.UR != nil {
    if err := p.UR.Process( d3, p ); err != nil {
      return err
    }
  }

  return nil
}

func (p *Pport) String() string {
  s := fmt.Sprintf("Pport ts=%s version=%s\n", p.TS, p.Version )

  if p.UR != nil {
    s += p.UR.String()
  }

  return s
}
