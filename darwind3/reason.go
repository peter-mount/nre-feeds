package darwind3

// Type used to represent a cancellation or late running reason
type DisruptionReason struct {
  // A Darwin Reason Code. 0 = none
  Reason  int       `xml:",chardata"`
  // Optional TIPLOC where the reason refers to, e.g. "signalling failure at Cheadle Hulme"
  Tiploc  string    `xml:"tiploc,attr,omitempty"`
  // If true, the tiploc attribute should be interpreted as "near",
  // e.g. "signalling failure near Cheadle Hulme".
  Near    bool      `xml:"near,attr,omitempty"`
}
