package models

import (
	"github.com/Aiscom-LLC/meals-api/domain"
	uuid "github.com/satori/go.uuid"
)

// ClientSchedulesCatering struct for joined catering table
type ClientSchedulesCatering struct {
	domain.Base
	Day           int       `json:"day"`
	Start         string    `json:"start"`
	End           string    `json:"end"`
	IsWorking     bool      `json:"isWorking"`
	ClientID      uuid.UUID `json:"-" swaggerignore:"true"`
	CateringStart string    `json:"cateringStart"`
	CateringEnd   string    `json:"cateringEnd"`
} //@name GetClientSchedulesResponse
