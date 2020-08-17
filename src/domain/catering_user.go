package domain

import (
	uuid "github.com/satori/go.uuid"
)

type CateringUser struct {
	Base
	CateringID uuid.UUID `json:"cateringId"`
	UserID     uuid.UUID `json:"userId"`
}

type CateringUserRepository interface {
	GetByKey(key, value string) (User, error)
}
