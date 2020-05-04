package dev

import (
	"fmt"
	"go_api/db/seeds/data"
	"go_api/src/config"
	"go_api/src/models/seeds"
)

// CreateUsers will populate users table with random users
func CreateUsers() {
	seedExists := config.DB.Where("name = ?", "init users").First(&seeds.Seed{}).Error
	if seedExists != nil {
		seed := seeds.Seed{
			Name: "init users",
		}
		users := data.CreateUsersData(50)
		for i := range users {
			config.DB.Create(&users[i])
		}
		config.DB.Create(&seed)
		fmt.Println("=== User seeds created ===")
	} else {
		fmt.Printf("Seed `init users` already exists")
	}
}
