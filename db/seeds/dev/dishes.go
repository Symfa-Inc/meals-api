package dev

import (
	"fmt"
	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/Aiscom-LLC/meals-api/src/utils"
)

// CreateDishes creates seeds for clients table
func CreateDishes() {
	cateringRepo := repository.NewCateringRepo()
	categoryRepo := repository.NewCategoryRepo()
	seedExists := config.DB.
		Where("name = ?", "init dishes").
		First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init dishes",
		}

		var dishesArray []domain.Dish
		cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
		categoryResult, _ := categoryRepo.GetByKey("name", "супы", cateringResult.ID.String())
		utils.JSONParse("/db/seeds/data/dishes.json", &dishesArray)

		for i := range dishesArray {
			dishesArray[i].CateringID = cateringResult.ID
			dishesArray[i].CategoryID = categoryResult.ID
			config.DB.Create(&dishesArray[i])
		}

		categoryResult2, _ := categoryRepo.GetByKey("name", "гарнир", cateringResult.ID.String())
		for i := range dishesArray {
			dishesArray[i].CateringID = cateringResult.ID
			dishesArray[i].CategoryID = categoryResult2.ID
			config.DB.Create(&dishesArray[i])
		}

		config.DB.Create(&seed)
		fmt.Println("=== Dishes seeds created ===")
	} else {
		fmt.Printf("Seed `init dishes` already exists \n")
	}
}
