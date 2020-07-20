package usecase

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go_api/src/domain"
	"go_api/src/repository"
	"go_api/src/schemes/request"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
	"time"
)

type meal struct{}

func NewMeal() *meal {
	return &meal{}
}

var mealRepo = repository.NewMealRepo()

// AddMeals add meals
// @Summary Add days for catering
// @Tags catering meals
// @Produce json
// @Param id path string false "Catering ID"
// @Param payload body request.AddMeal false "meal reading"
// @Success 201 {object} request.AddMeal "created meal"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/meals [post]
func (m meal) Add(c *gin.Context) {
	var path types.PathId
	var body request.AddMeal
	mealDishRepo := repository.NewMealDishesRepo()

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	parsedId, _ := uuid.FromString(path.ID)
	meal := domain.Meal{
		Date:       body.Date,
		CateringID: parsedId,
	}

	t := 24 * time.Hour
	difference := body.Date.Sub(time.Now().Truncate(t)).Hours()

	if difference < 0 {
		utils.CreateError(http.StatusBadRequest, "item has wrong date (can't use previous dates)", c)
		return
	}

	if err := mealRepo.Find(meal); err != nil {
		utils.CreateError(http.StatusBadRequest, "item already exist", c)
		return
	}

	for _, dishId := range body.Dishes {
		_, err, code := dishRepo.FindById(path.ID, dishId)
		if err != nil {
			utils.CreateError(code, err.Error(), c)
			return
		}
	}

	mealItem, err := mealRepo.Add(meal)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	for _, dishId := range body.Dishes {
		dishIdParsed, _ := uuid.FromString(dishId)
		mealDish := domain.MealDish{
			MealID: mealItem.(*domain.Meal).ID,
			DishID: dishIdParsed,
		}
		if err := mealDishRepo.Add(mealDish); err != nil {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}
	}

	c.JSON(http.StatusCreated, mealItem)
}

// GetMeals returns array of meals
// @Summary Get list of categories with dishes for passed meal ID
// @Tags catering meals
// @Produce json
// @Param date query string false "Meal Date in 2020-01-01T00:00:00Z format"
// @Param id path string false "Catering ID"
// @Success 200 {array} domain.Dish "dishes for passed day"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/meals [get]
func (m meal) Get(c *gin.Context) {
	var query types.DateQuery
	var pathUri types.PathId

	if err := utils.RequestBinderUri(&pathUri, c); err != nil {
		return
	}

	_, err := cateringRepo.GetByKey("id", pathUri.ID)

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

	result, mealId, err, code := mealRepo.Get(date, pathUri.ID)
	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mealId": mealId,
		"dishes": result,
	})
}

// UpdateMeal updates the menu for provided day, takes an array of dish Ids
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
func (m meal) Update(c *gin.Context) {
	var path types.PathMeal
	var body request.UpdateMeal
	mealDishRepo := repository.NewMealDishesRepo()

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	mealItem, err, code := mealRepo.GetByKey("id", path.MealId)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	if err := mealDishRepo.Delete(mealItem.ID.String()); err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	for _, dishId := range body.Dishes {
		dishIdParsed, _ := uuid.FromString(dishId)
		mealDish := domain.MealDish{
			MealID: mealItem.ID,
			DishID: dishIdParsed,
		}
		if err := mealDishRepo.Add(mealDish); err != nil {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}
	}
	c.Status(http.StatusNoContent)
}
