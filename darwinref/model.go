// Darwin Reference Data Model

package darwinref

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
