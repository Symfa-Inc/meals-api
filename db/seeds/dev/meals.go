package dev

import (
	"fmt"
	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

// CreateMeals creates seeds for meals table
func CreateMeals() {
	cateringRepo := repository.NewCateringRepo()
	clientRepo := repository.NewClientRepo()
	seedExist := config.DB.Where("name = ?", "init meals").First(&domain.Seed{}).Error
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	client, _ := clientRepo.GetByKey("name", "Dymi")
	if seedExist != nil {
		seed := domain.Seed{
			Name: "init meals",
		}
		var wg sync.WaitGroup
		wg.Add(7)

		t := 24 * time.Hour
		for i := 0; i < 7; i++ {
			go func(i int) {
				defer wg.Done()
				mealID := uuid.NewV4()
				meal := domain.Meal{
					Date:       time.Now().AddDate(0, 0, i).Truncate(t),
					CateringID: cateringResult.ID,
					ClientID:   client.ID,
					MealID:     mealID,
					Version:    "V.1",
					Person:     "Super Admin",
				}
				config.DB.Create(&meal)
			}(i)
		}

		wg.Wait()
		config.DB.Create(&seed)
		fmt.Println("=== Meals seeds created ===")
	} else {
		fmt.Printf("Seed `init caterings` already exists \n")
	}
}
