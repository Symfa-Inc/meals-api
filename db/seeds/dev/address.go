package dev

import (
	"fmt"
	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/utils"
	"sync"
)

// CreateAddresses creates seeds for addresses table
func CreateAddresses() {
	seedExists := config.DB.
		Where("name = ?", "init addresses").
		First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init addresses",
		}

		var addressesArray []domain.Address
		utils.JSONParse("/db/seeds/data/addresses.json", &addressesArray)
		var clientsArray []domain.Client
		config.DB.Find(&clientsArray).Limit(len(addressesArray))

		var wg sync.WaitGroup
		wg.Add(len(addressesArray))

		for i := range addressesArray {
			go func(i int) {
				defer wg.Done()
				addressesArray[i].ClientID = clientsArray[i].ID
				config.DB.Create(&addressesArray[i])
			}(i)
		}

		wg.Wait()
		config.DB.Create(&seed)
		fmt.Println("=== Addresses seeds created ===")
	} else {
		fmt.Printf("Seed `init addresses` already exists \n")
	}
}
