# darwind3
--
    import "github.com/peter-mount/darwin/darwind3"

darwind3 handles the real time push port feed

## Usage

#### type DarwinD3

```go
type DarwinD3 struct {
}
```


#### func (*DarwinD3) SetupRest

```go
func (d *DarwinD3) SetupRest(c *rest.ServerContext)
```

#### func (*DarwinD3) TestHandler

```go
func (d *DarwinD3) TestHandler(r *rest.Rest) error
```
Test handle used to test xml locally via rest

#### type DeactivatedSchedule

```go
type DeactivatedSchedule struct {
	XMLName xml.Name `json:"-" xml:"deactivated"`
	RID     string   `xml:"rid,attr"`
}
```

Notification that a Train Schedule is now deactivated in Darwin.

#### func (*DeactivatedSchedule) Process

```go
func (p *DeactivatedSchedule) Process(d3 *DarwinD3, r *Pport) error
```
Processor interface

#### type DisruptionReason

```go
type DisruptionReason struct {
	// A Darwin Reason Code. 0 = none
	Reason int `xml:",chardata"`
	// Optional TIPLOC where the reason refers to, e.g. "signalling failure at Cheadle Hulme"
	Tiploc string `xml:"tiploc,attr,omitempty"`
	// If true, the tiploc attribute should be interpreted as "near",
	// e.g. "signalling failure near Cheadle Hulme".
	Near bool `xml:"near,attr,omitempty"`
}
```

Type used to represent a cancellation or late running reason

#### type Location

```go
type Location struct {
	// Type of location, OR OPOR IP OPIP PP DT or OPDT
	Type string
	// Tiploc of this location
	Tiploc string
	// TIPLOC of False Destination to be used at this location
	FalseDestination string
	// Current Activity Codes
	ActivityType string
	// Planned Activity Codes (if different to current activities)
	PlannedActivity string
	// Is this service cancelled at this location
	Cancelled bool
	// Public Scheduled Time of Arrival
	Pta string
	// Public Scheduled Time of Departure
	Ptd string
	// Working Scheduled Time of Arrival
	Wta string
	// Working Scheduled Time of Departure
	Wtd string
	// Working Scheduled Time of Passing
	Wtp string
	// A delay value that is implied by a change to the service's route.
	// This value has been added to the forecast lateness of the service at
	// the previous schedule location when calculating the expected lateness
	// of arrival at this location.
	RDelay int
}
```


#### func (*Location) UnmarshalXML

```go
func (s *Location) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error
```

#### type Pport

```go
type Pport struct {
	XMLName xml.Name  `json:"-" xml:"Pport"`
	TS      time.Time `json:"ts" xml:"ts,attr"`
	Version string    `json:"version" xml:"version,attr"`
	Actions []Processor
}
```

The Pport element

#### func (*Pport) Process

```go
func (p *Pport) Process(d3 *DarwinD3, r *Pport) error
```
Process this message

#### func (*Pport) UnmarshalXML

```go
func (s *Pport) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error
```

#### type Processor

```go
type Processor interface {
	Process(*DarwinD3, *Pport) error
}
```

Process a messafe

#### type Schedule

```go
type Schedule struct {
	RID     string
	UID     string
	TrainId string
	SSD     string
	Toc     string
	// Default P
	Status string
	// Default OO
	TrainCat string
	// Default true
	PassengerService bool
	// Default true
	Active bool
	// Default false
	Deleted bool
	// Default false
	Charter bool
	// Cancel reason
	CancelReason DisruptionReason
	// The locations in this schedule
	Locations []*Location
}
```

Train Schedule

#### func (*Schedule) Process

```go
func (p *Schedule) Process(d3 *DarwinD3, r *Pport) error
```
Processor interface

#### func (*Schedule) UnmarshalXML

```go
func (s *Schedule) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error
```

#### type TS

```go
type TS struct {
	XMLName xml.Name `json:"-" xml:"TS"`
	// RTTI unique Train Identifier
	RID string `json:"rid" xml:"rid,attr"`
	// Train UID
	UID string `json:"uid" xml:"uid,attr"`
	// Scheduled Start Date
	SSD string `json:"ssd" xml:"ssd,attr"`
	// Indicates whether a train that divides is working with portions in
	// reverse to their normal formation. The value applies to the whole train.
	// Darwin will not validate that a divide association actually exists for this service.
	ReverseFormation bool             `json:"isReverseFormation,omitempty" xml:"isReverseFormation,attr,omitempty"`
	LateReason       DisruptionReason `xml:"LateReason"`
}
```

Train Status. Update to the "real time" forecast data for a service.

#### func (*TS) Process

```go
func (p *TS) Process(d3 *DarwinD3, r *Pport) error
```
Processor interface

#### type UR

```go
type UR struct {
	XMLName       xml.Name `json:"-" xml:"uR"`
	UpdateOrigin  string   `xml:"updateOrigin,attr,omitempty"`
	RequestSource string   `xml:"requestSource,attr,omitempty"`
	RequestId     string   `xml:"requestId,attr,omitempty"`
	Actions       []Processor
}
```

Update Response

#### func (*UR) Process

```go
func (p *UR) Process(d3 *DarwinD3, r *Pport) error
```
Process this message

#### func (*UR) UnmarshalXML

```go
func (s *UR) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error
```
