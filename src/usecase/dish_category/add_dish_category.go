package dish_category

import "github.com/gin-gonic/gin"

// AddDishCategory add dish category in DB
// returns 201 if success and 4xx if request failed
// @Summary Returns created category
// @Produce json
// @Accept json
// @Tags catering dish-categories
// @Param body body models.DishCategory false "Category Name"
// @Success 201 {object} models.Catering
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/dish-category [post]
func AddDishCategory(c *gin.Context) {

}
