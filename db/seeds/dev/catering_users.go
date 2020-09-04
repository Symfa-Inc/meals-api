package dev

import (
	"fmt"

	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/repository/enums"
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
		users, _ := userRepo.GetAllByKey("role", enums.UserRoleEnum.CateringAdmin)

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
