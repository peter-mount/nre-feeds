package darwind3

// Platform number with associated flags
type Platform struct {
  // Defines a platform number
  Platform          string    `json:"plat,omitempty" xml:",chardata"`
  // True if the platform number is confirmed.
  Confirmed         bool      `json:"confirmed,omitempty" xml:"conf,attr,omitempty"`
  // Platform number is suppressed and should not be displayed.
  Suppressed        bool      `json:"suppressed,omitempty" xml:"platsup,attr,omitempty"`
  // Whether a CIS, or Darwin Workstation, has set platform suppression at this location.
  CISSuppressed     bool      `json:"cisSuppressed,omitempty" xml:"cisPlatsup,attr,omitempty"`
  // The source of the platfom number. P = Planned, A = Automatic, M = Manual.
  // Default is P
  Source            string    `json:"source,omitempty" xml:"platsrc,attr,omitempty"`
}

func (a *Platform ) Equals( b *Platform ) bool {
  return b != nil &&
         a.Platform == b.Platform &&
         a.Confirmed == b.Confirmed &&
         a.Suppressed == b.Suppressed &&
         a.CISSuppressed == b.CISSuppressed &&
         a.Source == b.Source
}
