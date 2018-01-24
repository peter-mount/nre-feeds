package main

import (
  "darwind3"
)

func (c *Config) initPushPort() error {

  if( c.PushPort.Enabled ) {
    c.PushPort.d3 = &darwind3.DarwinD3{}

    c.PushPort.d3.SetupRest( c.Server.ctx.Context( "/pushPort" ) )
  }

  return nil
}
