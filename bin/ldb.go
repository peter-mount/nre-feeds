package main

import (
  "ldb"
)

func (c *Config) initLdb() error {

  // LDB only valid with pushport
  if c.PushPort.Enabled && c.LDB.Enabled {
    c.LDB.ldb = &ldb.LDB{
      Darwin: c.Database.pushPort,
      Reference: c.Database.reference,
    }
    if err := c.LDB.ldb.Init(); err != nil {
      return err
    }

    c.LDB.ldb.RegisterRest( c.Server.ctx.Context( "/ldb" ) )

    // Expire old messages every 15 minutes
    c.cron.AddFunc( "0 0/15 * * * *", c.LDB.ldb.Darwin.ExpireStationMessages )

    // Expire old schedules every 15 minutes
    c.cron.AddFunc( "0 * * * * *", c.LDB.ldb.Stations.Cleanup )
  }

  return nil
}
