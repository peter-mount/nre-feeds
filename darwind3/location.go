package darwind3

import (
  "encoding/xml"
  "strconv"
)

type Location struct {
  // Type of location, OR OPOR IP OPIP PP DT or OPDT
  Type              string
  // Tiploc of this location
  Tiploc            string
  // TIPLOC of False Destination to be used at this location
  FalseDestination  string
  // Current Activity Codes
  ActivityType      string
  // Planned Activity Codes (if different to current activities)
  PlannedActivity   string
  // Is this service cancelled at this location
  Cancelled         bool
  // Public Scheduled Time of Arrival
  Pta               string
  // Public Scheduled Time of Departure
  Ptd               string
  // Working Scheduled Time of Arrival
  Wta               string
  // Working Scheduled Time of Departure
  Wtd               string
  // Working Scheduled Time of Passing
  Wtp               string
  // A delay value that is implied by a change to the service's route.
  // This value has been added to the forecast lateness of the service at
  // the previous schedule location when calculating the expected lateness
  // of arrival at this location.
  RDelay            int
}

func (s *Location) UnmarshalXML( decoder *xml.Decoder, start xml.StartElement ) error {

  for _, attr := range start.Attr {
    switch attr.Name.Local {
      case "tpl":
        s.Tiploc = attr.Value

      case "act":
        s.ActivityType = attr.Value

      case "planAct":
        s.PlannedActivity = attr.Value

      case "Cancelled":
        s.Cancelled = attr.Value == "true"

      case "pta":
        s.Pta = attr.Value

      case "ptd":
        s.Ptd = attr.Value

      case "wta":
        s.Wta = attr.Value

      case "wtd":
        s.Wtd = attr.Value

      case "wtp":
        s.Wtp = attr.Value

      case "fd":
        s.FalseDestination = attr.Value

      case "rdelay":
        if v, err := strconv.Atoi( attr.Value ); err != nil {
          return err
        } else {
          s.RDelay = v
        }
    }
  }

  return decoder.Skip()
}
