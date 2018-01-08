// Darwin Reference Data Model

package darwinref

// The reference xml wrapper
type PportTimetableRef struct {
  Locations           []*Location             `xml:"LocationRef"`
  Toc                 []*Toc                  `xml:"TocRef"`
  LateRunningReasons    *LateRunningReasons   `xml:"LateRunningReasons"`
  CancellationReasons   *CancellationReasons  `xml:"CancellationReasons"`
  CISSource           []*CISSource            `xml:"CISSource"`
}

// Defines a location, i.e. Station or passing point
type Location struct {
  Tiploc      string            `xml:"tpl,attr"`
  Crs         string            `xml:"crs,attr"`
  Toc         string            `xml:"toc,attr"`
  Name        string            `xml:"locname,attr"`
}

// A rail operator
type Toc struct {
  Toc         string            `xml:"toc,attr"`
  Name        string            `xml:"tocname,attr"`
  Url         string            `xml:"url,attr"`
}

// Late Running Reasons
type LateRunningReasons struct {
  Reason    []*Reason            `xml:"Reason"`
}

// Cancellation Reasons
type CancellationReasons struct {
  Reason    []Reason            `xml:"Reason"`
}

// A reason, shared by LateRunningReasons and CancellationReasons
type Reason struct {
  Code        int               `xml:"code,attr"`
  Text        string            `xml:"reasontext,attr"`
}

// Via text
type Via struct {
  At      string        `xml:"at,attr"`
  Dest    string        `xml:"dest,attr"`
  Loc1    string        `xml:"loc1,attr"`
  Loc2    string        `xml:"loc2,attr"`
  Text    string        `xml:"viatext,attr`
}

type CISSource struct {
  Code    string        `xml:"code,attr"`
  Name    string        `xml:"name,attr"`
}
