package bin

import (
  "github.com/peter-mount/golib/rest"
)

func (c *Config) InitServer() error {

  if c.Server.Port < 1 || c.Server.Port > 65534 {
    c.Server.Port = 8080
  }

  // The webserver & base context path
  c.Server.server = rest.NewServer( c.Server.Port )
  c.Server.server.Headers = c.Server.Headers
  c.Server.server.Origins = c.Server.Origins
  c.Server.server.Methods = c.Server.Methods

  c.Server.Ctx = c.Server.server.Context( c.Server.Context )

  return nil
}

func (c *Config) Start() error {
  c.Cron.Start()

  return c.Server.server.Start()
}
