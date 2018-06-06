// Package that handles FTP updates from the NRE FTP server
package darwinupdate

import (
  "fmt"
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/nre-feeds/bin"
)

type DarwinUpdate struct {
  // The server name
  Server  string
  // The ftp user
  User    string
  // The ftp password for the NRE ftp server
  Pass    string
  config *bin.Config
}

func (a *DarwinUpdate) Name() string {
  return "DarwinUpdate"
}

func (a *DarwinUpdate) Init( k *kernel.Kernel ) error {
  service, err := k.AddService( &bin.Config{} )
  if err != nil {
    return err
  }
  a.config = (service).(*bin.Config)

  return nil
}

func (a *DarwinUpdate) PostInit() error {

  if a.config.Ftp.Enabled {
    a.Server = a.config.Ftp.Server
    a.User = a.config.Ftp.User
    a.Pass = a.config.Ftp.Password
    if a.Server == "" {
      a.Server = "datafeeds.nationalrail.co.uk:21"
    }
    if a.User == "" {
      a.User = "ftpuser"
    }

    if a.Pass == "" {
      return fmt.Errorf( "Ftp is enabled but credentials are required" )
    }
  }

  return nil
}
