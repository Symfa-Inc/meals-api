package dev

import (
	"fmt"

	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/Aiscom-LLC/meals-api/src/types"
)

// CreateCateringUsers will populate catering users table
func CreateCateringUsers() {
	seedExists := config.DB.Where("name = ?", "init catering users").First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init catering users",
		}

		cateringRepo := repository.NewCateringRepo()
		catering, _ := cateringRepo.GetByKey("name", "Twiist")

		userRepo := repository.NewUserRepo()
		users, _ := userRepo.GetAllByKey("role", types.UserRoleEnum.CateringAdmin)

		for _, user := range users {
			cateringUser := domain.CateringUser{
				CateringID: catering.ID,
				UserID:     user.ID,
			}
			config.DB.Create(&cateringUser)
		}
		config.DB.Create(&seed)

		fmt.Println("=== Catering user seeds created ===")
	} else {
		fmt.Printf("Seed `init catering users` already exists \n")
	}
}
