package meal

import uuid "github.com/satori/go.uuid"

type AddMealRequest struct {
	ID uuid.UUID `json:"id" binding:"required"`
}
