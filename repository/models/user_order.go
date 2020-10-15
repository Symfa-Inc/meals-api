package models

import uuid "github.com/satori/go.uuid"

// OrderItem struct for response
type OrderItem struct {
	ID     uuid.UUID `json:"id" gorm:"column:dish_id"`
	Image  *string   `json:"image" gorm:"column:path"`
	Price  int       `json:"price"`
	Name   string    `json:"name"`
	Amount int       `json:"amount"`
}

// UserOrder struct for response
type UserOrder struct {
	Items   []OrderItem `json:"items"`
	Status  string      `json:"status"`
	Total   float32     `json:"total"`
	OrderID uuid.UUID   `json:"orderId" gorm:"column:order_id"`
}
