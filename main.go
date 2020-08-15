package main

import (
	_ "github.com/Aiscom-LLC/meals-api/docs"
	"github.com/Aiscom-LLC/meals-api/src/delivery"
	_ "github.com/Aiscom-LLC/meals-api/src/init"
)

// @title AIS Catering
// @version 1.0.0
func main() {
	r := delivery.SetupRouter()
	r.Run()
}
