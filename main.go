package main

import (
	"fmt"
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
	if err := os.Mkdir("logs", 0777); err != nil {
		fmt.Println(err)
	}
	f, _ := os.Create("logs/gin.log: " + time.Now().UTC().String())
	gin.DefaultWriter = io.MultiWriter(f)

	r := delivery.SetupRouter()
	if err := r.Run(); err != nil {
		panic(err)
	}
}
