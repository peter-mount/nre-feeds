// Internal library used for the binary webservices
package bin

import (
  "flag"
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/rabbitmq"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "path/filepath"
  "fmt"
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
  }                         `yaml:"ftp"`

  KB struct {
    Username      string    `yaml:"username"`
    Password      string    `yaml:"password"`
    DataDir       string    `yaml:"datadir"`
  }

  RabbitMQ      rabbitmq.RabbitMQ `yaml:"rabbitmq"`

  D3 struct {
    // Set to true to use the darwintt service to try to resulve unknown schedules
    ResolveSched    bool    `yaml:"resolveSchedules"`
    // The queue name to create
    QueueName       string  `yaml:"queueName"`
    // The routingKey to listen for inbound d3 messages
    RoutingKey      string  `yaml:"routingKey"`
    // Prefix to the routingKeys used by the Event subsystem
    EventKeyPrefix  string  `yaml:"eventKeyPrefix"`
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
  }                         `yaml:"server"`

  Statistics struct {
    Log           bool      `yaml:"log"`
    Rest          string    `yaml:"rest"`
    Schedule      string    `yaml:"schedule"`
  }                         `yaml:"statistics"`

  configFile   *string
}

func (a *Config) Name() string {
  return "Config"
}

func (a *Config) Init( k *kernel.Kernel ) error {
  a.configFile = flag.String( "c", "", "The config file to use" )

  return nil
}

func (a *Config) PostInit() error {
  // Verify then load the config file
  if *a.configFile == "" {
    return fmt.Errorf( "No default config defined, provide with -c" )
  }

  if filename, err := filepath.Abs( *a.configFile ); err != nil {
    return err
  } else if in, err := ioutil.ReadFile( filename ); err != nil {
    return err
  } else if err := yaml.Unmarshal( in, a ); err != nil {
    return err
  }

  // Ensure the database path is correct
  if a.Database.Path == "" {
    a.Database.Path = "/database/"
  }

  if path, err := filepath.Abs( a.Database.Path ); err != nil {
    return err
  } else {
    a.Database.Path = path
  }

  if a.Database.Path[len(a.Database.Path)-1] != '/' {
    a.Database.Path = a.Database.Path + "/"
  }

  return nil
}

// DbPath ensures the database name is set. If the name is not absolute then it's
// taken as being relative to the database path in config.
// s The required filename
// d The filename to use if s is ""
func (c *Config) DbPath( s *string, d string ) *Config {
  if *s == "" {
    *s = d
  }

  if (*s)[0] != '/' {
    *s = c.Database.Path + *s
  }

  return c
}
