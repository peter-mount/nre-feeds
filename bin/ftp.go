package main

import (
  "darwinupdate"
)

func (c *Config) defaultValue( s *string, d string ) *Config {
  if *s == "" {
    *s = d
  }
  return c
}

func (c *Config) initFtp() error {

  // Force ftp offline if no password is set
  if c.Ftp.Password == "" {
    c.Ftp.Enabled = false
  }

  c.defaultValue( &c.Ftp.Server, "datafeeds.nationalrail.co.uk:21" ).
    defaultValue( &c.Ftp.User, "ftpuser" )

  if c.Ftp.Enabled {
    c.Ftp.update = &darwinupdate.DarwinUpdate{
      Ref: c.Database.reference,
      TT: c.Database.timetable,
      Server: c.Ftp.Server,
      User: c.Ftp.User,
      Pass: c.Ftp.Password,
    }

    c.Ftp.update.SetupRest( c.Server.ctx.Context( "/update" ) )

    if c.Ftp.Schedule != "" {
      c.Ftp.update.SetupSchedule( c.cron, c.Ftp.Schedule )
    }

    // Finally check to see if we need to import now
    c.Ftp.update.InitialImport()
  }

  return nil
}
