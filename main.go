package main

import (
	_ "go_api/docs"
	"go_api/src/config"
	"go_api/src/delivery"
)

// @title AIS Catering
// @version 1.0.0
func main() {
	config.CRON.Cron.Start()
	r := delivery.SetupRouter()
	r.Run()
}
