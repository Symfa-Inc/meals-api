package dev

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/models/seeds"
	"go_api/src/models/user"
	"go_api/src/types"
	"go_api/src/utils"
)

const seedName string = "init admin"

// CreateAdmin init admin user
func CreateAdmin() {
	seedExists := config.DB.Where("name = ?", seedName).First(&seeds.Seed{}).Error
	if seedExists != nil {
		seed := seeds.Seed{
			Name: seedName,
		}

		superAdmin := users.User{
			FirstName: "super",
			LastName:  "admin",
			Email:     "admin@meals.com",
			Password:  utils.HashString("Password12!"),
			Role:      types.UserRoleEnum.SuperAdmin,
		}

		config.DB.Create(&superAdmin)
		config.DB.Create(&seed)
		fmt.Println("=== Admin seed created ===")
	} else {
		fmt.Printf("Seed `%s` already exists \n", seedName)
	}
}
