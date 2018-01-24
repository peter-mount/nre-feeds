# darwinrest
--
    import "github.com/peter-mount/darwin/darwinrest"

darwinrest provides some additional rest services which use all of the other
packages in forming their results

## Usage

#### type DarwinRest

```go
type DarwinRest struct {
	Ref *darwinref.DarwinReference
	TT  *darwintimetable.DarwinTimetable
}
```


#### func (*DarwinRest) JourneyHandler

```go
func (rs *DarwinRest) JourneyHandler(r *rest.Rest) error
```
JourneyHandler returns a Journey from the timetable and any reference data

#### func (DarwinRest) RegisterRest

```go
func (r DarwinRest) RegisterRest(c *rest.ServerContext)
```
RegisterRest registers the rest endpoints into a ServerContext
