package main

import (
  "darwind3"
  "ldb"
  "darwinref"
  "darwintimetable"
  "darwinupdate"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/golib/statistics"
  "gopkg.in/robfig/cron.v2"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "path/filepath"
)

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
    update       *darwinupdate.DarwinUpdate
  }                         `yaml:"ftp"`

  PushPort struct {
    Enabled       bool      `yaml:"enabled"`
    ResolveSched  bool      `yaml:"resolveSchedules"`
  }                         `yaml:"pushPort"`

  LDB struct {
    Enabled       bool      `yaml:"enabled"`
    ldb          *ldb.LDB
  }                         `yaml:"ldb"`

  Server struct {
    Context       string    `yaml:"context"`
    Port          int       `yaml:"port"`
    // Web Server
    server       *rest.Server
    // Base Context
    ctx          *rest.ServerContext
  }                         `yaml:"server"`

  Statistics struct {
    Log           bool      `yaml:"log"`
    Rest          string    `yaml:"rest"`
    Schedule      string    `yaml:"schedule"`
    statistics   *statistics.Statistics
  }                         `yaml:"statistics"`

  // Cron
  cron         *cron.Cron
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
