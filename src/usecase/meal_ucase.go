package usecase

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go_api/src/domain"
	"go_api/src/repository"
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
// @Success 201 {object} domain.Meal "meal reading"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/meals [post]
func (m meal) Add(c *gin.Context) {
	var path types.PathId
	var body domain.Meal

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	parsedId, _ := uuid.FromString(path.ID)
	t := 24 * time.Hour
	body.CateringID = parsedId
	difference := body.Date.Sub(time.Now().Truncate(t)).Hours()

	if difference < 0 {
		utils.CreateError(http.StatusBadRequest, "item has wrong date (can't use previous dates)", c)
		return
	}

	if err := mealRepo.Find(body); err != nil {
		utils.CreateError(http.StatusBadRequest, "item already exist", c)
		return
	}

	body.CateringID = parsedId
	mealItem, err := mealRepo.Add(body)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusCreated, mealItem)
}

// GetMeals returns array of meals
// @Summary Get list of categories with dishes for passed meal ID
// @Tags catering meals
// @Produce json
// @Param mealId query string false "Meal ID"
// @Param id path string false "Catering ID"
// @Success 200 {object} response.GetMealsResponse "categories with dishes for passed day"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/meals [get]
func (m meal) Get(c *gin.Context) {
	var mealQuery types.MealIdQuery
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

	if err := utils.RequestBinderQuery(&mealQuery, c); err != nil {
		return
	}

	result, err, code := mealRepo.Get(mealQuery.MealId, pathUri.ID)
	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateMeal godoc
// @Summary Returns updated meal
// @Produce json
// @Accept json
// @Tags catering meals
// @Param id path string true "Catering ID"
// @Param mealId path string true "Meal ID"
// @Param body body request.AddMeal false "Meal date"
// @Success 204 "Successfully updated"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/meals/{mealId} [put]
func (m meal) Update(c *gin.Context) {
	var path types.PathMeal
	var body domain.Meal

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	err, code := mealRepo.Update(path, body)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}
