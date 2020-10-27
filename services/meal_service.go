package services

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	uuid "github.com/satori/go.uuid"
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

func (m *MealService) Add(path url.PathClient, body models.AddMeal, user interface{}) ([]models.GetMeal, int, error) {

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
		return []models.GetMeal{}, http.StatusBadRequest, errors.New("item has wrong date (can't use previous dates)")
	}

	meals, code, err := mealRepo.Get(body.Date, path.ID, path.ClientID)

	if err != nil {
		return []models.GetMeal{}, code, err
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
			return []models.GetMeal{}, code, err
		}
	}

	if err := mealRepo.Add(meal); err != nil {
		return []models.GetMeal{}, code, err
	}

	for _, dishID := range body.Dishes {
		dishIDParsed, _ := uuid.FromString(dishID)
		mealDish := domain.MealDish{
			MealID: meal.ID,
			DishID: dishIDParsed,
		}
		if err := mealDishRepo.Add(mealDish); err != nil {
			return []models.GetMeal{}, http.StatusBadRequest, err
		}
	}

	result, code, err := mealRepo.Get(body.Date, path.ID, path.ClientID)

	return result, code, err
}

var cateringRepo = repository.NewCateringRepo()

func (m *MealService) Get(query url.DateQuery, path url.PathClient) ([]models.GetMeal, int, error) {
	_, err := cateringRepo.GetByKey("id", path.ID)

	if err != nil {
		if err.Error() == "record not found" {
			return []models.GetMeal{}, http.StatusNotFound, err
		}
		return []models.GetMeal{}, http.StatusBadRequest, err
	}

	mealDate, err := time.Parse(time.RFC3339, query.Date)
	if err != nil {
		return []models.GetMeal{}, http.StatusBadRequest, errors.New("can't parse the date")
	}

	result, code, err := mealRepo.Get(mealDate, path.ID, path.ClientID)

	return result, code, err
}

func (m *MealService) CopyMeals(path url.PathClient, body models.CopyMealToDate) ([]models.GetMeal, int, error) {
	meals, code, err := mealRepo.Get(body.Date, path.ID, path.ClientID)

	if err != nil {
		return []models.GetMeal{}, code, err
	}

	mealExist, _, _ := mealRepo.Get(body.ToDate, path.ID, path.ClientID)
	if len(mealExist) != 0 {
		return []models.GetMeal{}, http.StatusBadRequest, errors.New("meals for current day already exist")
	}

	for meal := range meals {
		for dish := range meals[meal].Result {
			_, code, err := dishRepo.FindByID(path.ID, meals[meal].Result[dish].ID.String())

			if err != nil {
				return []models.GetMeal{}, code, err
			}
		}

		mealID := meals[meal].MealID
		mealResult, _, _ := mealRepo.GetByKey("meal_id", mealID.String())
		mealResult.Date = body.ToDate
		mealResult.MealID = uuid.NewV4()

		if err := mealRepo.Add(&mealResult); err != nil {
			return []models.GetMeal{}, code, err
		}

		for dish := range meals[meal].Result {
			meals[meal].Result[dish].ID = uuid.NewV4()

			if err := dishRepo.Add(meals[meal].Result[dish].CateringID.String(), &meals[meal].Result[dish]); err != nil {
				return []models.GetMeal{}, http.StatusBadRequest, err
			}

			mealDish := domain.MealDish{
				MealID: mealResult.ID,
				DishID: meals[meal].Result[dish].ID,
			}

			if err := mealDishRepo.Add(mealDish); err != nil {
				return []models.GetMeal{}, http.StatusBadRequest, err
			}
		}
	}

	result, code, err := mealRepo.Get(body.ToDate, path.ID, path.ClientID)

	return result, code, err
}

func (m *MealService) CopyWeek(path url.PathClient, body models.CopyMealToWeek) (int, error) {
	if len(body.Date) != len(body.ToWeek) {
		return http.StatusBadRequest, errors.New("non valid data")
	}

	today := time.Now()
	threeWeeks := float64(504) // 504 hours in 3 weeks in total
	if body.ToWeek[4].Sub(today).Hours() > threeWeeks {
		return http.StatusBadRequest, errors.New("can't use data after 3 week")
	}
	for date := range body.Date {
		meals, code, err := mealRepo.Get(body.Date[date], path.ID, path.ClientID)

		if err != nil {
			return code, err
		}

		mealExist, _, _ := mealRepo.Get(body.ToWeek[date], path.ID, path.ClientID)

		if len(mealExist) != 0 {
			return http.StatusBadRequest, fmt.Errorf("meals for %v already exist", body.ToWeek[date].Weekday())
		}

		for meal := range meals {
			mealID := meals[meal].MealID
			mealResult, _, _ := mealRepo.GetByKey("meal_id", mealID.String())
			mealResult.Date = body.ToWeek[date]
			mealResult.MealID = uuid.NewV4()

			if err := mealRepo.Add(&mealResult); err != nil {
				return code, err
			}

			for dish := range meals[meal].Result {
				meals[meal].Result[dish].ID = uuid.NewV4()
				cateringID := meals[meal].Result[dish].CateringID.String()
				currentDish := meals[meal].Result[dish]

				if err := dishRepo.Add(cateringID, &currentDish); err != nil {
					return http.StatusBadRequest, err
				}

				mealDish := domain.MealDish{
					MealID: mealResult.ID,
					DishID: meals[meal].Result[dish].ID,
				}

				if err := mealDishRepo.Add(mealDish); err != nil {
					return http.StatusBadRequest, err
				}
			}
		}
	}

	return http.StatusCreated, nil
}
