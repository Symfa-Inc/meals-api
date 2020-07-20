package dev

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/utils"
	"sync"
)

func CreateClients() {
	seedExists := config.DB.Where("name = ?", "init clients").First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init clients",
		}

		var clientsArray []domain.Client
		utils.JsonParse("/db/seeds/data/clients.json", &clientsArray)

		var wg sync.WaitGroup
		wg.Add(len(clientsArray))

		for i := range clientsArray {
			go func(i int) {
				defer wg.Done()
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
