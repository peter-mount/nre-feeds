# bin
--
    import "github.com/peter-mount/nre-feeds/bin"

Internal library used for the binary webservices

Internal library used for the binary webservices

## Usage

```go
const (
	VERSION = "@@version@@"
)
```

#### func  RunApplication

```go
func RunApplication(app func(*Config) (func(), error))
```
RunApplication starts the common services then runs the supplied function to
configure the specific application. As long as it returns nil for error then the
http server is started. The optional function in the return will, if not nil, be
called when the application shuts down.

#### type Config

```go
type Config struct {
	// URL prefixes for lookups to the reference microservices
	Services struct {
		DarwinD3  string `yaml:"darwind3"`
		Reference string `yaml:"reference"`
		Timetable string `yaml:"timetable"`
	} `yaml:"services"`

	Database struct {
		Path string `yaml:path`
		// Darwin Reference
		Reference string `yaml:"reference"`
		// Darwin Timetable
		Timetable string `yaml:"timetable"`
		// Darwin PushPort
		PushPort string `yaml:"pushPort"`
	} `yaml:"database"`

	Ftp struct {
		Enabled  bool   `yaml:"enabled"`
		Server   string `yaml:"server"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Schedule string `yaml:"schedule"`
		Update   *darwinupdate.DarwinUpdate
	} `yaml:"ftp"`

	RabbitMQ rabbitmq.RabbitMQ `yaml:"rabbitmq"`

	D3 struct {
		ResolveSched bool   `yaml:"resolveSchedules"`
		QueueName    string `yaml:"queueName"`
		RoutingKey   string `yaml:"routingKey"`
	} `yaml:"d3"`

	Server struct {
		// Root context path, defaults to ""
		Context string `yaml:"context"`
		// The port to run on, defaults to 80
		Port int `yaml:"port"`
		// The permitted headers
		Headers []string
		// The permitted Origins
		Origins []string
		// The permitted methods
		Methods []string

		// Base Context
		Ctx *rest.ServerContext
	} `yaml:"server"`

	Statistics struct {
		Log      bool   `yaml:"log"`
		Rest     string `yaml:"rest"`
		Schedule string `yaml:"schedule"`
	} `yaml:"statistics"`

	// Cron
	Cron *cron.Cron
}
```

Common configuration used to read config.yaml

#### func (*Config) DbPath

```go
func (c *Config) DbPath(s *string, d string) *Config
```
DbPath ensures the database name is set. If the name is not absolute then it's
taken as being relative to the database path in config. s The required filename
d The filename to use if s is ""

#### func (*Config) InitFtp

```go
func (c *Config) InitFtp() error
```
InitFtp initialises the ftp client. This is exposed as it's not usually used
within every microservice
