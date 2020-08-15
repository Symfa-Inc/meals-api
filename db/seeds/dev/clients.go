package dev

import (
	"fmt"
	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/Aiscom-LLC/meals-api/src/utils"
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
		config.DB.Find(&cateringsArray).Limit(len(clientsArray))
		twiistCatering, _ := cateringRepo.GetByKey("name", "Twiist")

		for i := range clientsArray {
			clientsArray[i].AutoApproveOrders = false
			if i == 0 {
				clientsArray[i].CateringID = twiistCatering.ID
				config.DB.Create(&clientsArray[i])
			} else if cateringsArray[i].Name != "Twiist" {
				clientsArray[i].CateringID = cateringsArray[i].ID
				config.DB.Create(&clientsArray[i])
			}
		}

		config.DB.Create(&seed)
		fmt.Println("=== Clients seeds created ===")
	} else {
		fmt.Printf("Seed `init clients` already exists \n")
	}
}
