package darwind3

import (
  "encoding/xml"
)

// Train Status. Update to the "real time" forecast data for a service.
type TS struct {
  XMLName           xml.Name  `json:"-" xml:"TS"`
  // RTTI unique Train Identifier
  RID               string    `json:"rid" xml:"rid,attr"`
  // Train UID
  UID               string    `json:"uid" xml:"uid,attr"`
  // Scheduled Start Date
  SSD               string    `json:"ssd" xml:"ssd,attr"`
  // Indicates whether a train that divides is working with portions in
  // reverse to their normal formation. The value applies to the whole train.
  // Darwin will not validate that a divide association actually exists for this service.
  ReverseFormation  bool      `json:"isReverseFormation,omitempty" xml:"isReverseFormation,attr,omitempty"`
  LateReason        DisruptionReason  `xml:"LateReason"`
}
