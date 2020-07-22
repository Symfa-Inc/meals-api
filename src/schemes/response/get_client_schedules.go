package response

import (
	uuid "github.com/satori/go.uuid"
	"go_api/src/domain"
)

// GetClientSchedules response scheme
type GetClientSchedules struct {
	domain.Base
	Day           int       `json:"day"`
	Start         string    `json:"start"`
	End           string    `json:"end"`
	IsWorking     bool      `json:"isWorking"`
	ClientID      uuid.UUID `json:"-" swaggerignore:"true"`
	CateringStart string    `json:"cateringStart"`
	CateringEnd   string    `json:"cateringEnd"`
} //@name GetClientSchedulesResponse
