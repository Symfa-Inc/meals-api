package dev

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/repository"
	"go_api/src/types"
	"go_api/src/utils"
)

// CreateUsers will populate users table with random users
func CreateUsers() {
	seedExists := config.DB.Where("name = ?", "init users").First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init users",
		}

		hashedPassword := utils.HashString("Password12!")
		cateringRepo := repository.NewCateringRepo()
		catering, _ := cateringRepo.GetByKey("name", "Twiist")

		clientRepo := repository.NewClientRepo()
		client, _ := clientRepo.GetByKey("name", "Dymi")

		var userArray []domain.User
		var clientArray []domain.Client
		utils.JSONParse("/db/seeds/data/users.json", &userArray)
		utils.JSONParse("/db/seeds/data/clients.json", &clientArray)

		for i := range userArray {
			if i < 3 {
				cateringClient, _ := clientRepo.GetByKey("name", clientArray[i].Name)
				userArray[i].CompanyType = &types.CompanyTypesEnum.Catering
				userArray[i].CateringID = &cateringClient.CateringID
				userArray[i].ClientID = &cateringClient.ID
				userArray[i].Password = hashedPassword
				userArray[i].Status = &types.StatusTypesEnum.Active
				config.DB.Create(&userArray[i])
			} else {
				userArray[i].CompanyType = &types.CompanyTypesEnum.Client
				userArray[i].Floor = &i
				userArray[i].CateringID = &catering.ID
				userArray[i].ClientID = &client.ID
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
