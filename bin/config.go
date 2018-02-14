// Internal library used for the binary webservices
package bin

import (
  "darwind3"
  "ldb"
  "darwinref"
  "darwintimetable"
  "darwinupdate"
  "github.com/peter-mount/golib/rabbitmq"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/golib/statistics"
  "gopkg.in/robfig/cron.v2"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "path/filepath"
)

// Common configuration used to read config.yaml
type Config struct {
  Database struct {
    Path          string    `yaml:path`
    // Darwin Reference
    Reference     string    `yaml:"reference"`
    reference    *darwinref.DarwinReference
    // Darwin Timetable
    Timetable     string    `yaml:"timetable"`
    timetable    *darwintimetable.DarwinTimetable
    // Darwin PushPort
    PushPort      string    `yaml:"pushPort"`
    pushPort     *darwind3.DarwinD3
  }                         `yaml:"database"`

  Ftp struct {
    Enabled       bool      `yaml:"enabled"`
    Server        string    `yaml:"server"`
    User          string    `yaml:"user"`
    Password      string    `yaml:"password"`
    Schedule      string    `yaml:"schedule"`
    Update       *darwinupdate.DarwinUpdate
  }                         `yaml:"ftp"`

  PushPort struct {
    Enabled       bool      `yaml:"enabled"`
    ResolveSched  bool      `yaml:"resolveSchedules"`
    QueueName     string    `yaml:"queueName"`
    RoutingKey    string    `yaml:"routingKey"`
    RabbitMQ      rabbitmq.RabbitMQ `yaml:"rabbitmq"`
  }                         `yaml:"pushPort"`

  LDB struct {
    Enabled       bool      `yaml:"enabled"`
    ldb          *ldb.LDB
  }                         `yaml:"ldb"`

  Server struct {
    Context       string    `yaml:"context"`
    // The port to run on, defaults to 8080
    Port          int       `yaml:"port"`
    // The permitted headers
    Headers     []string
    // The permitted Origins
    Origins     []string
    // The permitted methods
    Methods     []string
    // Web Server
    server       *rest.Server
    // Base Context
    Ctx          *rest.ServerContext
  }                         `yaml:"server"`

  Statistics struct {
    Log           bool      `yaml:"log"`
    Rest          string    `yaml:"rest"`
    Schedule      string    `yaml:"schedule"`
    statistics   *statistics.Statistics
  }                         `yaml:"statistics"`

  // Cron
  Cron         *cron.Cron
}

// ReadFile reads the provided file and imports yaml config
func (c *Config) ReadFile( configFile string ) error {
  if filename, err := filepath.Abs( configFile ); err != nil {
    return err
  } else if in, err := ioutil.ReadFile( filename ); err != nil {
    return err
  } else {
    return c.Unmarshal( in )
  }
}

// Unmarshal reads yaml from a byte slice
func (c *Config) Unmarshal( in []byte ) error {
  if err := yaml.Unmarshal( in, c ); err != nil {
    return err
  }
  return nil
}
