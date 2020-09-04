package models

import (
	"github.com/Aiscom-LLC/meals-api/domain"
	uuid "github.com/satori/go.uuid"
)

type GetCateringUser struct {
	ID                  uuid.UUID `json:"id"`
	domain.UserCatering `json:"catering"`
	FirstName           string  `json:"firstName"`
	LastName            string  `json:"lastName"`
	Email               string  `json:"email"`
	Role                string  `json:"role"`
	Status              *string `json:"status"`
}