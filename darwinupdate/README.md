# darwinupdate
--
    import "github.com/peter-mount/nre-feeds/darwinupdate"

Package that handles FTP updates from the NRE FTP server

## Usage

#### func  FtpLs

```go
func FtpLs(con *ftp.ServerConn, path string) error
```
FtpLs utility to log the files in a path

#### type DarwinUpdate

```go
type DarwinUpdate struct {
	// The server name
	Server string
	// The ftp user
	User string
	// The ftp password for the NRE ftp server
	Pass string
}
```


#### func (*DarwinUpdate) Ftp

```go
func (u *DarwinUpdate) Ftp(f func(*ftp.ServerConn) error) error
```

#### func (*DarwinUpdate) ImportRequiredTimetable

```go
func (u *DarwinUpdate) ImportRequiredTimetable(v interface{ TimetableId() string }) bool
```

#### func (*DarwinUpdate) ReferenceUpdate

```go
func (u *DarwinUpdate) ReferenceUpdate(ref *darwinref.DarwinReference) error
```

#### func (*DarwinUpdate) TimetableUpdate

```go
func (u *DarwinUpdate) TimetableUpdate(tt *darwintimetable.DarwinTimetable) error
```
