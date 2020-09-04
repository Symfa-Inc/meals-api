package services

import (
	"errors"
	"github.com/Aiscom-LLC/meals-api/api/swagger"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
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

func (o *OrderService) Add(query string, order swagger.OrderRequest, claims jwt.MapClaims) (swagger.UserOrder, int, error) {
	userRepo := repository.NewUserRepo()
	var userID string

	id := claims["id"].(string)

	user, _ := userRepo.GetByKey("id", id)

	if user.ID == uuid.Nil {
		userID = ""
	} else {
		userID = user.ID.String()
	}
	for i, dish := range order.Items {
		if dish.Amount == 0 {
			return swagger.UserOrder{}, http.StatusBadRequest, errors.New("can't add dish with 0 amount")
		}
		for j := i + 1; j < len(order.Items); j++ {
			if dish.DishID == order.Items[j].DishID {
				return swagger.UserOrder{}, http.StatusBadRequest, errors.New("can't add 2 same dishes, please increment amount field instead")
			}
		}
	}

	date, err := time.Parse(time.RFC3339, query)

	if err != nil {
		return swagger.UserOrder{}, http.StatusBadRequest, err
	}

	difference := date.Sub(time.Now().Truncate(time.Hour * 24)).Hours()

	if difference < 0 {
		return swagger.UserOrder{}, http.StatusBadRequest, errors.New("can't add order to previous date")
	}

	userOrder, err := orderRepo.Add(userID, date, order)

	if err != nil {
		return swagger.UserOrder{}, http.StatusBadRequest, err
	}

	return userOrder, 0, nil
}
