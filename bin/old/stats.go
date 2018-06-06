package bin

import (
  "github.com/peter-mount/golib/statistics"
)

// initStats initialises the statistics subsystem
func (c *Config) initStats() error {

  c.Statistics.statistics = &statistics.Statistics{
    Log: c.Statistics.Log,
    Schedule: c.Statistics.Schedule,
    // Use our cron rather than create a second one
    Cron: c.Cron,
  }

  c.Statistics.statistics.Configure()

  if c.Statistics.Rest != "" {
    c.Server.Ctx.HandleFunc( "/stats", statistics.StatsRestHandler ).Methods( "GET" )
  }

  return nil
}
