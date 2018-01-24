package main

import (
  "github.com/peter-mount/golib/statistics"
)

func (c *Config) initStats() error {

  c.Statistics.statistics = &statistics.Statistics{
    Log: c.Statistics.Log,
    Schedule: c.Statistics.Schedule,
  }

  c.Statistics.statistics.Configure()

  return nil
}
