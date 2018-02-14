package darwinupdate

import (
  "github.com/jlaffaye/ftp"
  "log"
  "time"
)

func (u *DarwinUpdate) Ftp( f func( *ftp.ServerConn ) error ) error {
  log.Println( "FTP: Connecting" )
  if con, err := ftp.DialTimeout( u.Server, time.Minute ); err != nil {
    return err
  } else {
    log.Println( "FTP: Login" )
    if err := con.Login( u.User, u.Pass ); err != nil {
      log.Println( "FTP:", err )
      con.Quit()
      return nil
    }

    if err := f( con ); err != nil {
      log.Println( "FTP: Quit", err )
      con.Quit()
      return err
    }

    log.Println( "FTP: Quit" )
    return con.Quit()
  }
}

// FtpLs utility to log the files in a path
func FtpLs( con *ftp.ServerConn, path string ) error {
  log.Println( "LS", path )
  if entries, err := con.List( path ); err != nil {
    return err
  } else {
    for i, e := range entries {
      log.Printf( "%d %s %s %d %v\n", i, e.Name, e.Type, e.Size, e.Time )
    }
  }

  return nil
}
