package dev

import (
	"fmt"
	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/utils"
)

// CreateCaterings will populate table with random caterings
func CreateCaterings() {
	seedExists := config.DB.Where("name = ?", "init caterings").First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init caterings",
		}

		var cateringArray []domain.Catering
		utils.JSONParse("/db/seeds/data/caterings.json", &cateringArray)

		for i := range cateringArray {
			config.DB.Create(&cateringArray[i])
		}

		config.DB.Create(&seed)
		fmt.Println("=== Caterings seeds created ===")
	} else {
		fmt.Printf("Seed `init caterings` already exists \n")
	}
}
