package main

import (
	"io"
	"os"
	"time"

	"github.com/Aiscom-LLC/meals-api/src/delivery"
	_ "github.com/Aiscom-LLC/meals-api/src/init"
	"github.com/gin-gonic/gin"
)

// @title AIS Catering
// @version 1.0.0
func main() {
	os.Mkdir("logs", 0777)
	f, _ := os.Create("logs/gin.log: " + time.Now().UTC().String())
	gin.DefaultWriter = io.MultiWriter(f)

	r := delivery.SetupRouter()
	if err := r.Run(); err != nil {
		panic(err)
	}
}
