package bin

import (
  "gopkg.in/robfig/cron.v2"
)

func (c *Config) InitCron() error {

  c.Cron = cron.New()

  return nil
}
