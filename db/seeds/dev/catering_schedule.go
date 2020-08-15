package dev

import (
	"fmt"
	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/utils"
	"sync"
)

// CreateCateringSchedules will populate table with default catering schedules
func CreateCateringSchedules() {
	seedExists := config.DB.Where("name = ?", "init catering schedules").First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init catering schedules",
		}

		var wg sync.WaitGroup
		var caterings []domain.Catering
		config.DB.Find(&caterings)
		wg.Add(len(caterings))
		for i := range caterings {
			go func(i int) {
				defer wg.Done()
				utils.AddDefaultCateringSchedules(caterings[i].ID)
			}(i)
		}
		wg.Wait()
		config.DB.Create(&seed)
		fmt.Println("=== Catering schedules seeds created ===")
	} else {
		fmt.Printf("Seed `init catering schedules` already exists \n")
	}
}
