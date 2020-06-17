package meals

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go_api/src/models"
	"go_api/src/repository"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
	"strconv"
	"time"
)

// AddMeals godoc
// @Summary Add days for catering
// @Tags catering meals
// @Produce json
// @Param id path string false "Catering ID"
// @Param payload body request.AddMealRequestList false "array of meals"
// @Success 201 {array} models.Meal "array of meal readings"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/meals [post]
func AddMeals(c *gin.Context) {
	var path types.PathId
	var body []models.Meal
	var resultMealArray []*models.Meal

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

		if err := repository.FindMealDB(body[i]); err != nil {
			utils.CreateError(http.StatusBadRequest, "item "+strconv.Itoa(i+1)+" already exist", c)
			return
		}
	}

	for i := range body {
		body[i].CateringID = parsedId
		mealItem, err := repository.CreateMealDB(body[i])
		if err != nil {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}
		resultMealArray = append(resultMealArray, mealItem.Value.(*models.Meal))
	}

	c.JSON(http.StatusCreated, resultMealArray)
}
