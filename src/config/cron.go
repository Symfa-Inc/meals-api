package config

import "github.com/robfig/cron"

var CRON *cron.Cron

func init() {
	CRON = cron.New()
}
