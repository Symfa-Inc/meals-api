package swagger

import uuid "github.com/satori/go.uuid"

// Order struct for request scheme
type Order struct {
	DishID uuid.UUID `json:"dishId"`
	Amount int       `json:"amount"`
}

// OrderRequest struct for request scheme
type OrderRequest struct {
	Items   []Order `json:"items" binding:"required"`
	Comment string  `json:"comment"`
}
