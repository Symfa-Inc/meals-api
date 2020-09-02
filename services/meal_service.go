package services

import (
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/schemes/request"
	"github.com/Aiscom-LLC/meals-api/types"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
	"time"
)

// MealService struct
type MealService struct{}

// NewMealService return pointer to meal struct
// with all methods
func NewMealService() *MealService {
	return &MealService{}
}

var mealRepo = repository.NewMealRepo()
var dishRepo = repository.NewDishRepo()

func (m *MealService) Add(path types.PathClient, body request.AddMeal, user interface{}) ([]domain.Meal, int, error) {

	userName := user.(domain.User).FirstName + " " + user.(domain.User).LastName

	parsedCateringID, _ := uuid.FromString(path.ID)
	parsedClientID, _ := uuid.FromString(path.ClientID)
	meal := &domain.Meal{
		Date: body.Date,
		CateringID: parsedCateringID,
		ClientID: parsedClientID,
		Person: userName,
	}

	t := 24 * time.Hour
	difference := body.Date.Sub(time.Now().Truncate(t)).Hours()

	if difference < 0 {
		return []domain.Meal{}, http.StatusBadRequest, nil
	}

	meals, code, err := mealRepo.Get(body.Date, path.ID, path.ClientID)

	if err != nil {
		return []domain.Meal{}, code, err
	}

	if len(meals) != 0 {
		meal.MealID = meals[0].MealID
		meal.Version = "V." + strconv.Itoa(len(meals)+1)
	} else {
		MealID := uuid.NewV4()
		meal.MealID = MealID
		meal.Version = "V.1"
	}

	for _, dishID := range body.Dishes {
		_, code, err := dishRepo.FindByID(path.ID, dishID)
		if err != nil {
			return []domain.Meal{}, code, err
		}
	}

	if err := mealRepo.Add(meal); err != nil {
		return []domain.Meal{}, code, err
	}

	for _, dishID := range body.Dishes {
		dishIDParsed, _ := uuid.FromString(dishID)
		mealDish := domain.MealDish{
			MealID: meal.ID,
			DishID: dishIDParsed,
		}
		if err
	}

	return []domain.Meal{}, 0, nil
}