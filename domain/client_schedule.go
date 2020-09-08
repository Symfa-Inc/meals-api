package domain

import (
	uuid "github.com/satori/go.uuid"
)

// ClientSchedule struct
type ClientSchedule struct {
	Base
	Day       int       `json:"day"`
	Start     string    `json:"start"`
	End       string    `json:"end"`
	IsWorking bool      `json:"isWorking"`
	ClientID  uuid.UUID `json:"-" swaggerignore:"true"`
} //@name ClientSchedule
