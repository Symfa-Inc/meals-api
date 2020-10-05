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
	file, _ := os.Create("logs/gin: " + time.Now().UTC().String() + ".log")
	fileErr, _ := os.Create("logs/err: " + time.Now().UTC().String() + ".log")
	gin.DefaultWriter = io.MultiWriter(file)
	gin.DefaultErrorWriter = io.MultiWriter(fileErr)

	r := delivery.SetupRouter()
	if err := r.Run(); err != nil {
		panic(err)
	}
}
