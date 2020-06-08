package dev

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/models"
	"go_api/src/utils"
	"sync"
)

// CreateUsers will populate users table with random users
func CreateUsers() {
	seedExists := config.DB.Where("name = ?", "init users").First(&models.Seed{}).Error
	if seedExists != nil {
		seed := models.Seed{
			Name: "init users",
		}

		hashedPassword := utils.HashString("Password12!")

		var userArray []models.User
		utils.JsonParse("/db/seeds/data/users.json", &userArray)

		var wg sync.WaitGroup
		wg.Add(len(userArray))

		for i := range userArray {
			go func(i int) {
				defer wg.Done()
				userArray[i].Password = hashedPassword
				config.DB.Create(&userArray[i])
			}(i)
		}

		wg.Wait()
		config.DB.Create(&seed)

		fmt.Println("=== User seeds created ===")
	} else {
		fmt.Printf("Seed `init users` already exists \n")
	}
}
