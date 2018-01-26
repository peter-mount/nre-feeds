# darwintimetable
--
    import "github.com/peter-mount/darwin/darwintimetable"

Reference timetable


Reference timetable

Reference timetable


Reference timetable

Reference timetable

Reference timetable

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

#### type AssocService

```go
type AssocService struct {
	RID string `json:"rid" xml:"rid,attr"`
	Wta string `json:"wta" xml:"wta,attr"`
	Wtd string `json:"wtd" xml:"wtd,attr"`
	Wtp string `json:"wtp" xml:"wtp,attr"`
	Pta string `json:"pta" xml:"pta,attr"`
	Ptd string `json:"ptd" xml:"ptd,attr"`
}
```


#### func (*AssocService) Read

```go
func (t *AssocService) Read(c *codec.BinaryCodec)
```

#### func (*AssocService) Write

```go
func (t *AssocService) Write(c *codec.BinaryCodec)
```

#### type Association

```go
type Association struct {
	XMLName   xml.Name     `json:"-" xml:"Association"`
	Main      AssocService `json:"main" xml:"main"`
	Assoc     AssocService `json:"assoc" xml:"assoc"`
	Tiploc    string       `json:"tiploc" xml:"tiploc,attr"`
	Category  string       `json:"category" xml:"category,attr"`
	Cancelled bool         `json:"cancelled" xml:"isCancelled,attr"`
	Deleted   bool         `json:"deleted" xml:"isDeleted,attr"`
	// Date entry was inserted into the database
	Date time.Time `json:"date" xml:"date,attr"`
	// URL to this entity
	Self string `json:"self" xml:"self,attr,omitempty"`
}
```


#### func (*Association) Equals

```go
func (a *Association) Equals(b *Association) bool
```

#### func (*Association) Read

```go
func (t *Association) Read(c *codec.BinaryCodec)
```

#### func (*Association) Write

```go
func (t *Association) Write(c *codec.BinaryCodec)
```

#### type DT

```go
type DT struct {
	// SchedLocAttributes
	Tiploc    string `xml:"tpl,attr"`
	Act       string `xml:"act,attr"`
	PlanAct   string `xml:"planAct,attr"`
	Cancelled bool   `xml:"can,attr"`
	Platform  string `xml:"plat,attr"`
	// CallPtAttributes
	Pta string `xml:"pta,attr"`
	Ptd string `xml:"ptd,attr"`
	// Working times
	Wta string `xml:"wta,attr"`
	Wtd string `xml:"wtd,attr"`
	// Delay implied by a change to the services route
	RDelay string `xml:"rdelay,attr"`
}
```


#### func (*DT) Location

```go
func (s *DT) Location() *Location
```

#### type DarwinTimetable

```go
type DarwinTimetable struct {
}
```


#### func (*DarwinTimetable) Close

```go
func (r *DarwinTimetable) Close()
```
Close the database. If OpenDB() was used to open the db then that db is closed.
If UseDB() was used this simply detaches the DarwinReference from that DB. The
DB is not closed()

#### func (*DarwinTimetable) GetJourney

```go
func (r *DarwinTimetable) GetJourney(tx *bolt.Tx, rid string) (*Journey, bool)
```
GetJourney returns details of a Journey

#### func (*DarwinTimetable) GetJourneyBucket

```go
func (r *DarwinTimetable) GetJourneyBucket(bucket *bolt.Bucket, rid string) (*Journey, bool)
```

#### func (*DarwinTimetable) ImportHandler

```go
func (dt *DarwinTimetable) ImportHandler(r *rest.Rest) error
```

#### func (*DarwinTimetable) JourneyHandler

```go
func (dt *DarwinTimetable) JourneyHandler(r *rest.Rest) error
```

#### func (*DarwinTimetable) OpenDB

```go
func (r *DarwinTimetable) OpenDB(dbFile string) error
```
OpenDB opens a DarwinReference database.

#### func (*DarwinTimetable) PruneSchedules

```go
func (t *DarwinTimetable) PruneSchedules() (int, error)
```
PruneSchedules prunes all expired schedules

#### func (*DarwinTimetable) PruneSchedulesHandler

```go
func (dt *DarwinTimetable) PruneSchedulesHandler(r *rest.Rest) error
```

#### func (*DarwinTimetable) Read

```go
func (t *DarwinTimetable) Read(c *codec.BinaryCodec)
```

#### func (*DarwinTimetable) RegisterRest

```go
func (r *DarwinTimetable) RegisterRest(c *rest.ServerContext)
```
RegisterRest registers the rest endpoints into a ServerContext

#### func (*DarwinTimetable) ScheduleCleanup

```go
func (t *DarwinTimetable) ScheduleCleanup(c *cron.Cron)
```

#### func (*DarwinTimetable) TimetableId

```go
func (r *DarwinTimetable) TimetableId() string
```
Return's the timetableId for this reference dataset

#### func (*DarwinTimetable) UnmarshalXML

```go
func (t *DarwinTimetable) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error
```

#### func (*DarwinTimetable) Update

```go
func (r *DarwinTimetable) Update(f func(*bolt.Tx) error) error
```
Update performs a read write opertation on the database

#### func (*DarwinTimetable) UseDB

```go
func (r *DarwinTimetable) UseDB(db *bolt.DB) error
```
UseDB Allows an already open database to be used with DarwinReference.

#### func (*DarwinTimetable) View

```go
func (r *DarwinTimetable) View(f func(*bolt.Tx) error) error
```
View performs a readonly operation on the database

#### func (*DarwinTimetable) Write

```go
func (t *DarwinTimetable) Write(c *codec.BinaryCodec)
```

#### type IP

```go
type IP struct {
	// SchedLocAttributes
	Tiploc    string `xml:"tpl,attr"`
	Act       string `xml:"act,attr"`
	PlanAct   string `xml:"planAct,attr"`
	Cancelled bool   `xml:"can,attr"`
	Platform  string `xml:"plat,attr"`
	// CallPtAttributes
	Pta string `xml:"pta,attr"`
	Ptd string `xml:"ptd,attr"`
	// Working times
	Wta string `xml:"wta,attr"`
	Wtd string `xml:"wtd,attr"`
	// Delay implied by a change to the services route
	RDelay string `xml:"rdelay,attr"`
	// False destination to be used at this location
	FalseDest string `xml:"fd,attr"`
}
```


#### func (*IP) Location

```go
func (s *IP) Location() *Location
```

#### type Journey

```go
type Journey struct {
	XMLName   xml.Name `json:"-" xml:"Journey"`
	RID       string   `json:"rid" xml:"rid,attr"`
	UID       string   `json:"uid" xml:"uid,attr"`
	TrainID   string   `json:"trainId" xml:"trainId"`
	SSD       SSD      `json:"ssd" xml:"ssd,attr"`
	Toc       string   `json:"toc" xml:"toc,attr"`
	TrainCat  string   `json:"trainCat" xml:"trainCat,attr"`
	Passenger bool     `json:"isPassengerSvc" xml:"isPassengerSvc,attr"`
	// The schedule
	Schedule     []*Location `json:"locations" xml:location`
	CancelReason int         `json:"cancelReason" xml:"cancelReason,attr"`
	// Associations
	//Associations  []*Association  `xml:"-"`
	// Date entry was inserted into the database
	Date time.Time `json:"date" xml:"date,attr"`
	// URL to this entity
	Self string `json:"self" xml:"self,attr,omitempty"`
}
```


#### func (*Journey) Equals

```go
func (a *Journey) Equals(b *Journey) bool
```

#### func (*Journey) Read

```go
func (t *Journey) Read(c *codec.BinaryCodec)
```

#### func (*Journey) SetSelf

```go
func (t *Journey) SetSelf(r *rest.Rest)
```

#### func (*Journey) UnmarshalXML

```go
func (j *Journey) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error
```

#### func (*Journey) Write

```go
func (t *Journey) Write(c *codec.BinaryCodec)
```

#### type Location

```go
type Location struct {
	XMLName   xml.Name `json:"-" xml:"location"`
	Type      string   `json:"type" xml:"type,attr"`
	Tiploc    string   `json:"tpl" xml:"tpl,attr"`
	Act       string   `json:"act,omitempty" xml:"act,attr,omitempty"`
	PlanAct   string   `json:"planAct,omitempty" xml:"planAct,attr,omitempty"`
	Cancelled bool     `json:"cancelled,omitempty" xml:"can,attr,omitempty"`
	Platform  string   `json:"plat,omitempty" xml:"plat,attr,omitempty"`
	// CallPtAttributes
	Pta *PublicTime `json:"pta,omitempty" xml:"pta,attr,omitempty"`
	Ptd *PublicTime `json:"ptd,omitempty" xml:"ptd,attr,omitempty"`
	// Working times
	Wta *WorkingTime `json:"wta,omitempty" xml:"wta,attr,omitempty"`
	Wtd *WorkingTime `json:"wtd,omitempty" xml:"wtd,attr,omitempty"`
	Wtp *WorkingTime `json:"wtp,omitempty" xml:"wtp,attr,omitempty"`
	// Delay implied by a change to the services route
	RDelay string `json:"rdelay,omitempty" xml:"rdelay,attr,omitempty"`
	// False destination to be used at this location
	FalseDest string `json:"fd,omitempty" xml:"fd,attr,omitempty"`
}
```

Common location object used in persistence

#### func (*Location) Read

```go
func (t *Location) Read(c *codec.BinaryCodec)
```

#### func (*Location) Write

```go
func (t *Location) Write(c *codec.BinaryCodec)
```

#### type OPDT

```go
type OPDT struct {
	// SchedLocAttributes
	Tiploc    string `xml:"tpl,attr"`
	Act       string `xml:"act,attr"`
	PlanAct   string `xml:"planAct,attr"`
	Cancelled bool   `xml:"can,attr"`
	Platform  string `xml:"plat,attr"`
	// Working times
	Wta string `xml:"wta,attr"`
	Wtd string `xml:"wtd,attr"`
	// Delay implied by a change to the services route
	RDelay string `xml:"rdelay,attr"`
}
```


#### func (*OPDT) Location

```go
func (s *OPDT) Location() *Location
```

#### type OPIP

```go
type OPIP struct {
	// SchedLocAttributes
	Tiploc    string `xml:"tpl,attr"`
	Act       string `xml:"act,attr"`
	PlanAct   string `xml:"planAct,attr"`
	Cancelled bool   `xml:"can,attr"`
	Platform  string `xml:"plat,attr"`
	// Working times
	Wta string `xml:"wta,attr"`
	Wtd string `xml:"wtd,attr"`
	// Delay implied by a change to the services route
	RDelay string `xml:"rdelay,attr"`
}
```


#### func (*OPIP) Location

```go
func (s *OPIP) Location() *Location
```

#### type OPOR

```go
type OPOR struct {
	// SchedLocAttributes
	Tiploc    string `xml:"tpl,attr"`
	Act       string `xml:"act,attr"`
	PlanAct   string `xml:"planAct,attr"`
	Cancelled bool   `xml:"can,attr"`
	Platform  string `xml:"plat,attr"`
	// Working times
	Wta string `xml:"wta,attr"`
	Wtd string `xml:"wtd,attr"`
}
```


#### func (*OPOR) Location

```go
func (s *OPOR) Location() *Location
```

#### type OR

```go
type OR struct {
	// SchedLocAttributes
	Tiploc    string `xml:"tpl,attr"`
	Act       string `xml:"act,attr"`
	PlanAct   string `xml:"planAct,attr"`
	Cancelled bool   `xml:"can,attr"`
	Platform  string `xml:"plat,attr"`
	// CallPtAttributes
	Pta string `xml:"pta,attr"`
	Ptd string `xml:"ptd,attr"`
	// Working times
	Wta string `xml:"wta,attr"`
	Wtd string `xml:"wtd,attr"`
	// False destination to be used at this location
	FalseDest string `xml:"fd,attr"`
}
```


#### func (*OR) Location

```go
func (s *OR) Location() *Location
```

#### type PP

```go
type PP struct {
	// SchedLocAttributes
	Tiploc    string `xml:"tpl,attr"`
	Act       string `xml:"act,attr"`
	PlanAct   string `xml:"planAct,attr"`
	Cancelled bool   `xml:"can,attr"`
	Platform  string `xml:"plat,attr"`
	// Working times
	Wtp string `xml:"wtp,attr"`
	// Delay implied by a change to the services route
	RDelay string `xml:"rdelay,attr"`
}
```


#### func (*PP) Location

```go
func (s *PP) Location() *Location
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

#### func (*SSD) Write

```go
func (t *SSD) Write(c *codec.BinaryCodec)
```
BinaryCodec writer

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

#### func (*WorkingTime) Write

```go
func (t *WorkingTime) Write(c *codec.BinaryCodec)
```
BinaryCodec writer
