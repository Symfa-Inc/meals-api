package usecase

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/schemes/request"
	"github.com/Aiscom-LLC/meals-api/types"
	"github.com/Aiscom-LLC/meals-api/utils"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Meal struct
type Meal struct{}

// NewMeal returns pointer to meal struct
// with all methods
func NewMeal() *Meal {
	return &Meal{}
}

var mealRepo = repository.NewMealRepo()
var mealDishRepo = repository.NewMealDishesRepo()

// Add Creates meal for certain client
// @Summary Creates meal for certain client
// @Tags catering meals
// @Produce json
// @Param id path string false "Catering ID"
// @Param clientId path string false "Client ID"
// @Param payload body request.AddMeal false "meal reading"
// @Success 201 {object} request.AddMeal "created meal"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/clients/{clientId}/meals [post]
func (m Meal) Add(c *gin.Context) {
	var path types.PathClient
	var body request.AddMeal

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	ctxUser, _ := c.Get("user")
	ctxUserName := ctxUser.(domain.User).FirstName + " " + ctxUser.(domain.User).LastName

	parsedCateringID, _ := uuid.FromString(path.ID)
	parsedClientID, _ := uuid.FromString(path.ClientID)
	meal := &domain.Meal{
		Date:       body.Date,
		CateringID: parsedCateringID,
		ClientID:   parsedClientID,
		Person:     ctxUserName,
	}

	t := 24 * time.Hour
	difference := body.Date.Sub(time.Now().Truncate(t)).Hours()

	if difference < 0 {
		utils.CreateError(http.StatusBadRequest, errors.New("item has wrong date (can't use previous dates)"), c)
		return
	}

	meals, code, err := mealRepo.Get(body.Date, path.ID, path.ClientID)

	if err != nil {
		utils.CreateError(code, err, c)
		return
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
			utils.CreateError(code, err, c)
			return
		}
	}

	if err := mealRepo.Add(meal); err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	for _, dishID := range body.Dishes {
		dishIDParsed, _ := uuid.FromString(dishID)
		mealDish := domain.MealDish{
			MealID: meal.ID,
			DishID: dishIDParsed,
		}
		if err := mealDishRepo.Add(mealDish); err != nil {
			utils.CreateError(http.StatusBadRequest, err, c)
			return
		}
	}

	result, _, _ := mealRepo.Get(body.Date, path.ID, path.ClientID)

	c.JSON(http.StatusCreated, result)
}

// Get returns array of meals
// @Summary GetByRange list of categories with dishes for passed meal ID
// @Tags catering meals
// @Produce json
// @Param startDate query string false "Meal Start Date in 2020-01-01T00:00:00Z format"
// @Param endDate query string false "Meal End Date in 2020-01-01T00:00:00Z format"
// @Param id path string false "Catering ID"
// @Param clientId path string false "Client ID"
// @Success 200 {array} response.GetMeal "dishes for passed day"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/clients/{clientId}/meals [get]
func (m Meal) Get(c *gin.Context) {
	var query types.DateRangeQuery
	var path types.PathClient
	var cateringRepo = repository.NewCateringRepo()
	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	_, err := cateringRepo.GetByKey("id", path.ID)

	if err != nil {
		if err.Error() == "record not found" {
			utils.CreateError(http.StatusNotFound, err, c)
			return
		}
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	startDate, err := time.Parse(time.RFC3339, query.StartDate)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, errors.New("can't parse the date"), c)
		return
	}

	endDate, err := time.Parse(time.RFC3339, query.EndDate)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, errors.New("can't parse the date"), c)
		return
	}

	if startDate.After(endDate) {
		utils.CreateError(http.StatusBadRequest, errors.New("end date can't be earlier than start date"), c)
		return
	}

	result, code, err := mealRepo.GetByRange(startDate, endDate, path.ID, path.ClientID)
	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, result)
}
