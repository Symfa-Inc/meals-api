package usecase

import (
	"go_api/src/domain"
	"go_api/src/repository"
	"go_api/src/schemes/request"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
	"strconv"
	"time"

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
	ctxUserName := ctxUser.(domain.UserClientCatering).FirstName + " " + ctxUser.(domain.UserClientCatering).LastName

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
		utils.CreateError(http.StatusBadRequest, "item has wrong date (can't use previous dates)", c)
		return
	}

	meals, code, err := mealRepo.Get(body.Date, path.ID, path.ClientID)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	if len(meals) != 0 {
		meal.MealID = meals[0].MealID
		meal.Version = "V." + strconv.Itoa(len(meals)+1)
	} else {
		MealID, _ := uuid.NewV4()
		meal.MealID = MealID
		meal.Version = "V.1"
	}

	for _, dishID := range body.Dishes {
		_, code, err := dishRepo.FindByID(path.ID, dishID)
		if err != nil {
			utils.CreateError(code, err.Error(), c)
			return
		}
	}

	if err := mealRepo.Add(meal); err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	for _, dishID := range body.Dishes {
		dishIDParsed, _ := uuid.FromString(dishID)
		mealDish := domain.MealDish{
			MealID: meal.ID,
			DishID: dishIDParsed,
		}
		if err := mealDishRepo.Add(mealDish); err != nil {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}
	}

	c.JSON(http.StatusCreated, meal)
}

// Get returns array of meals
// @Summary Get list of categories with dishes for passed meal ID
// @Tags catering meals
// @Produce json
// @Param date query string false "Meal Date in 2020-01-01T00:00:00Z format"
// @Param id path string false "Catering ID"
// @Param clientId path string false "Client ID"
// @Success 200 {array} response.GetMeal "dishes for passed day"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/clients/{clientId}/meals [get]
func (m Meal) Get(c *gin.Context) {
	var query types.DateQuery
	var path types.PathClient

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	_, err := cateringRepo.GetByKey("id", path.ID)

	if err != nil {
		if err.Error() == "record not found" {
			utils.CreateError(http.StatusNotFound, err.Error(), c)
			return
		}
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	date, err := time.Parse(time.RFC3339, query.Date)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, "can't parse the date", c)
		return
	}

	result, code, err := mealRepo.Get(date, path.ID, path.ClientID)
	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, result)
}

// Update updates the menu for provided day, takes an array of dish Ids
// @Summary Returns 204 if success and 4xx if error
// @Produce json
// @Accept json
// @Tags catering meals
// @Param id path string true "Catering ID"
// @Param mealId path string false "Meal ID"
// @Param body body request.UpdateMeal false "Meal date"
// @Success 204 "Successfully updated"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/meals/{mealId} [put]
/*func (m Meal) Update(c *gin.Context) {
	var path types.PathMeal
	var body request.UpdateMeal

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	mealItem, code, err := mealRepo.GetByKey("id", path.MealID)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	var categoryIDs []uuid.UUID
	for _, dishID := range body.Dishes {
		result, _, _ := dishRepo.FindByID(path.ID, dishID)
		categoryIDs = append(categoryIDs, result.CategoryID)
	}
	duplicates := make(map[uuid.UUID]int)
	for _, id := range categoryIDs {
		_, exist := duplicates[id]
		if exist {
			duplicates[id]++
			if duplicates[id] > 10 {
				utils.CreateError(http.StatusBadRequest, "can't add more than 10 dishes in single category for current day", c)
				return
			}
		} else {
			duplicates[id] = 1
		}
	}

	if err := mealDishRepo.Delete(mealItem.ID.String()); err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	for _, dishID := range body.Dishes {
		dishIDParsed, _ := uuid.FromString(dishID)
		mealDish := domain.MealDish{
			MealID: mealItem.ID,
			DishID: dishIDParsed,
		}
		if err := mealDishRepo.Add(mealDish); err != nil {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}
	}
	c.Status(http.StatusNoContent)
}*/
