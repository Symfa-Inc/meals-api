package domain

import (
	"github.com/Aiscom-LLC/meals-api/repository/models"
	uuid "github.com/satori/go.uuid"
)

// OrderDishes struct for DB
type OrderDishes struct {
	Base
	OrderID uuid.UUID
	DishID  uuid.UUID
	Amount  int
}

// OrderDishRepository is order interface for repository
type OrderDishRepository interface {
	CancelOrder(userID, orderID string) (int, error)
	GetUserOrder(userID, date string) (models.UserOrder, int, error)
	GetOrders(cateringID, clientID, date, companyType string) (models.SummaryOrderResult, int, error)
	ApproveOrders(clientID, date string) error
	GetOrdersStatus(clientID, date string) *string
}
