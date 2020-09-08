package domain

import (
	uuid "github.com/satori/go.uuid"
)

// CateringSchedule struct
type CateringSchedule struct {
	Base
	Day        int       `json:"day"`
	Start      string    `json:"start"`
	End        string    `json:"end"`
	IsWorking  bool      `json:"isWorking"`
	CateringID uuid.UUID `json:"-" swaggerignore:"true"`
} //@name CateringSchedule
