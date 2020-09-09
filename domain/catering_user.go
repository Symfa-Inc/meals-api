package domain

import (
	uuid "github.com/satori/go.uuid"
)

// CateringUser struct
type CateringUser struct {
	Base
	CateringID uuid.UUID `json:"cateringId"`
	UserID     uuid.UUID `json:"userId"`
}
