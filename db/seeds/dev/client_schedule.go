package dev

import (
	"fmt"
	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/domain"
	"sync"
)

// CreateClientSchedules will populate table with default client schedules
func CreateClientSchedules() {
	seedExists := config.DB.Where("name = ?", "init client schedules").First(&domain.Seed{}).Error
	if seedExists != nil {
		seed := domain.Seed{
			Name: "init client schedules",
		}

		var clients []domain.Client
		config.DB.Find(&clients)
		var wg sync.WaitGroup
		wg.Add(len(clients))
		for i := range clients {
			go func(i int) {
				defer wg.Done()
				var schedules []domain.CateringSchedule
				config.DB.Where("catering_id = ?", clients[i].CateringID).Find(&schedules)
				for _, schedule := range schedules {
					clientSchedule := domain.ClientSchedule{
						Day:       schedule.Day,
						Start:     schedule.Start,
						End:       schedule.End,
						IsWorking: schedule.IsWorking,
						ClientID:  clients[i].ID,
					}
					config.DB.Create(&clientSchedule)
				}
			}(i)
		}

		wg.Wait()
		config.DB.Create(&seed)
		fmt.Println("=== Client schedules seeds created ===")
	} else {
		fmt.Printf("Seed `init client schedules` already exists \n")
	}
}
