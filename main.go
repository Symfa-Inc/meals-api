package main

import (
	"github.com/Aiscom-LLC/meals-api/api"
	_ "github.com/Aiscom-LLC/meals-api/docs"
	_ "github.com/Aiscom-LLC/meals-api/init"
)

// @title AIS Catering
// @version 1.0.0
func main() {
	r := api.SetupRouter()
	if err := r.Run(); err != nil {
		panic(err)
	}
}
