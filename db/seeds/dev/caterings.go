package dev

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/models"
	"go_api/src/utils"
	"sync"
)

// CreateCaterings will populate table with random caterings
func CreateCaterings() {
	seedExists := config.DB.Where("name = ?", "init caterings").First(&models.Seed{}).Error
	if seedExists != nil {
		seed := models.Seed{
			Name: "init caterings",
		}

		var cateringArray []models.Catering
		utils.JsonParse("/db/seeds/data/caterings.json", &cateringArray)

		var wg sync.WaitGroup
		wg.Add(len(cateringArray))

		for i := range cateringArray {
			go func(i int) {
				defer wg.Done()
				config.DB.Create(&cateringArray[i])
			}(i)
		}

		wg.Wait()
		config.DB.Create(&seed)
		fmt.Println("=== Caterings seeds created ===")
	} else {
		fmt.Printf("Seed `init caterings` already exists \n")
	}
}
