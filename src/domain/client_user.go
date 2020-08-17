package domain

import uuid "github.com/satori/go.uuid"

type ClientUser struct {
	Base
	ClientID uuid.UUID `json:"clientId"`
	UserID   uuid.UUID `json:"userId"`
	Enum     int
	Floor    int `json:"floor"`
}
