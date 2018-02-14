package bin

import (
  "darwinupdate"
)

func (c *Config) defaultValue( s *string, d string ) *Config {
  if *s == "" {
    *s = d
  }
  return c
}

func (c *Config) InitFtp() error {

  // Force ftp offline if no password is set
  if c.Ftp.Password == "" {
    c.Ftp.Enabled = false
  }

  c.defaultValue( &c.Ftp.Server, "datafeeds.nationalrail.co.uk:21" ).
    defaultValue( &c.Ftp.User, "ftpuser" )

  if c.Ftp.Enabled {
    // Create the updater
    c.Ftp.Update = &darwinupdate.DarwinUpdate{
      Server: c.Ftp.Server,
      User: c.Ftp.User,
      Pass: c.Ftp.Password,
    }
  }

  return nil
}
