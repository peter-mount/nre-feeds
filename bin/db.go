package main

import (
  "darwind3"
  "ldb"
  "darwinref"
  "darwintimetable"
  "path/filepath"
)

func (c *Config) dbPath( s *string, d string ) *Config {
  if *s == "" {
    *s = c.Database.Path + d
  } else if (*s)[0] != '/' {
    *s = c.Database.Path + *s
  }
  return c
}

func (c *Config) initDb() error {

  if c.Database.Path == "" {
    c.Database.Path = "/database/"
  }

  if path, err := filepath.Abs( c.Database.Path ); err != nil {
    return err
  } else {
    c.Database.Path = path
  }

  if c.Database.Path[len(c.Database.Path)-1] != '/' {
    c.Database.Path = c.Database.Path + "/"
  }

  c.Database.reference = &darwinref.DarwinReference{}
  c.dbPath( &c.Database.Reference, "dwref.db" )
  if err := c.Database.reference.OpenDB( c.Database.Reference ); err != nil {
    return err
  }
  c.Database.reference.RegisterRest( c.Server.ctx.Context( "/ref" ) )

  c.dbPath( &c.Database.Timetable, "dwtt.db" )
  c.Database.timetable = &darwintimetable.DarwinTimetable{}
  if err := c.Database.timetable.OpenDB( c.Database.Timetable ); err != nil {
    return err
  }
  c.Database.timetable.RegisterRest( c.Server.ctx.Context( "/timetable" ) )
  c.Database.timetable.ScheduleCleanup( c.cron )

  if c.PushPort.Enabled {
    c.dbPath( &c.Database.PushPort, "dwlive.db" )
    c.Database.pushPort = &darwind3.DarwinD3{}
    if err := c.Database.pushPort.OpenDB( c.Database.PushPort ); err != nil {
      return err
    }

    if c.PushPort.ResolveSched {
      // Allow D3 to resolve schedules from the timetable
      c.Database.pushPort.Timetable = c.Database.timetable
    }

    c.Database.pushPort.RegisterRest( c.Server.ctx.Context( "/live" ) )
  }

  // LDB only valid with pushport
  if c.PushPort.Enabled && c.LDB.Enabled {
    c.dbPath( &c.Database.LDB, "dwldb.db" )
    c.Database.ldb = &ldb.LDB{
      Darwin: c.Database.pushPort,
      Reference: c.Database.reference,
    }
    if err := c.Database.ldb.OpenDB( c.Database.LDB ); err != nil {
      return err
    }
    c.Database.ldb.RegisterRest( c.Server.ctx.Context( "/ldb" ) )
  }

  return nil
}
