package services

import (
	"errors"
	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/schemes/request"
	"github.com/Aiscom-LLC/meals-api/types"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// OrderService struct
type OrderService struct{}

// NewOrderService return pointer to order struct
// with all methods
func NewOrderService() *OrderService {
	return &OrderService{}
}

var orderRepo = repository.NewOrderRepo()

func (o *OrderService) Add(c *gin.Context, order request.OrderRequest, path types.PathID) {
	var query = c.Query("date")


	_, err := middleware.Passport().GetClaimsFromJWT(c)
	if err != nil {
		return
	}

	for i, dish := range order.Items {
		if dish.Amount == 0 {
			utils.CreateError(http.StatusBadRequest, errors.New("can't add dish with 0 amount"), c)
			return
		}
		for j := i + 1; j < len(order.Items); j++ {
			if dish.DishID == order.Items[j].DishID {
				utils.CreateError(http.StatusBadRequest, errors.New("can't add 2 same dishes, please increment amount field instead"), c)
				return
			}
		}
	}

	date, err := time.Parse(time.RFC3339, query)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	difference := date.Sub(time.Now().Truncate(time.Hour * 24)).Hours()

	if difference < 0 {
		utils.CreateError(http.StatusBadRequest, errors.New("can't add order to previous date"), c)
		return
	}

	userOrder, err := orderRepo.Add(path.ID, date, order)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	c.JSON(http.StatusCreated, userOrder)
}
