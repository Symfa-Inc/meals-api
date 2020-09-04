package services

import (
	"errors"
	"github.com/Aiscom-LLC/meals-api/api/swagger"
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
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
var mealDishRepo = repository.NewMealDishesRepo()

func (m *MealService) Add(path url.PathClient, body swagger.AddMeal, user interface{}) ([]swagger.GetMeal, int, error) {

	userName := user.(domain.User).FirstName + " " + user.(domain.User).LastName

	parsedCateringID, _ := uuid.FromString(path.ID)
	parsedClientID, _ := uuid.FromString(path.ClientID)
	meal := &domain.Meal{
		Date:       body.Date,
		CateringID: parsedCateringID,
		ClientID:   parsedClientID,
		Person:     userName,
	}

	t := 24 * time.Hour
	difference := body.Date.Sub(time.Now().Truncate(t)).Hours()

	if difference < 0 {
		return []swagger.GetMeal{}, http.StatusBadRequest, errors.New("item has wrong date (can't use previous dates)")
	}

	meals, code, err := mealRepo.Get(body.Date, path.ID, path.ClientID)

	if err != nil {
		return []swagger.GetMeal{}, code, err
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
			return []swagger.GetMeal{}, code, err
		}
	}

	if err := mealRepo.Add(meal); err != nil {
		return []swagger.GetMeal{}, code, err
	}

	for _, dishID := range body.Dishes {
		dishIDParsed, _ := uuid.FromString(dishID)
		mealDish := domain.MealDish{
			MealID: meal.ID,
			DishID: dishIDParsed,
		}
		if err := mealDishRepo.Add(mealDish); err != nil {
			return []swagger.GetMeal{}, http.StatusBadRequest, err
		}
	}

	result, code, err := mealRepo.Get(body.Date, path.ID, path.ClientID)

	return result, code, err
}

var cateringRepo = repository.NewCateringRepo()

func (m *MealService) Get(query url.DateRangeQuery, path url.PathClient) ([]swagger.GetMeal, int, error) {
	_, err := cateringRepo.GetByKey("id", path.ID)

	if err != nil {
		if err.Error() == "record not found" {
			return []swagger.GetMeal{}, http.StatusNotFound, err
		}
		return []swagger.GetMeal{}, http.StatusBadRequest, err
	}

	startDate, err := time.Parse(time.RFC3339, query.StartDate)
	if err != nil {
		return []swagger.GetMeal{}, http.StatusBadRequest, errors.New("can't parse the date")
	}

	endDate, err := time.Parse(time.RFC3339, query.EndDate)
	if err != nil {
		return []swagger.GetMeal{}, http.StatusBadRequest, errors.New("can't parse the date")
	}

	if startDate.After(endDate) {
		return []swagger.GetMeal{}, http.StatusBadRequest, errors.New("end date can't be earlier than start date")
	}

	result, code, err := mealRepo.GetByRange(startDate, endDate, path.ID, path.ClientID)

	return result, code, err
}
