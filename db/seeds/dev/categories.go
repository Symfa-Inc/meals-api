package dev

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/repository"
	"go_api/src/utils"
	"sync"
)

func CreateCategories() {
	cateringRepo := repository.NewCateringRepo()
	seedExists := config.DB.
		Where("name = ?", "init dish_categories").
		First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init dish_categories",
		}

		cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
		var categoriesArray []domain.Category
		utils.JsonParse("/db/seeds/data/categories.json", &categoriesArray)

		var wg sync.WaitGroup
		wg.Add(len(categoriesArray))

		for i := range categoriesArray {
			go func(i int) {
				defer wg.Done()
				categoriesArray[i].CateringID = cateringResult.ID
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
