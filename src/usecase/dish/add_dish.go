package dish

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go_api/src/models"
	"go_api/src/repository"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
)

// AddDish adds dish for catering with provided ID
// @Summary Add dish for catering
// @Tags catering dish
// @Produce json
// @Param id path string false "Catering ID"
// @Param payload body request.AddDish false "dish object"
// @Success 201 {object} models.Dish "dish object with ID"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/dish [post]
func AddDish(c *gin.Context) {
	var path types.PathId
	var body models.Dish

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	body.CateringID, _ = uuid.FromString(path.ID)

	result, err := repository.CreateDish(path.ID, body)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	fmt.Println(result)
	fmt.Println(err)
	c.JSON(http.StatusCreated, result)
}
