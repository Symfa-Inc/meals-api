package config

import "github.com/robfig/cron"

// CRON cron struct
var CRON struct {
	Cron    *cron.Cron
	Entries []map[string][]cron.EntryID
}

func init() {
	CRON.Cron = cron.New()
}
