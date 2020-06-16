package dev

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/models"
	"go_api/src/repository/catering"
	"sync"
	"time"
)

func CreateMeals() {
	seedExist := config.DB.Where("name = ?", "init meals").First(&models.Seed{}).Error
	cateringResult, _ := catering.GetCateringByKey("name", "Twiist")
	if seedExist != nil {
		seed := models.Seed{
			Name: "init meals",
		}
		var wg sync.WaitGroup
		wg.Add(7)

		for i := 0; i < 7; i++ {
			go func(i int) {
				defer wg.Done()
				meal := models.Meal{
					Date:       time.Now().AddDate(0, 0, i),
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
