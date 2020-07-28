package domain

import uuid "github.com/satori/go.uuid"

// OrderDishes struct for DB
type OrderDishes struct {
	Base
	OrderID uuid.UUID
	DishID  uuid.UUID
	Amount  int
}
