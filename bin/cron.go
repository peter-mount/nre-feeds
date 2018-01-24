package main

import (
  "gopkg.in/robfig/cron.v2"
)

func (c *Config) initCron() error {

  c.cron = cron.New()

  return nil
}
