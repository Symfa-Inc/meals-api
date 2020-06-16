package meal

import (
	"errors"
	"github.com/satori/go.uuid"
	"go_api/src/config"
	"go_api/src/models"
	"go_api/src/types"
	"time"
)

const layout = "2006-01-02T15:04:05"

func CreateMealsDB(query types.StartEndDateQuery, id string) ([]models.Meal, error)   {
	timeStart, _ := time.Parse(layout, query.StartDate)
	timeEnd, _ := time.Parse(layout, query.EndDate)

	days := timeEnd.Sub(timeStart).Hours() / 24

	if days < 0 {
		return nil, errors.New("EndDate is earlier than StartDate")
	}

	parsedId, _ := uuid.FromString(id)

	var mealsArray []models.Meal

	for i := 0; i <= int(days); i++ {
		meal := models.Meal{
			Date:       timeStart.AddDate(0,0, i),
			CateringID: parsedId,
		}
		config.DB.Create(&meal)
		mealsArray = append(mealsArray, meal)
	}
	return mealsArray, nil
}

//func GetMealsDB(query types.StartEndDateQuery, id string) ([]models.Meal, error)

