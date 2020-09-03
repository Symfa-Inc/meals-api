package response

import (
	uuid "github.com/satori/go.uuid"
)

// UserOrder struct for response
type UserOrder struct {
	Items   []OrderItem `json:"items"`
	Status  string      `json:"status"`
	Total   float32         `json:"total"`
	OrderID uuid.UUID   `json:"orderId" gorm:"column:order_id"`
}

// OrderItem struct for response
type OrderItem struct {
	ID     uuid.UUID `json:"id" gorm:"column:dish_id"`
	Image  *string   `json:"image" gorm:"column:path"`
	Price  float32       `json:"price"`
	Name   string    `json:"name"`
	Amount int       `json:"amount"`
}
