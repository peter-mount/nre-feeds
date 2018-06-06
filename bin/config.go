// Internal library used for the binary webservices
package bin

import (
  "github.com/peter-mount/golib/rabbitmq"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/golib/statistics"
  "github.com/peter-mount/nre-feeds/darwinupdate"
  "gopkg.in/robfig/cron.v2"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "path/filepath"
)

// Common configuration used to read config.yaml
type Config struct {
  // URL prefixes for lookups to the reference microservices
  Services struct {
    DarwinD3  string        `yaml:"darwind3"`
    Reference string        `yaml:"reference"`
    Timetable string        `yaml:"timetable"`
  }                         `yaml:"services"`

  Database struct {
    Path          string    `yaml:path`
    // Darwin Reference
    Reference     string    `yaml:"reference"`
    // Darwin Timetable
    Timetable     string    `yaml:"timetable"`
    // Darwin PushPort
    PushPort      string    `yaml:"pushPort"`
  }                         `yaml:"database"`

  Ftp struct {
    Enabled       bool      `yaml:"enabled"`
    Server        string    `yaml:"server"`
    User          string    `yaml:"user"`
    Password      string    `yaml:"password"`
    Schedule      string    `yaml:"schedule"`
    Update       *darwinupdate.DarwinUpdate
  }                         `yaml:"ftp"`

  RabbitMQ      rabbitmq.RabbitMQ `yaml:"rabbitmq"`

  D3 struct {
    ResolveSched  bool      `yaml:"resolveSchedules"`
    QueueName     string    `yaml:"queueName"`
    RoutingKey    string    `yaml:"routingKey"`
  }                         `yaml:"d3"`

  Server struct {
    // Root context path, defaults to ""
    Context       string    `yaml:"context"`
    // The port to run on, defaults to 80
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
func (c *Config) readFile( configFile string ) error {
  if filename, err := filepath.Abs( configFile ); err != nil {
    return err
  } else if in, err := ioutil.ReadFile( filename ); err != nil {
    return err
  } else if err := yaml.Unmarshal( in, c ); err != nil {
    return err
  }
  return nil
}
