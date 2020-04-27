package config

import (
	"os"
	"log"
	"github.com/joho/godotenv"
)

type env struct {
	DbHost string
	DbPort string
	DbUser string
	DbPassword string
	DbName string
}

// Env is env procjet struct
var Env env

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

 	Env = env{
		DbHost: os.Getenv("DB_HOST"),
		DbPort: os.Getenv("DB_PORT"),
		DbUser: os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName: os.Getenv("DB_NAME"), 
	}
}