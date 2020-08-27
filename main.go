package main

import (
	_ "github.com/Aiscom-LLC/meals-api/docs"
	"github.com/Aiscom-LLC/meals-api/delivery"
	_ "github.com/Aiscom-LLC/meals-api/init"
)

// @title AIS Catering
// @version 1.0.0
func main() {
	r := delivery.SetupRouter()
	if err := r.Run(); err != nil {
		panic(err)
	}
}
