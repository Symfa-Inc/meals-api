package dev

import (
	"fmt"
	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/Aiscom-LLC/meals-api/src/utils"
	"sync"
)

// CreateCategories creates seeds for categories table
func CreateCategories() {
	cateringRepo := repository.NewCateringRepo()
	clientRepo := repository.NewClientRepo()
	seedExists := config.DB.
		Where("name = ?", "init dish_categories").
		First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init dish_categories",
		}

		catering, _ := cateringRepo.GetByKey("name", "Twiist")
		client, _ := clientRepo.GetByKey("name", "Dymi")
		var categoriesArray []domain.Category
		utils.JSONParse("/db/seeds/data/categories.json", &categoriesArray)

		var wg sync.WaitGroup
		wg.Add(len(categoriesArray))

		for i := range categoriesArray {
			go func(i int) {
				defer wg.Done()
				categoriesArray[i].CateringID = catering.ID
				categoriesArray[i].ClientID = client.ID
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
