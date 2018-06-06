package bin

import (
  "gopkg.in/robfig/cron.v2"
)

// initCron initialises the Cron scheduler
func (c *Config) initCron() error {

  c.Cron = cron.New()

  return nil
}
