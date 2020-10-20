package dev

import (
	"fmt"
	"time"

	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
)

// CreateMealDishes creates seeds for meal_dishes table
func CreateMealDishes() {
	cateringRepo := repository.NewCateringRepo()
	categoryRepo := repository.NewCategoryRepo()
	dishRepo := repository.NewDishRepo()
	mealRepo := repository.NewMealRepo()

	seedExists := config.DB.
		Where("name = ?", "init meals dishes").
		First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init meal dishes",
		}

		var dishesArray []domain.Dish

		cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
		cateringID := cateringResult.ID.String()

		categoryResult, _ := categoryRepo.GetByKey("name", "супы", cateringID)

		categoryID := categoryResult.ID.String()

		t := time.Hour * 24
		dishesArray, _, _ = dishRepo.Get(cateringID, categoryID)

		mealResult, _, _ := mealRepo.GetByKey("date", time.Now().AddDate(0, 0, 0).Truncate(t).Format(time.RFC3339))
		var mealDish domain.MealDish
		for i := range dishesArray {
			mealDish.DishID = dishesArray[i].ID
			mealDish.MealID = mealResult.ID
			config.DB.Create(&mealDish)
		}

		config.DB.Create(&seed)
		fmt.Println("=== Meal Dishes seeds created ===")
	} else {
		fmt.Printf("Seed `init meal dishes` already exists \n")
	}
}
