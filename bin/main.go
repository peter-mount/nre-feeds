// CIF Rest server
package main

import (
  "darwinref"
  "darwinrest"
  "darwintimetable"
  "darwinupdate"
  "flag"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/golib/statistics"
  "gopkg.in/robfig/cron.v2"
  "log"
  "os"
  "os/signal"
  "syscall"
)

func main() {
  log.Println( "darwin v0.1" )

  refFile := flag.String( "ref", "", "The reference database file" )
  ttFile := flag.String( "timetable", "", "The timetable database file" )
  ftpPassword := flag.String( "ftp", "", "The FTP Password at National Rail" )

  // Port for the webserver
  port := flag.Int( "p", 8080, "Port to use" )

  flag.Parse()

  crontab := cron.New()

  stats := statistics.Statistics{ Log: true }
  stats.Configure()

  server := rest.NewServer( *port )

  var ref *darwinref.DarwinReference
  if *refFile != "" {
    ref = &darwinref.DarwinReference{}

    if err := ref.OpenDB( *refFile ); err != nil {
      log.Fatal( err )
    }

    ref.RegisterRest( server.Context( "/ref" ) )
  }

  var tt *darwintimetable.DarwinTimetable
  if *ttFile != "" {
    tt = &darwintimetable.DarwinTimetable{}

    if err := tt.OpenDB( *ttFile ); err != nil {
      log.Fatal( err )
    }

    tt.RegisterRest( server.Context( "/timetable" ) )

    tt.ScheduleCleanup( crontab )
  }

  if *ftpPassword != "" {
    ftp := &darwinupdate.DarwinUpdate{
      Ref: ref,
      TT: tt,
      Pass: *ftpPassword,
    }

    ftp.Setup( server.Context( "/update" ), crontab )
  }

  rst := &darwinrest.DarwinRest{
    Ref: ref,
    TT: tt,
  }
  // These apply to the root
  rst.RegisterRest( server.Context( "" ) )

  // Listen to signals & close the db before exiting
  // SIGINT for ^C, SIGTERM for docker stopping the container
  sigs := make( chan os.Signal, 1 )
  signal.Notify( sigs, syscall.SIGINT, syscall.SIGTERM )
  go func() {
    sig := <-sigs
    log.Println( "Signal", sig )

    crontab.Stop()

    if( ref != nil ) {
      ref.Close()
    }

    if( tt != nil ) {
      tt.Close()
    }

    log.Println( "Database closed" )

    os.Exit( 0 )
  }()

  crontab.Start()
  server.Start()
}
