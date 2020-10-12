package config

import (
	"github.com/robfig/cron"
)

// CRON cron struct
var CRON struct {
	Cron    *cron.Cron
	Entries []map[string]map[int]string
}

func init() {
	CRON.Cron = cron.New()
}
