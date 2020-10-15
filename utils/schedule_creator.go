package utils

import (
	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/domain"
	uuid "github.com/satori/go.uuid"
)

// AddDefaultCateringSchedules adds default
// for provided catering id
func AddDefaultCateringSchedules(cateringID uuid.UUID) {
	for i := 0; i < 7; i++ {
		if i < 5 {
			schedule := domain.CateringSchedule{
				Day:        i,
				Start:      "00:00",
				End:        "16:45",
				IsWorking:  true,
				CateringID: cateringID,
			}
			config.DB.Create(&schedule)
		} else {
			schedule := domain.CateringSchedule{
				Day:        i,
				Start:      "00:00",
				End:        "16:45",
				IsWorking:  false,
				CateringID: cateringID,
			}
			config.DB.Create(&schedule)
		}
	}
}
