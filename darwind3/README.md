# darwind3
--
    import "github.com/peter-mount/darwin/darwind3"

darwind3 handles the real time push port feed

## Usage

#### type CircularTimes

```go
type CircularTimes struct {
	// The time for this location.
	// This is calculated as the first value defined below in the following
	// sequence: Wtd, Wta, Wtp, Ptd & Pta.
	Time darwintimetable.WorkingTime `json:"time"`
	// Public Scheduled Time of Arrival
	Pta *darwintimetable.PublicTime `json:"pta,omitempty"`
	// Public Scheduled Time of Departure
	Ptd *darwintimetable.PublicTime `json:"ptd,omitempty"`
	// Working Scheduled Time of Arrival
	Wta *darwintimetable.WorkingTime `json:"wta,omitempty"`
	// Working Scheduled Time of Departure
	Wtd *darwintimetable.WorkingTime `json:"wtd,omitempty"`
	// Working Scheduled Time of Passing
	Wtp *darwintimetable.WorkingTime `json:"wtp,omitempty"`
}
```

A scheduled time used to distinguish a location on circular routes. Note that
all scheduled time attributes are marked as optional, but at least one must
always be supplied. Only one value is required, and typically this should be the
wtd value. However, for locations that have no wtd, or for clients that deal
exclusively with public times, another value that is valid for the location may
be supplied.

#### func (*CircularTimes) Compare

```go
func (a *CircularTimes) Compare(b *CircularTimes) bool
```
Compare compares two Locations by their times

#### func (*CircularTimes) Equals

```go
func (a *CircularTimes) Equals(b *CircularTimes) bool
```

#### func (*CircularTimes) Read

```go
func (t *CircularTimes) Read(c *codec.BinaryCodec)
```

#### func (*CircularTimes) String

```go
func (l *CircularTimes) String() string
```

#### func (*CircularTimes) UnmarshalXMLAttributes

```go
func (t *CircularTimes) UnmarshalXMLAttributes(start xml.StartElement)
```
UnmarshalXMLAttributes reads from an arbitary start element

#### func (*CircularTimes) UpdateTime

```go
func (l *CircularTimes) UpdateTime()
```
UpdateTime updates the Time field used for sequencing the location. This is the
the first one of these set in the following order: Wtd, Wta, Wtp, Ptd, Pta Note
this value is not persisted as it's a generated value

#### func (*CircularTimes) Write

```go
func (t *CircularTimes) Write(c *codec.BinaryCodec)
```

#### type DarwinD3

```go
type DarwinD3 struct {

	// Optional link to DarwinTimetable for resolving schedules.
	Timetable *darwintimetable.DarwinTimetable
}
```


#### func (*DarwinD3) Close

```go
func (r *DarwinD3) Close()
```
Close the database. If OpenDB() was used to open the db then that db is closed.
If UseDB() was used this simply detaches the DarwinReference from that DB. The
DB is not closed()

#### func (*DarwinD3) GetSchedule

```go
func (r *DarwinD3) GetSchedule(tx *bolt.Tx, rid string) *Schedule
```
GetSchedule retrieves a schedule or nil if not found

#### func (*DarwinD3) OpenDB

```go
func (r *DarwinD3) OpenDB(dbFile string) error
```
OpenDB opens a DarwinReference database.

#### func (*DarwinD3) ProcessUpdate

```go
func (d *DarwinD3) ProcessUpdate(p *Pport, f func(*Transaction) error) error
```

#### func (*DarwinD3) RegisterRest

```go
func (d *DarwinD3) RegisterRest(c *rest.ServerContext)
```

#### func (*DarwinD3) ScheduleHandler

```go
func (d *DarwinD3) ScheduleHandler(r *rest.Rest) error
```

#### func (*DarwinD3) TestHandler

```go
func (d *DarwinD3) TestHandler(r *rest.Rest) error
```
Test handle used to test xml locally via rest

#### func (*DarwinD3) Update

```go
func (r *DarwinD3) Update(f func(*bolt.Tx) error) error
```
Update performs a read write opertation on the database

#### func (*DarwinD3) View

```go
func (r *DarwinD3) View(f func(*bolt.Tx) error) error
```
View performs a readonly operation on the database

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
func (p *DeactivatedSchedule) Process(tx *Transaction) error
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

#### func (*DisruptionReason) Equals

```go
func (a *DisruptionReason) Equals(b *DisruptionReason) bool
```

#### func (*DisruptionReason) Read

```go
func (t *DisruptionReason) Read(c *codec.BinaryCodec)
```

#### func (*DisruptionReason) Write

```go
func (t *DisruptionReason) Write(c *codec.BinaryCodec)
```

#### type Location

```go
type Location struct {
	// Type of location, OR OPOR IP OPIP PP DT or OPDT
	Type string `json:"type"`
	// Tiploc of this location
	Tiploc string `json:"tiploc"`
	// The times for this entry
	Times CircularTimes `json:"timetable"`
	// TIPLOC of False Destination to be used at this location
	FalseDestination string `json:"FalseDestination,omitempty"`
	// Is this service cancelled at this location
	Cancelled bool `json:"cancelled,omitempty"`
	// The Planned data for this location
	// i.e. information planned in advance
	Planned struct {
		// Current Activity Codes
		ActivityType string `json:"activity,omitempty"`
		// Planned Activity Codes (if different to current activities)
		PlannedActivity string `json:"plannedActivity,omitempty"`
		// A delay value that is implied by a change to the service's route.
		// This value has been added to the forecast lateness of the service at
		// the previous schedule location when calculating the expected lateness
		// of arrival at this location.
		RDelay int `json:"rDelay,omitempty"`
	} `json:"planned"`
	// The Forecast data at this location
	// i.e. information that changes in real time
	Forecast struct {
		// Forecast data for the arrival at this location
		Arrival TSTime `json:"arr"`
		// Forecast data for the departure at this location
		Departure TSTime `json:"dep"`
		// Forecast data for the pass of this location
		Pass TSTime `json:"pass"`
		// Current platform number
		Platform Platform `json:"plat"`
		// The service is suppressed at this location.
		Suppressed bool `json:"suppressed,omitempty"`
		// The length of the service at this location on departure
		// (or arrival at destination).
		// The default value of zero indicates that the length is unknown.
		Length int `json:"length,omitempty"`
		// Indicates from which end of the train stock will be detached.
		// The value is set to “true” if stock will be detached from the front of
		// the train at this location. It will be set at each location where stock
		// will be detached from the front.
		// Darwin will not validate that a stock detachment activity code applies
		// at this location.
		DetachFront bool `json:"detachFront,omitempty"`
		// The train order at this location (1, 2 or 3). 0 Means no TrainOrder has been set
		TrainOrder *TrainOrder `json:"trainOrder,omitempty"`
	} `json:"forecast"`
}
```

A location in a schedule. This is formed of the entries from a schedule and is
updated by any incoming Forecasts.

As schedules can be circular (i.e. start and end at the same station) then the
unique key is Tiploc and CircularTimes.Time.

Location's within a Schedule are sorted by CircularTimes.Time accounting for
crossing over midnight.

#### func (*Location) Compare

```go
func (a *Location) Compare(b *Location) bool
```
Compare compares two Locations by their times

#### func (*Location) EqualInSchedule

```go
func (a *Location) EqualInSchedule(b *Location) bool
```
Equals compares two Locations based on their Tiploc & time. This is used when
trying to locate a location that's been updated

#### func (*Location) Equals

```go
func (a *Location) Equals(b *Location) bool
```
Equals compares two Locations in their entirety

#### func (*Location) Read

```go
func (t *Location) Read(c *codec.BinaryCodec)
```

#### func (*Location) String

```go
func (l *Location) String() string
```

#### func (*Location) UnmarshalXML

```go
func (s *Location) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error
```

#### func (*Location) Write

```go
func (t *Location) Write(c *codec.BinaryCodec)
```

#### type Platform

```go
type Platform struct {
	// Defines a platform number
	Platform string `json:"plat,omitempty" xml:",chardata"`
	// True if the platform number is confirmed.
	Confirmed bool `json:"confirmed,omitempty" xml:"conf,attr,omitempty"`
	// Platform number is suppressed and should not be displayed.
	Suppressed bool `json:"suppressed,omitempty" xml:"platsup,attr,omitempty"`
	// Whether a CIS, or Darwin Workstation, has set platform suppression at this location.
	CISSuppressed bool `json:"cisSuppressed,omitempty" xml:"cisPlatsup,attr,omitempty"`
	// The source of the platfom number. P = Planned, A = Automatic, M = Manual.
	// Default is P
	Source string `json:"source,omitempty" xml:"platsrc,attr,omitempty"`
}
```

Platform number with associated flags

#### func (*Platform) Equals

```go
func (a *Platform) Equals(b *Platform) bool
```

#### func (*Platform) Read

```go
func (t *Platform) Read(c *codec.BinaryCodec)
```

#### func (*Platform) Write

```go
func (t *Platform) Write(c *codec.BinaryCodec)
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
func (p *Pport) Process(d3 *DarwinD3) error
```
Process this message

#### func (*Pport) UnmarshalXML

```go
func (s *Pport) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error
```

#### type Processor

```go
type Processor interface {
	Process(*Transaction) error
}
```

Processor interface used by some types used when processing a message and
updating our internal state

#### type Schedule

```go
type Schedule struct {
	RID     string              `json:"rid"`
	UID     string              `json:"uid"`
	TrainId string              `json:"trainId"`
	SSD     darwintimetable.SSD `json:"ssd"`
	Toc     string              `json:"toc"`
	// Default P
	Status string `json:"status"`
	// Default OO
	TrainCat string `json:"trainCat"`
	// Default true
	PassengerService bool `json:"passengerService,omitempty"`
	// Default true
	Active bool `json:"active,omitempty"`
	// Default false
	Deleted bool `json:"deleted,omitempty"`
	// Default false
	Charter bool `json:"charter,omitempty"`
	// Cancel reason
	CancelReason DisruptionReason `json:"cancelReason"`
	// The locations in this schedule
	Locations []*Location `json:"locations"`
	// Usually this is the date we insert into the db but here we use the TS time
	// as returned from darwin
	Date time.Time `json:"date,omitempty" xml:"date,attr,omitempty"`
	// URL to this entity
	Self string `json:"self,omitempty" xml:"self,attr,omitempty"`
}
```

Train Schedule

#### func  ScheduleFromBytes

```go
func ScheduleFromBytes(b []byte) *Schedule
```
ScheduleFromBytes returns a schedule based on a slice or nil if none

#### func (*Schedule) Bytes

```go
func (s *Schedule) Bytes() ([]byte, error)
```
Bytes returns the schedule as an encoded byte slice

#### func (*Schedule) Defaults

```go
func (s *Schedule) Defaults()
```
Defaults sets the default values for a schedule

#### func (*Schedule) Equals

```go
func (a *Schedule) Equals(b *Schedule) bool
```

#### func (*Schedule) Process

```go
func (p *Schedule) Process(tx *Transaction) error
```
Process processes an inbound schedule importing or merging it with the current
Schedule in the database

#### func (*Schedule) Read

```go
func (t *Schedule) Read(c *codec.BinaryCodec)
```

#### func (*Schedule) SetSelf

```go
func (s *Schedule) SetSelf(r *rest.Rest)
```

#### func (*Schedule) Sort

```go
func (s *Schedule) Sort()
```
Sort sorts the locations in a schedule by time order

#### func (*Schedule) String

```go
func (p *Schedule) String() string
```

#### func (*Schedule) UnmarshalXML

```go
func (s *Schedule) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error
```

#### func (*Schedule) Write

```go
func (t *Schedule) Write(c *codec.BinaryCodec)
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
	SSD darwintimetable.SSD `json:"ssd" xml:"ssd,attr"`
	// Indicates whether a train that divides is working with portions in
	// reverse to their normal formation. The value applies to the whole train.
	// Darwin will not validate that a divide association actually exists for
	// this service.
	ReverseFormation bool `json:"isReverseFormation,omitempty" xml:"isReverseFormation,attr,omitempty"`
	//Late running reason for this service.
	// The reason applies to all locations of this service.
	LateReason DisruptionReason `xml:"LateReason"`
	// The locations in this update
	Locations []*Location
}
```

Train Status. Update to the "real time" forecast data for a service.

#### func (*TS) Process

```go
func (p *TS) Process(tx *Transaction) error
```
Process processes an inbound Train Status update, merging it with an existing
schedule in the database

#### func (*TS) UnmarshalXML

```go
func (s *TS) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error
```

#### type TSTime

```go
type TSTime struct {
	// Estimated Time. For locations with a public activity,
	// this will be based on the "public schedule".
	// For all other activities, it will be based on the "working schedule".
	ET string `json:"et,omitempty" xml:"et,attr,omitempty"`
	// The manually applied lower limit that has been applied to the estimated
	// time at this location. The estimated time will not be set lower than this
	// value, but may be set higher.
	ETMin string `json:"etMin,omitempty" xml:"etmin,attr,omitempty"`
	// Indicates that an unknown delay forecast has been set for the estimated
	// time at this location. Note that this value indicates where a manual
	// unknown delay forecast has been set, whereas it is the "delayed"
	// attribute that indicates that the actual forecast is "unknown delay".
	ETUnknown bool `json:"etUnknown,omitempty" xml:"etUnknown,attr,omitempty"`
	// The estimated time based on the "working schedule".
	// This will only be set for public activities and when it also differs
	// from the estimated time based on the "public schedule".
	WET string `json:"wet,omitempty" xml:"wet,attr,omitempty"`
	// Actual Time
	AT string `json:"at,omitempty" xml:"at,attr,omitempty"`
	// If true, indicates that an actual time ("at") value has just been removed
	// and replaced by an estimated time ("et").
	// Note that this attribute will only be set to "true" once, when the actual
	// time is removed, and will not be set in any snapshot.
	ATRemoved bool `json:"atRemoved,omitempty" xml:"atRemoved,attr,omitempty"`
	// Indicates that this estimated time is a forecast of "unknown delay".
	// Displayed  as "Delayed" in LDB.
	// Note that this value indicates that this forecast is "unknown delay",
	// whereas it is the "etUnknown" attribute that indicates where the manual
	// unknown delay forecast has been set.
	Delayed bool `json:"delayed,omitempty" xml:"delayed,attr,omitempty"`
	// The source of the forecast or actual time.
	Src string `json:"src,omitempty" xml:"src,attr,omitempty"`
	// The RTTI CIS code of the CIS instance if the src is a CIS.
	SrcInst string `json:"srcInst,omitempty" xml:"srcInst,attr,omitempty"`
}
```

Type describing time-based forecast attributes for a TS arrival/departure/pass

#### func (*TSTime) Equals

```go
func (a *TSTime) Equals(b *TSTime) bool
```

#### func (*TSTime) Read

```go
func (t *TSTime) Read(c *codec.BinaryCodec)
```

#### func (*TSTime) Write

```go
func (t *TSTime) Write(c *codec.BinaryCodec)
```

#### type TrainOrder

```go
type TrainOrder struct {
	Order int `json:"order" xml:"order,attr"`
	// The platform number where the train order applies
	Platform string `json:"plat,omitempty" xml:"plat,attr,omitempty"`
}
```

Defines the expected Train order at a platform

#### type Transaction

```go
type Transaction struct {
}
```


#### func (*Transaction) DeleteSchedule

```go
func (tx *Transaction) DeleteSchedule(rid string) error
```
DeleteSchedule deletes a schedule from the database

#### func (*Transaction) GetSchedule

```go
func (tx *Transaction) GetSchedule(rid string) *Schedule
```
GetSchedule retrieves a schedule or nil if not found

#### func (*Transaction) PutSchedule

```go
func (tx *Transaction) PutSchedule(s *Schedule) error
```
PutSchedule persists a Schedule

#### func (*Transaction) ResolveSchedule

```go
func (d *Transaction) ResolveSchedule(rid string) *Schedule
```
ResolveSchedule attempts to retrieve a schedule from the timetable. If
DarwinD3.Timetable is not set then this always returns nil

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
func (p *UR) Process(tx *Transaction) error
```
Process this message

#### func (*UR) UnmarshalXML

```go
func (s *UR) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error
```
