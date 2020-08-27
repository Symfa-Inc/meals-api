package dev

import (
	"fmt"

	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/types"
)

// CreateClientUsers will populate client users table
func CreateClientUsers() {
	seedExists := config.DB.Where("name = ?", "init client users").First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init client users",
		}

		clientRepo := repository.NewClientRepo()
		client, _ := clientRepo.GetByKey("name", "Dymi")

		userRepo := repository.NewUserRepo()
		adminUsers, _ := userRepo.GetAllByKey("role", types.UserRoleEnum.ClientAdmin)
		users, _ := userRepo.GetAllByKey("role", types.UserRoleEnum.User)

		for i, user := range users {
			clientUser := domain.ClientUser{
				ClientID: client.ID,
				UserID:   user.ID,
				Floor:    i + 1,
			}
			config.DB.Create(&clientUser)
		}

		for i, user := range adminUsers {
			clientUser := domain.ClientUser{
				ClientID: client.ID,
				UserID:   user.ID,
				Floor:    i + 1,
			}
			config.DB.Create(&clientUser)
		}

		config.DB.Create(&seed)

		fmt.Println("=== Catering user seeds created ===")
	} else {
		fmt.Printf("Seed `init catering users` already exists \n")
	}
}
