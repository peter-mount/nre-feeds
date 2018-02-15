# ldb
--
    import "github.com/peter-mount/darwin/ldb"

LDB - Live Departure Boards

## Usage

#### type LDB

```go
type LDB struct {
	// Link to D3
	Darwin string
	// Link to reference
	Reference string
	// Eventing
	EventManager *darwind3.DarwinEventManager
	// The managed stations
	Stations *Stations
}
```


#### func (*LDB) GetStationCrs

```go
func (d *LDB) GetStationCrs(crs string) *Station
```
GetStationCrs returns the Station instance by CRS or nil if not found Unlike
GetStationTiploc this will not create a station if it's not found

#### func (*LDB) GetStationTiploc

```go
func (d *LDB) GetStationTiploc(tiploc string) *Station
```
GetStationTiploc returns the Station instance by Tiploc or nil if not found.
Note: If we don't have an entry then this will create one

#### func (*LDB) Init

```go
func (d *LDB) Init() error
```

#### func (*LDB) RegisterRest

```go
func (d *LDB) RegisterRest(c *rest.ServerContext)
```

#### type Service

```go
type Service struct {
	// The RID of this service
	RID string `json:"rid"`
	// The destination
	Destination string `json:"destination"`
	// Via text
	Via string `json:"via,omitempty"`
	// Service Start Date
	SSD darwintimetable.SSD `json:"ssd"`
	// The trainId (headcode)
	TrainId string `json:"trainId"`
	// The operator of this service
	Toc string `json:"toc"`
	// Is a passenger service
	PassengerService bool `json:"passengerService,omitempty"`
	// Is a charter service
	Charter bool `json:"charter,omitempty"`
	// Cancel running reason for this service. The reason applies to all locations
	// of this service which are marked as cancelled
	CancelReason darwind3.DisruptionReason `json:"cancelReason"`
	// Late running reason for this service. The reason applies to all locations
	// of this service which are not marked as cancelled
	LateReason darwind3.DisruptionReason `json:"lateReason"`
	// The "time" for this service
	Location *darwind3.Location `json:"location"`
	// The time this entry was set
	Date time.Time `json:"date,omitempty" xml:"date,attr,omitempty"`
	// URL to the train detail page
	Self string `json:"self,omitempty" xml:"self,attr,omitempty"`
}
```

A representation of a service at a location

#### func (*Service) Clone

```go
func (a *Service) Clone() *Service
```
Clone returns a copy of this Service

#### func (*Service) Compare

```go
func (a *Service) Compare(b *Service) bool
```
Compare two services by the times at a location

#### func (*Service) MarshalJSON

```go
func (t *Service) MarshalJSON() ([]byte, error)
```

#### func (*Service) Timestamp

```go
func (s *Service) Timestamp() time.Time
```
Timestamp returns the time.Time of this service based on the SSD and Location's
Time. TODO this does not currently handle midnight correctly

#### type Station

```go
type Station struct {
	// The location details for this station
	Locations []*darwinref.Location
	Crs       string
}
```

The holder for a station's departure boards

#### func (*Station) Update

```go
func (s *Station) Update(f func() error) error
```
Perform an action on the station with an exclusive lock

#### type Stations

```go
type Stations struct {
}
```

Manages all stations

#### func  NewStations

```go
func NewStations() *Stations
```

#### func (*Stations) Cleanup

```go
func (st *Stations) Cleanup()
```
Cleanup removes any old schedules still in memory for each station

#### func (*Stations) Update

```go
func (s *Stations) Update(f func() error) error
```
Perform an action on the Stations instance with an exclusive lock
