package dev

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/repository"
	"time"
)

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
		var dishesArray2 []domain.Dish

		cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
		cateringId := cateringResult.ID.String()

		categoryResult, _ := categoryRepo.GetByKey("name", "супы", cateringId)
		categoryResult2, _ := categoryRepo.GetByKey("name", "гарнир", cateringId)

		categoryId := categoryResult.ID.String()
		categoryId2 := categoryResult2.ID.String()

		t := time.Hour * 24
		dishesArray, _, _ = dishRepo.Get(cateringId, categoryId)
		dishesArray2, _, _ = dishRepo.Get(cateringId, categoryId2)

		mealResult, _, _ := mealRepo.GetByKey("date", time.Now().AddDate(0, 0, 0).Truncate(t).Format(time.RFC3339))
		var mealDish domain.MealDish
		for i := range dishesArray {
			mealDish.DishID = dishesArray[i].ID
			mealDish.MealID = mealResult.ID
			config.DB.Create(&mealDish)
		}

		for i := range dishesArray {
			mealDish.DishID = dishesArray2[i].ID
			mealDish.MealID = mealResult.ID
			config.DB.Create(&mealDish)
		}

		config.DB.Create(&seed)
		fmt.Println("=== Meal Dishes seeds created ===")
	} else {
		fmt.Printf("Seed `init meal dishes` already exists \n")
	}
}
