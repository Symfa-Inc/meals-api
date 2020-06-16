package meals

import (
	"github.com/gin-gonic/gin"
	"go_api/src/repository/meal"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
)

// AddDays godoc
// @Summary Add days for catering
// @Tags meals
// @Produce json
// @Param startDate query string true "example: 2006-01-02T00:00:00"
// @Param endDate query string true "example: 2006-01-09T00:00:00"
// @Param body body types.PathId false "Catering ID"
// @Success 201 {array} models.Meal "array of meal readings"
// @Failure 400 {object} types.Error "Error"
// @Router /meals [post]
func AddMeals(c *gin.Context) {
	var query types.StartEndDateQuery
	var body types.PathId

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	mealsArray, err := meal.CreateMealsDB(query, body.ID)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusCreated, mealsArray)
}
