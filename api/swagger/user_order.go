package swagger

import (
	uuid "github.com/satori/go.uuid"
)

// UserOrder struct for response
type UserOrder struct {
	Items   []OrderItem `json:"items"`
	Status  string      `json:"status"`
	Total   int         `json:"total"`
	OrderID uuid.UUID   `json:"orderId" gorm:"type:column:order_id"`
}

// OrderItem struct for response
type OrderItem struct {
	ID     uuid.UUID `json:"id" gorm:"column:dish_id"`
	Image  *string   `json:"image" gorm:"column:path"`
	Price  int       `json:"price"`
	Name   string    `json:"name"`
	Amount int       `json:"amount"`
}
