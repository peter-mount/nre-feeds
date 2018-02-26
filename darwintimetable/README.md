# darwintimetable
--
    import "github.com/peter-mount/nre-feeds/darwintimetable"

Reference timetable


Reference timetable

Reference timetable


Reference timetable

Reference timetable

Reference timetable

## Usage

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

#### type DarwinTimetableClient

```go
type DarwinTimetableClient struct {
	// The url prefix, e.g. "http://localhost:8080" of the remote service
	// Note no trailing "/" as the client will add a patch starting with "/"
	Url string
}
```

A remove client to the DarwinTimetable microservice

#### func (*DarwinTimetableClient) GetJourney

```go
func (c *DarwinTimetableClient) GetJourney(rid string) (*Journey, error)
```
GetJourney returns a Journey by making an HTTP call to a remote instance of
DarwinTimetable

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
	SSD       util.SSD `json:"ssd" xml:"ssd,attr"`
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
	Pta *util.PublicTime `json:"pta,omitempty" xml:"pta,attr,omitempty"`
	Ptd *util.PublicTime `json:"ptd,omitempty" xml:"ptd,attr,omitempty"`
	// Working times
	Wta *util.WorkingTime `json:"wta,omitempty" xml:"wta,attr,omitempty"`
	Wtd *util.WorkingTime `json:"wtd,omitempty" xml:"wtd,attr,omitempty"`
	Wtp *util.WorkingTime `json:"wtp,omitempty" xml:"wtp,attr,omitempty"`
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
