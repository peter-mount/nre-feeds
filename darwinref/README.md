# darwinref
--
    import "github.com/peter-mount/darwin/darwinref"

Unmarshal the Darwin Reference XML

## Usage

#### type CISSource

```go
type CISSource struct {
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
}
```


#### type CrsResponse

```go
type CrsResponse struct {
	XMLName xml.Name    `json:"-" xml:"crs"`
	Crs     string      `json:"crs" xml:"crs,attr"`
	Tiploc  []*Location `json:"locations,omitempty" xml:"LocationRef"`
	Self    string      `json:"self,omitempty" xml:"self,attr,omitempty"`
}
```


#### type DarwinReference

```go
type DarwinReference struct {
}
```

Processed reference format

#### func (*DarwinReference) AllReasonCancelHandler

```go
func (dr *DarwinReference) AllReasonCancelHandler(r *rest.Rest) error
```

#### func (*DarwinReference) AllReasonLateHandler

```go
func (dr *DarwinReference) AllReasonLateHandler(r *rest.Rest) error
```

#### func (*DarwinReference) AllTocsHandler

```go
func (dr *DarwinReference) AllTocsHandler(r *rest.Rest) error
```

#### func (*DarwinReference) Close

```go
func (r *DarwinReference) Close()
```
Close the database. If OpenDB() was used to open the db then that db is closed.
If UseDB() was used this simply detaches the DarwinReference from that DB. The
DB is not closed()

#### func (*DarwinReference) CrsHandler

```go
func (dr *DarwinReference) CrsHandler(r *rest.Rest) error
```

#### func (*DarwinReference) GetCancellationReason

```go
func (r *DarwinReference) GetCancellationReason(tx *bolt.Tx, id int) (*Reason, bool)
```
GetToc returns details of a TOC

#### func (*DarwinReference) GetCrs

```go
func (r *DarwinReference) GetCrs(tx *bolt.Tx, t string) ([]*Location, bool)
```
Return a *Location for a tiploc

#### func (*DarwinReference) GetCrsBucket

```go
func (r *DarwinReference) GetCrsBucket(crsbucket *bolt.Bucket, tiplocbucket *bolt.Bucket, crs string) ([]*Location, bool)
```

#### func (*DarwinReference) GetLateReason

```go
func (r *DarwinReference) GetLateReason(tx *bolt.Tx, id int) (*Reason, bool)
```
GetToc returns details of a TOC

#### func (*DarwinReference) GetReasonBucket

```go
func (r *DarwinReference) GetReasonBucket(bucket *bolt.Bucket, id int) (*Reason, bool)
```

#### func (*DarwinReference) GetTiploc

```go
func (r *DarwinReference) GetTiploc(tx *bolt.Tx, tpl string) (*Location, bool)
```
Return a *Location for a tiploc

#### func (*DarwinReference) GetTiplocBucket

```go
func (r *DarwinReference) GetTiplocBucket(bucket *bolt.Bucket, tpl string) (*Location, bool)
```

#### func (*DarwinReference) GetToc

```go
func (r *DarwinReference) GetToc(tx *bolt.Tx, toc string) (*Toc, bool)
```
GetToc returns details of a TOC

#### func (*DarwinReference) GetTocBucket

```go
func (r *DarwinReference) GetTocBucket(bucket *bolt.Bucket, tpl string) (*Toc, bool)
```

#### func (*DarwinReference) GetVia

```go
func (r *DarwinReference) GetVia(tx *bolt.Tx, at string, dest string, loc1 string, loc2 string) (*Via, bool)
```
GetToc returns details of a TOC

#### func (*DarwinReference) GetViaBucket

```go
func (r *DarwinReference) GetViaBucket(bucket *bolt.Bucket, at string, dest string, loc1 string, loc2 string) (*Via, bool)
```

#### func (*DarwinReference) ImportHandler

```go
func (dr *DarwinReference) ImportHandler(r *rest.Rest) error
```

#### func (*DarwinReference) OpenDB

```go
func (r *DarwinReference) OpenDB(dbFile string) error
```
OpenDB opens a DarwinReference database.

#### func (*DarwinReference) Read

```go
func (t *DarwinReference) Read(c *codec.BinaryCodec)
```

#### func (*DarwinReference) ReasonCancelHandler

```go
func (dr *DarwinReference) ReasonCancelHandler(r *rest.Rest) error
```

#### func (*DarwinReference) ReasonLateHandler

```go
func (dr *DarwinReference) ReasonLateHandler(r *rest.Rest) error
```

#### func (DarwinReference) RegisterRest

```go
func (r DarwinReference) RegisterRest(c *rest.ServerContext)
```
RegisterRest registers the rest endpoints into a ServerContext

#### func (*DarwinReference) TimetableId

```go
func (r *DarwinReference) TimetableId() string
```
Return's the timetableId for this reference dataset

#### func (*DarwinReference) TiplocHandler

```go
func (dr *DarwinReference) TiplocHandler(r *rest.Rest) error
```

#### func (*DarwinReference) TocHandler

```go
func (dr *DarwinReference) TocHandler(r *rest.Rest) error
```

#### func (*DarwinReference) UnmarshalXML

```go
func (r *DarwinReference) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error
```

#### func (*DarwinReference) Update

```go
func (r *DarwinReference) Update(f func(*bolt.Tx) error) error
```
Update performs a read write opertation on the database

#### func (*DarwinReference) UseDB

```go
func (r *DarwinReference) UseDB(db *bolt.DB) error
```
UseDB Allows an already open database to be used with DarwinReference.

#### func (*DarwinReference) ViaHandler

```go
func (dr *DarwinReference) ViaHandler(r *rest.Rest) error
```
ViaHandler returns the unique instance of a via entry

#### func (*DarwinReference) View

```go
func (r *DarwinReference) View(f func(*bolt.Tx) error) error
```
View performs a readonly operation on the database

#### func (*DarwinReference) Write

```go
func (t *DarwinReference) Write(c *codec.BinaryCodec)
```

#### type Location

```go
type Location struct {
	XMLName xml.Name `json:"-" xml:"LocationRef"`
	Tiploc  string   `json:"tpl" xml:"tpl,attr"`
	Crs     string   `json:"crs,omitempty" xml:"crs,attr,omitempty"`
	Toc     string   `json:"toc,omitempty" xml:"toc,attr,omitempty"`
	Name    string   `json:"locname" xml:"locname,attr"`
	// Date entry was inserted into the database
	Date time.Time `json:"date" xml:"date,attr"`
	// URL to this entity
	Self string `json:"self,omitempty" xml:"self,attr,omitempty"`
}
```

Defines a location, i.e. Station or passing point

#### func (*Location) Equals

```go
func (a *Location) Equals(b *Location) bool
```

#### func (*Location) Read

```go
func (t *Location) Read(c *codec.BinaryCodec)
```

#### func (*Location) SetSelf

```go
func (t *Location) SetSelf(r *rest.Rest)
```
SetSelf sets the Self field to match this request

#### func (*Location) Write

```go
func (t *Location) Write(c *codec.BinaryCodec)
```

#### type LocationMap

```go
type LocationMap struct {
}
```


#### func  NewLocationMap

```go
func NewLocationMap() *LocationMap
```

#### func (*LocationMap) Add

```go
func (r *LocationMap) Add(t *Location)
```
AddTiploc adds a Tiploc to the response

#### func (*LocationMap) AddAll

```go
func (r *LocationMap) AddAll(t []*Location)
```
AddTiplocs adds an array of Tiploc's to the response

#### func (*LocationMap) AddTiploc

```go
func (r *LocationMap) AddTiploc(dr *DarwinReference, tx *bolt.Tx, t string)
```

#### func (*LocationMap) AddTiplocs

```go
func (r *LocationMap) AddTiplocs(dr *DarwinReference, tx *bolt.Tx, ts []string)
```

#### func (*LocationMap) ForEach

```go
func (r *LocationMap) ForEach(f func(*Location))
```

#### func (*LocationMap) Get

```go
func (r *LocationMap) Get(n string) (*Location, bool)
```

#### func (*LocationMap) MarshalJSON

```go
func (t *LocationMap) MarshalJSON() ([]byte, error)
```

#### func (*LocationMap) Self

```go
func (r *LocationMap) Self(rs *rest.Rest)
```
Self sets the Self field to match this request

#### type Reason

```go
type Reason struct {
	XMLName   xml.Name `json:"-" xml:"Reason"`
	Code      int      `json:"code" xml:"code,attr"`
	Text      string   `json:"reasontext" xml:"reasontext,attr"`
	Cancelled bool     `json:"canc" xml:"canc,attr"`
	// Date entry was inserted into the database
	Date time.Time `json:"date" xml:"date,attr"`
	// URL to this entity
	Self string `json:"self" xml:"self,attr,omitempty"`
}
```

A reason, shared by LateRunningReasons and CancellationReasons

#### func (*Reason) Equals

```go
func (a *Reason) Equals(b *Reason) bool
```

#### func (*Reason) Read

```go
func (t *Reason) Read(c *codec.BinaryCodec)
```

#### func (*Reason) SetSelf

```go
func (t *Reason) SetSelf(r *rest.Rest)
```
SetSelf sets the Self field to match this request

#### func (*Reason) Write

```go
func (t *Reason) Write(c *codec.BinaryCodec)
```

#### type ReasonsResponse

```go
type ReasonsResponse struct {
	XMLName xml.Name  `json:"-" xml:"reasons"`
	Reasons []*Reason `json:"reasons,omitempty" xml:"Reason"`
	Self    string    `json:"self,omitempty" xml:"self,attr,omitempty"`
}
```


#### type Toc

```go
type Toc struct {
	XMLName xml.Name `json:"-" xml:"TocRef"`
	Toc     string   `json:"toc" xml:"toc,attr"`
	Name    string   `json:"tocname" xml:"tocname,attr"`
	Url     string   `json:"url" xml:"url,attr"`
	// Date entry was inserted into the database
	Date time.Time `json:"date" xml:"date,attr"`
	// URL to this entity
	Self string `json:"self" xml:"self,attr,omitempty"`
}
```

A rail operator

#### func (*Toc) Equals

```go
func (a *Toc) Equals(b *Toc) bool
```

#### func (*Toc) Read

```go
func (t *Toc) Read(c *codec.BinaryCodec)
```

#### func (*Toc) SetSelf

```go
func (t *Toc) SetSelf(r *rest.Rest)
```

#### func (*Toc) Write

```go
func (t *Toc) Write(c *codec.BinaryCodec)
```

#### type TocMap

```go
type TocMap struct {
}
```


#### func  NewTocMap

```go
func NewTocMap() *TocMap
```

#### func (*TocMap) Add

```go
func (r *TocMap) Add(t *Toc)
```
AddTiploc adds a Tiploc to the response

#### func (*TocMap) AddAll

```go
func (r *TocMap) AddAll(t []*Toc)
```
AddTiplocs adds an array of Tiploc's to the response

#### func (*TocMap) AddToc

```go
func (r *TocMap) AddToc(dr *DarwinReference, tx *bolt.Tx, t string)
```

#### func (*TocMap) AddTocs

```go
func (r *TocMap) AddTocs(dr *DarwinReference, tx *bolt.Tx, ts []string)
```

#### func (*TocMap) Get

```go
func (r *TocMap) Get(n string) (*Toc, bool)
```

#### func (*TocMap) MarshalJSON

```go
func (t *TocMap) MarshalJSON() ([]byte, error)
```

#### func (*TocMap) Self

```go
func (r *TocMap) Self(rs *rest.Rest)
```
Self sets the Self field to match this request

#### type TocsResponse

```go
type TocsResponse struct {
	XMLName xml.Name `json:"-" xml:"tocs"`
	Toc     []*Toc   `json:"tocs,omitempty" xml:"TocRef"`
	Self    string   `json:"self,omitempty" xml:"self,attr,omitempty"`
}
```


#### type Via

```go
type Via struct {
	XMLName xml.Name `json:"-" xml:"Via"`
	At      string   `json:"at" xml:"at,attr"`
	Dest    string   `json:"dest" xml:"dest,attr"`
	Loc1    string   `json:"loc1" xml:"loc1,attr"`
	Loc2    string   `json:"loc2,omitempty" xml:"loc2,attr,omitempty"`
	Text    string   `json:"text" xml:"viatext,attr"`
	// Date entry was inserted into the database
	Date time.Time `json:"date" xml:"date,attr"`
	// URL to this entity
	Self string `json:"self" xml:"self,attr,omitempty"`
}
```

Via text

#### func (*Via) Equals

```go
func (v *Via) Equals(o *Via) bool
```
Are two Via's equal

#### func (*Via) Read

```go
func (v *Via) Read(c *codec.BinaryCodec)
```

#### func (*Via) SetSelf

```go
func (v *Via) SetSelf(r *rest.Rest)
```
SetSelf sets the Self field to match this request

#### func (*Via) String

```go
func (v *Via) String() string
```

#### func (*Via) Write

```go
func (v *Via) Write(c *codec.BinaryCodec)
```
