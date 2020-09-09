package models

import (
	uuid "github.com/satori/go.uuid"
)

type GetCateringUser struct {
	ID           uuid.UUID `json:"id"`
	UserCatering `json:"catering"`
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	Email        string  `json:"email"`
	Role         string  `json:"role"`
	Status       *string `json:"status"`
}
