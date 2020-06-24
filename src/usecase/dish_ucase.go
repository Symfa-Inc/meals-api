package usecase

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go_api/src/domain"
	"go_api/src/repository"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
)

type dish struct{}

func NewDish() *dish {
	return &dish{}
}

var dishRepo = repository.NewDishRepo()

// AddDish adds dish for catering with provided ID
// @Summary Add dish for catering
// @Tags catering dishes
// @Produce json
// @Param id path string false "Catering ID"
// @Param payload body request.AddDish false "dish object"
// @Success 204 "Successfully created"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/dishes [post]
func (d dish) Add(c *gin.Context) {
	var path types.PathId
	var body domain.Dish

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	body.CateringID, _ = uuid.FromString(path.ID)

	_, err := dishRepo.Add(path.ID, body)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}
