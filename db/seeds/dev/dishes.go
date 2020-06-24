package dev

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/repository"
	"go_api/src/utils"
	"sync"
)

func CreateDishes() {
	cateringRepo := repository.NewCateringRepo()
	dishCategoryRepo := repository.NewDishCategoryRepo()
	seedExists := config.DB.
		Where("name = ?", "init dishes").
		First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init dishes",
		}

		var dishesArray []domain.Dish
		cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
		dishCategoryResult, _ := dishCategoryRepo.GetByKey("name", "супы", cateringResult.ID.String())
		utils.JsonParse("/db/seeds/data/dishes.json", &dishesArray)

		var wg sync.WaitGroup
		wg.Add(len(dishesArray))

		for i := range dishesArray {
			go func(i int) {
				defer wg.Done()
				dishesArray[i].CateringID = cateringResult.ID
				dishesArray[i].DishCategoryID = dishCategoryResult.ID
				config.DB.Create(&dishesArray[i])
			}(i)
		}

		wg.Wait()
		config.DB.Create(&seed)
		fmt.Println("=== Dishes seeds created ===")
	} else {
		fmt.Printf("Seed `init dishes` already exists \n")
	}
}
