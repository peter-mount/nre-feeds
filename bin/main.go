// Internal library used for the binary webservices
package bin

import (
  "flag"
  "log"
  "os"
  "runtime"
)

// RunApplication starts the common services then runs the supplied function
// to configure the specific application. As long as it returns nil for error
// then the http server is started.
// The optional function in the return will, if not nil, be called when the
// application shuts down.
func RunApplication( app func( *Config ) ( func(), error ) ) {

  log.Printf( "%s %s %s(%s)", os.Args[0], VERSION, runtime.GOOS, runtime.GOARCH )

  configFile := flag.String( "c", "", "The config file to use" )

  flag.Parse()

  if *configFile == "" {
    log.Fatal( "No default config defined, provide with -c" )
  }

  config := &Config{}

  if err := config.readFile( *configFile ); err != nil {
    log.Fatal( err )
  }

  if err := config.initCron(); err != nil {
    log.Fatal( err )
  }

  if err := config.initServer(); err != nil {
    log.Fatal( err )
  }

  if err := config.initStats(); err != nil {
    log.Fatal( err )
  }

  if err := config.initDb(); err != nil {
    log.Fatal( err )
  }

  if close, err := app( config ); err != nil {
    log.Fatal( err )
  } else if err := config.initShutdown( close ); err != nil {
    log.Fatal( err )
  }

  config.Cron.Start()

  if err := config.Server.server.Start(); err != nil {
    log.Fatal( err )
  }
}
