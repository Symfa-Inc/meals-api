package dev

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/repository"
	"go_api/src/utils"
)

// CreateClients creates seeds for clients table
func CreateClients() {
	cateringRepo := repository.NewCateringRepo()
	seedExists := config.DB.Where("name = ?", "init clients").First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init clients",
		}

		var clientsArray []domain.Client
		utils.JSONParse("/db/seeds/data/clients.json", &clientsArray)
		var cateringsArray []domain.Catering
		twiistCatering, _ := cateringRepo.GetByKey("name", "Twiist")
		config.DB.Find(&cateringsArray).Limit(len(clientsArray))

		for i := range clientsArray {
			if i == 0 {
				clientsArray[i].CateringID = twiistCatering.ID
				config.DB.Create(&clientsArray[i])
			} else {
				if cateringsArray[i].Name != "Twiist" {
					clientsArray[i].CateringID = cateringsArray[i].ID
					config.DB.Create(&clientsArray[i])
				}
			}
		}

		config.DB.Create(&seed)
		fmt.Println("=== Clients seeds created ===")
	} else {
		fmt.Printf("Seed `init clients` already exists \n")
	}
}
