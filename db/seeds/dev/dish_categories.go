package dev

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/models"
	"go_api/src/repository"
	"go_api/src/utils"
	"sync"
)

func CreateDishCategories() {
	seedExists := config.DB.
		Where("name = ?", "init dish_categories").
		First(&models.Seed{}).Error
	if seedExists != nil {
		seed := models.Seed{
			Name: "init dish_categories",
		}

		cateringResult, _ := repository.GetCateringByKey("name", "Twiist")
		var categoriesArray []models.DishCategory
		utils.JsonParse("/db/seeds/data/dish_categories.json", &categoriesArray)

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
