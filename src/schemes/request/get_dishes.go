package request

import uuid "github.com/satori/go.uuid"

type GetDishes struct {
	CategoryID uuid.UUID `json:"categoryId" binding:"required"`
}
