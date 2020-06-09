package main

import (
	_ "go_api/docs"
	_ "go_api/src/config"
	"go_api/src/delivery"
)

// @title AIS Catering
// @version 1.0.0
func main() {
	r := delivery.SetupRouter()
	r.Run()
}
