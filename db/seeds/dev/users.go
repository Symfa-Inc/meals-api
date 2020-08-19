package dev

import (
	"fmt"

	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/types"
	"github.com/Aiscom-LLC/meals-api/src/utils"
)

// CreateUsers will populate users table with random users
func CreateUsers() {
	seedExists := config.DB.Where("name = ?", "init users").First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init users",
		}

		hashedPassword := utils.HashString("Password12!")

		var userArray []domain.User
		var clientArray []domain.Client
		utils.JSONParse("/db/seeds/data/users.json", &userArray)
		utils.JSONParse("/db/seeds/data/clients.json", &clientArray)

		for i := range userArray {
			if i < 3 {
				userArray[i].Password = hashedPassword
				userArray[i].Status = &types.StatusTypesEnum.Active
				config.DB.Create(&userArray[i])
			} else {
				userArray[i].Password = hashedPassword
				userArray[i].Status = &types.StatusTypesEnum.Active
				config.DB.Create(&userArray[i])
			}
		}
		config.DB.Create(&seed)

		fmt.Println("=== User seeds created ===")
	} else {
		fmt.Printf("Seed `init users` already exists \n")
	}
}
