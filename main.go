package main

import (
	_ "github.com/Aiscom-LLC/meals-api/docs"
	"github.com/Aiscom-LLC/meals-api/src/delivery/api"
	_ "github.com/Aiscom-LLC/meals-api/src/init"
)

// @title AIS Catering
// @version 1.0.0
func main() {
	r := api.SetupRouter()
	if err := r.Run(); err != nil {
		panic(err)
	}
}
