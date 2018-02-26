# util
--
    import "github.com/peter-mount/nre-feeds/util"

Utility types used to store common data like times

## Usage

#### func  PublicTimeEquals

```go
func PublicTimeEquals(a *PublicTime, b *PublicTime) bool
```

#### func  PublicTimeWrite

```go
func PublicTimeWrite(c *codec.BinaryCodec, t *PublicTime)
```
PublicTimeWrite is a workaround for writing null times. If the pointer is null
then a time is written where IsZero()==true

#### func  WorkingTimeEquals

```go
func WorkingTimeEquals(a *WorkingTime, b *WorkingTime) bool
```
WorkingTimeEquals compares equality between two WorkingTimes. Unlike
WorkingTime.Equals() this will return true if both are null, otherwise both must
not be null and equal to be true

#### func  WorkingTimeWrite

```go
func WorkingTimeWrite(c *codec.BinaryCodec, t *WorkingTime)
```
WorkingTimeWrite is a workaround for writing null times. If the pointer is null
then a time is written where IsZero()==true

#### type CircularTimes

```go
type CircularTimes struct {
	// The time for this location.
	// This is calculated as the first value defined below in the following
	// sequence: Wtd, Wta, Wtp, Ptd & Pta.
	Time WorkingTime `json:"time"`
	// Public Scheduled Time of Arrival
	Pta *PublicTime `json:"pta,omitempty"`
	// Public Scheduled Time of Departure
	Ptd *PublicTime `json:"ptd,omitempty"`
	// Working Scheduled Time of Arrival
	Wta *WorkingTime `json:"wta,omitempty"`
	// Working Scheduled Time of Departure
	Wtd *WorkingTime `json:"wtd,omitempty"`
	// Working Scheduled Time of Passing
	Wtp *WorkingTime `json:"wtp,omitempty"`
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

#### func (*CircularTimes) IsPass

```go
func (t *CircularTimes) IsPass() bool
```
IsPass returns true if the instance represents a pass at a station

#### func (*CircularTimes) IsPublic

```go
func (t *CircularTimes) IsPublic() bool
```
IsPublic returns true of the instance contains public times

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

#### type PublicTime

```go
type PublicTime struct {
}
```

Public Timetable time Note: 00:00 is not possible as in CIF that means no-time

#### func  NewPublicTime

```go
func NewPublicTime(s string) *PublicTime
```
NewPublicTime returns a new PublicTime instance from a string of format "HH:MM"

#### func  PublicTimeRead

```go
func PublicTimeRead(c *codec.BinaryCodec) *PublicTime
```
PublicTimeRead is a workaround issue where a custom type cannot be omitempty in
JSON unless it's a nil So instead of using BinaryCodec.Read( v ), we call this &
set the return value in the struct as a pointer.

#### func (*PublicTime) Compare

```go
func (a *PublicTime) Compare(b *PublicTime) bool
```
Compare a PublicTime against another, accounting for crossing midnight. The
rules for handling crossing midnight are: < -6 hours = crossed midnight < 0 back
in time < 18 hours increasing time > 18 hours back in time & crossing midnight

#### func (*PublicTime) Equals

```go
func (a *PublicTime) Equals(b *PublicTime) bool
```

#### func (*PublicTime) Get

```go
func (t *PublicTime) Get() int
```
Get returns the PublicTime in minutes of the day

#### func (*PublicTime) IsZero

```go
func (t *PublicTime) IsZero() bool
```
IsZero returns true if the time is not present

#### func (*PublicTime) MarshalJSON

```go
func (t *PublicTime) MarshalJSON() ([]byte, error)
```
Custom JSON Marshaler. This will write null or the time as "HH:MM"

#### func (*PublicTime) MarshalXMLAttr

```go
func (t *PublicTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error)
```
Custom XML Marshaler.

#### func (*PublicTime) Parse

```go
func (v *PublicTime) Parse(s string)
```

#### func (*PublicTime) Read

```go
func (t *PublicTime) Read(c *codec.BinaryCodec)
```
BinaryCodec reader

#### func (*PublicTime) Set

```go
func (t *PublicTime) Set(v int)
```
Set sets the PublicTime in minutes of the day

#### func (*PublicTime) String

```go
func (t *PublicTime) String() string
```
String returns a PublicTime in HH:MM format or 5 blank spaces if it's not set.

#### func (*PublicTime) UnmarshalJSON

```go
func (t *PublicTime) UnmarshalJSON(b []byte) error
```

#### func (*PublicTime) Write

```go
func (t *PublicTime) Write(c *codec.BinaryCodec)
```
BinaryCodec writer

#### type SSD

```go
type SSD struct {
}
```


#### func (*SSD) Before

```go
func (s *SSD) Before(t time.Time) bool
```
Before is an SSD before a specified time

#### func (*SSD) Equals

```go
func (a *SSD) Equals(b *SSD) bool
```

#### func (*SSD) MarshalJSON

```go
func (t *SSD) MarshalJSON() ([]byte, error)
```
Custom JSON Marshaler.

#### func (*SSD) MarshalXMLAttr

```go
func (t *SSD) MarshalXMLAttr(name xml.Name) (xml.Attr, error)
```
Custom XML Marshaler.

#### func (*SSD) Parse

```go
func (t *SSD) Parse(s string)
```

#### func (*SSD) Read

```go
func (t *SSD) Read(c *codec.BinaryCodec)
```
BinaryCodec reader

#### func (*SSD) String

```go
func (t *SSD) String() string
```
String returns a SSD in "YYYY-MM-DD" format

#### func (*SSD) Time

```go
func (t *SSD) Time() time.Time
```

#### func (*SSD) UnmarshalJSON

```go
func (t *SSD) UnmarshalJSON(b []byte) error
```

#### func (*SSD) Write

```go
func (t *SSD) Write(c *codec.BinaryCodec)
```
BinaryCodec writer

#### type TSTime

```go
type TSTime struct {
	// Estimated Time. For locations with a public activity,
	// this will be based on the "public schedule".
	// For all other activities, it will be based on the "working schedule".
	ET *WorkingTime `json:"et,omitempty" xml:"et,attr,omitempty"`
	// The manually applied lower limit that has been applied to the estimated
	// time at this location. The estimated time will not be set lower than this
	// value, but may be set higher.
	ETMin *WorkingTime `json:"etMin,omitempty" xml:"etmin,attr,omitempty"`
	// Indicates that an unknown delay forecast has been set for the estimated
	// time at this location. Note that this value indicates where a manual
	// unknown delay forecast has been set, whereas it is the "delayed"
	// attribute that indicates that the actual forecast is "unknown delay".
	ETUnknown bool `json:"etUnknown,omitempty" xml:"etUnknown,attr,omitempty"`
	// The estimated time based on the "working schedule".
	// This will only be set for public activities and when it also differs
	// from the estimated time based on the "public schedule".
	WET *WorkingTime `json:"wet,omitempty" xml:"wet,attr,omitempty"`
	// Actual Time
	AT *WorkingTime `json:"at,omitempty" xml:"at,attr,omitempty"`
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

#### func (*TSTime) Compare

```go
func (a *TSTime) Compare(b *TSTime) bool
```
Compare compares two TSTime's. This will use the value returned by TSTime.Time()

#### func (*TSTime) Equals

```go
func (a *TSTime) Equals(b *TSTime) bool
```

#### func (*TSTime) IsSet

```go
func (t *TSTime) IsSet() bool
```
IsSet returns true if at least ET or AT is set

#### func (*TSTime) MarshalJSON

```go
func (t *TSTime) MarshalJSON() ([]byte, error)
```

#### func (*TSTime) Read

```go
func (t *TSTime) Read(c *codec.BinaryCodec)
```

#### func (*TSTime) Time

```go
func (t *TSTime) Time() *WorkingTime
```
Time returns the appropirate time from TSTime to use in displays. This is the
first one set of AT, ET or nil if neither is set.

#### func (*TSTime) UnmarshalXML

```go
func (s *TSTime) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error
```

#### func (*TSTime) Write

```go
func (t *TSTime) Write(c *codec.BinaryCodec)
```

#### type WorkingTime

```go
type WorkingTime struct {
}
```

Working Timetable time. WorkingTime is similar to PublciTime, except we can have
seconds. In the Working Timetable, the seconds can be either 0 or 30.

#### func  NewWorkingTime

```go
func NewWorkingTime(s string) *WorkingTime
```
NewWorkingTime returns a new WorkingTime instance from a string of format
"HH:MM:SS"

#### func  WorkingTimeRead

```go
func WorkingTimeRead(c *codec.BinaryCodec) *WorkingTime
```
WorkingTimeRead is a workaround issue where a custom type cannot be omitempty in
JSON unless it's a nil So instead of using BinaryCodec.Read( v ), we call this &
set the return value in the struct as a pointer.

#### func (*WorkingTime) Compare

```go
func (a *WorkingTime) Compare(b *WorkingTime) bool
```
Compare a WorkingTime against another, accounting for crossing midnight. The
rules for handling crossing midnight are: < -6 hours = crossed midnight < 0 back
in time < 18 hours increasing time > 18 hours back in time & crossing midnight

#### func (*WorkingTime) Equals

```go
func (a *WorkingTime) Equals(b *WorkingTime) bool
```

#### func (*WorkingTime) Get

```go
func (t *WorkingTime) Get() int
```
Get returns the WorkingTime in seconds of the day

#### func (*WorkingTime) IsZero

```go
func (t *WorkingTime) IsZero() bool
```
IsZero returns true if the time is not present

#### func (*WorkingTime) MarshalJSON

```go
func (t *WorkingTime) MarshalJSON() ([]byte, error)
```
Custom JSON Marshaler. This will write null or the time as "HH:MM:SS"

#### func (*WorkingTime) MarshalXMLAttr

```go
func (t *WorkingTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error)
```
Custom XML Marshaler.

#### func (*WorkingTime) Parse

```go
func (v *WorkingTime) Parse(s string)
```

#### func (*WorkingTime) Read

```go
func (t *WorkingTime) Read(c *codec.BinaryCodec)
```
BinaryCodec reader

#### func (*WorkingTime) Set

```go
func (t *WorkingTime) Set(v int)
```
Set sets the WorkingTime in seconds of the day

#### func (*WorkingTime) String

```go
func (t *WorkingTime) String() string
```
String returns a PublicTime in HH:MM:SS format or 8 blank spaces if it's not
set.

#### func (*WorkingTime) UnmarshalJSON

```go
func (t *WorkingTime) UnmarshalJSON(b []byte) error
```

#### func (*WorkingTime) Write

```go
func (t *WorkingTime) Write(c *codec.BinaryCodec)
```
BinaryCodec writer
