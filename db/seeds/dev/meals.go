package dev

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/repository"
	"sync"
	"time"
)

func CreateMeals() {
	cateringRepo := repository.NewCateringRepo()
	seedExist := config.DB.Where("name = ?", "init meals").First(&domain.Seed{}).Error
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
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
				meal := domain.Meal{
					Date:       time.Now().AddDate(0, 0, i).Truncate(t),
					CateringID: cateringResult.ID,
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
