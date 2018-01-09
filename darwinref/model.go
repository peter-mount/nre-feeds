// Darwin Reference Data Model

package darwinref

// Processed reference format
type DarwinReference struct {
  timetableId         string
  // Map of all locations by tiploc
  tiploc              map[string]*Location
  // Map of all locations by CRS/3Alpha code
  crs                 map[string][]*Location
  // Map of Toc (Operator) codes
  toc                 map[string]*Toc
  // Reasons for a train being late
  lateRunningReasons  map[int]string
  // Reasons for a train being cancelled at a location
  cancellationReasons map[int]string
  // CIS source
  cisSource           map[string]string
  // via texts, map of at+","+ dest then array of possibilities
  via                 map[string][]*Via
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
  Text    string        `xml:"viatext,attr"`
}

type CISSource struct {
  Code    string        `xml:"code,attr"`
  Name    string        `xml:"name,attr"`
}
