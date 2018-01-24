package main

import (
  "github.com/peter-mount/golib/statistics"
)

func (c *Config) initStats() error {

  c.Statistics.statistics = &statistics.Statistics{
    Log: c.Statistics.Log,
    Schedule: c.Statistics.Schedule,
    // Use our cron rather than create a second one
    Cron: c.cron,
  }

  c.Statistics.statistics.Configure()

  if c.Statistics.Rest != "" {
    c.Server.ctx.HandleFunc( "/stats", statistics.StatsRestHandler ).Methods( "GET" )
  }

  return nil
}
