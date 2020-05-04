package main

import (
	_ "go_api/docs"
	_ "go_api/src/config"
	"go_api/src/mux"
)
// @title AIS Catering
// @version 1.0.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	r := mux.SetupRouter()
	r.Run()
}
