package main

import (
	_ "go_api/docs"
	"go_api/src/delivery"
	_ "go_api/src/init"
)

// @title AIS Catering
// @version 1.0.0
func main() {
	r := delivery.SetupRouter()
	r.Run()
}
