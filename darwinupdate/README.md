# darwinupdate
--
    import "github.com/peter-mount/darwin/darwinupdate"

Package that handles FTP updates

## Usage

#### func  FtpLs

```go
func FtpLs(con *ftp.ServerConn, path string) error
```
FtpLs utility to log the files in a path

#### type DarwinUpdate

```go
type DarwinUpdate struct {
	// DarwinReference instance or nil
	Ref *darwinref.DarwinReference
	// DarwinTimetable instance or nil
	TT *darwintimetable.DarwinTimetable
	// The server name
	Server string
	// The ftp user
	User string
	// The ftp password for the NRE ftp server
	Pass string
}
```


#### func (*DarwinUpdate) InitialImport

```go
func (u *DarwinUpdate) InitialImport()
```

#### func (*DarwinUpdate) ReferenceHandler

```go
func (u *DarwinUpdate) ReferenceHandler(r *rest.Rest) error
```

#### func (*DarwinUpdate) ReferenceUpdate

```go
func (u *DarwinUpdate) ReferenceUpdate(con *ftp.ServerConn) error
```

#### func (*DarwinUpdate) SetupRest

```go
func (u *DarwinUpdate) SetupRest(c *rest.ServerContext)
```

#### func (*DarwinUpdate) SetupSchedule

```go
func (u *DarwinUpdate) SetupSchedule(cr *cron.Cron, schedule string)
```

#### func (*DarwinUpdate) TimetableHandler

```go
func (u *DarwinUpdate) TimetableHandler(r *rest.Rest) error
```

#### func (*DarwinUpdate) TimetableUpdate

```go
func (u *DarwinUpdate) TimetableUpdate(con *ftp.ServerConn) error
```

#### func (*DarwinUpdate) Update

```go
func (u *DarwinUpdate) Update(force bool) error
```
Update performs an update of all data
