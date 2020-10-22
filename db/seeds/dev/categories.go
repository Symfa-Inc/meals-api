package dev

import (
	"fmt"
	"sync"

	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/utils"
)

// CreateCategories creates seeds for categories table
func CreateCategories() {
	cateringRepo := repository.NewCateringRepo()
	seedExists := config.DB.
		Where("name = ?", "init dish_categories").
		First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init dish_categories",
		}

		catering, _ := cateringRepo.GetByKey("name", "Twiist")
		var categoriesArray []domain.Category
		utils.JSONParse("/db/seeds/data/categories.json", &categoriesArray)

		var wg sync.WaitGroup
		wg.Add(len(categoriesArray))

		for i := range categoriesArray {
			go func(i int) {
				defer wg.Done()
				categoriesArray[i].CateringID = catering.ID
				config.DB.Create(&categoriesArray[i])
			}(i)
		}

		wg.Wait()
		config.DB.Create(&seed)
		fmt.Println("=== Dish Categories seeds created ===")
	} else {
		fmt.Printf("Seed `init dish_categories` already exists \n")
	}
}
