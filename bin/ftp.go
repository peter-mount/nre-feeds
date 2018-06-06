package bin

import (
  "github.com/peter-mount/nre-feeds/darwinupdate"
)

func (c *Config) defaultValue( s *string, d string ) *Config {
  if *s == "" {
    *s = d
  }
  return c
}

// InitFtp initialises the ftp client. This is exposed as it's not usually used
// within every microservice
func (c *Config) InitFtp() error {

  // Force ftp offline if no password is set
  if c.Ftp.Password == "" {
    c.Ftp.Enabled = false
  }

  if c.Ftp.Enabled {
    c.defaultValue( &c.Ftp.Server, "datafeeds.nationalrail.co.uk:21" ).
    defaultValue( &c.Ftp.User, "ftpuser" )

    // Create the updater
    c.Ftp.Update = &darwinupdate.DarwinUpdate{
      Server: c.Ftp.Server,
      User: c.Ftp.User,
      Pass: c.Ftp.Password,
    }
  }

  return nil
}
