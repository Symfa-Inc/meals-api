package dev

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/utils"
	"sync"
)

// CreateClients creates seeds for clients table
func CreateClients() {
	seedExists := config.DB.Where("name = ?", "init clients").First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init clients",
		}

		var clientsArray []domain.Client
		utils.JSONParse("/db/seeds/data/clients.json", &clientsArray)
		var cateringsArray []domain.Catering
		config.DB.Find(&cateringsArray).Limit(len(clientsArray))

		var wg sync.WaitGroup
		wg.Add(len(clientsArray))

		for i := range clientsArray {
			go func(i int) {
				defer wg.Done()
				clientsArray[i].CateringID = cateringsArray[i].ID
				config.DB.Create(&clientsArray[i])
			}(i)
		}

		wg.Wait()
		config.DB.Create(&seed)
		fmt.Println("=== Clients seeds created ===")
	} else {
		fmt.Printf("Seed `init clients` already exists \n")
	}
}
