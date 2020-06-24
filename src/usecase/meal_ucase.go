package usecase

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go_api/src/domain"
	"go_api/src/repository"
	"go_api/src/schemes/response"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
	"strconv"
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
// @Param payload body request.AddMealList false "array of meals"
// @Success 201 {array} domain.Meal "array of meal readings"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/meals [post]
func (m meal) Add(c *gin.Context) {
	var path types.PathId
	var body []domain.Meal
	var resultMealArray []*domain.Meal

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	parsedId, _ := uuid.FromString(path.ID)
	for i, _ := range body {
		t := 24 * time.Hour
		body[i].CateringID = parsedId
		difference := body[i].Date.Sub(time.Now().Truncate(t)).Hours()

		if difference < 0 {
			utils.CreateError(http.StatusBadRequest, "item "+strconv.Itoa(i+1)+" has wrong date (can't use previous dates)", c)
			return
		}

		if err := mealRepo.Find(body[i]); err != nil {
			utils.CreateError(http.StatusBadRequest, "item "+strconv.Itoa(i+1)+" already exist", c)
			return
		}
	}

	for i := range body {
		body[i].CateringID = parsedId
		mealItem, err := mealRepo.Add(body[i])
		if err != nil {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}
		resultMealArray = append(resultMealArray, mealItem.(*domain.Meal))
	}

	c.JSON(http.StatusCreated, resultMealArray)
}

// GetMeals returns array of meals
// @Summary Get list of meals withing provided date range
// @Tags catering meals
// @Produce json
// @Param startDate query string false "example: 2006-01-02"
// @Param endDate query string false "example: 2006-01-09"
// @Param limit query int false "limit of returned array"
// @Param id path string false "Catering ID"
// @Success 200 {object} response.GetMealsModel "array of meal readings"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/meals [get]
func (m meal) Get(c *gin.Context) {
	var limitQuery types.PaginationQuery
	var dateQuery types.StartEndDateQuery
	var pathUri types.PathId

	if err := utils.RequestBinderUri(&pathUri, c); err != nil {
		return
	}

	_, err := cateringRepo.GetByKey("id", pathUri.ID)

	if err != nil {
		if err.Error() == "record not found" {
			utils.CreateError(http.StatusNotFound, err.Error(), c)
			return
		} else {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}
	}

	if err := utils.RequestBinderQuery(&dateQuery, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&limitQuery, c); err != nil {
		return
	}

	result, total, err := mealRepo.Get(limitQuery.Limit, dateQuery, pathUri.ID)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, response.GetMealsModel{
		Items: result,
		Total: total,
	})
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

	err := mealRepo.Update(path, body)

	if err != nil {
		if err.Error() == "meal not found" {
			utils.CreateError(http.StatusNotFound, err.Error(), c)
			return
		}
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}
