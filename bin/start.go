package main

import (
  "darwinrest"
)

func (c *Config) start() error {

  // Shared rest services
  rst := &darwinrest.DarwinRest{
    Ref: c.Database.reference,
    TT: c.Database.timetable,
  }
  // These apply to the root context
  rst.RegisterRest( c.Server.ctx )

  c.cron.Start()

  return c.Server.server.Start()
}
