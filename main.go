package main

import (
	"github.com/Aiscom-LLC/meals-api/src/delivery"
	_ "github.com/Aiscom-LLC/meals-api/src/init"
)

// @title AIS Catering
// @version 1.0.0
func main() {

	r := delivery.SetupRouter()
	if err := r.Run(); err != nil {
		panic(err)
	}
}
