package darwind3

import (
  "encoding/xml"
  "fmt"
)

// Update Response
type UR struct {
  XMLName             xml.Name            `json:"-" xml:"uR"`
  UpdateOrigin        string              `xml:"updateOrigin,attr,omitempty"`
  RequestSource       string              `xml:"requestSource,attr,omitempty"`
  RequestId           string              `xml:"requestId,attr,omitempty"`
  // extension tns:DataResponse
  Schedule         []*Schedule            `xml:"schedule"`
  Deactivated      []*DeactivatedSchedule `xml:"deactivated"`
  /*
  Association      []*Association         `xml:"association"`
  TS     []*Schedule                      `xml:"TS"`
  StationMessage   []*StationMessage      `xml:"OW"`
  TrainAlert       []*TrainAlert          `xml:"trainAlert"`
  TrainOrder       []*TrainOrder          `xml:"trainOrder"`
  TrackingID       []*TrackingID          `xml:"trackingID"`
  Alarm            []*Alarm               `xml:"alarm"`
  */
}

// Process this message
func (p *UR) Process( d3 *DarwinD3, r *Pport ) error {

  if len( p.Schedule ) > 0 {
    for _, s := range p.Schedule {
      if err:= s.Process( d3, r ); err != nil {
        return err
      }
    }
  }

  return nil
}

func (p *UR) String() string {
  s := fmt.Sprintf("UR updateOrigin=%s requestSource=%s requestId=%s\n", p.UpdateOrigin, p.RequestSource, p.RequestId )

  if len( p.Schedule ) > 0 {
    for _, e := range p.Schedule {
      s += e.String()
    }
  }

  return s
}
