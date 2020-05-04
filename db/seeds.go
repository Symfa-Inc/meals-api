package main

import (
	"go_api/db/seeds/dev"
)

func main() {
	dev.CreateAdmin()
	dev.CreateUsers()
}