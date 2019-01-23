package darwind3

// Type used to represent a cancellation or late running reason
type DisruptionReason struct {
  // A Darwin Reason Code. 0 = none
  Reason  int       `json:"reason" xml:",chardata"`
  // Optional TIPLOC where the reason refers to, e.g. "signalling failure at Cheadle Hulme"
  Tiploc  string    `json:"tiploc,omitempty" xml:"tiploc,attr,omitempty"`
  // If true, the tiploc attribute should be interpreted as "near",
  // e.g. "signalling failure near Cheadle Hulme".
  Near    bool      `json:"near,omitempty" xml:"near,attr,omitempty"`
}

func (a *DisruptionReason) Equals( b *DisruptionReason ) bool {
  return b != nil &&
         a.Reason == b.Reason &&
         a.Tiploc == b.Tiploc &&
         a.Near == b.Near
}
