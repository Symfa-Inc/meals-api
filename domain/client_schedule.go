package domain

import (
	"github.com/gin-gonic/gin"
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

// ClientSchedulesCatering struct for joined catering table
type ClientSchedulesCatering struct {
	Base
	Day           int       `json:"day"`
	Start         string    `json:"start"`
	End           string    `json:"end"`
	IsWorking     bool      `json:"isWorking"`
	ClientID      uuid.UUID `json:"-" swaggerignore:"true"`
	CateringStart string    `json:"cateringStart"`
	CateringEnd   string    `json:"cateringEnd"`
} //@name GetClientSchedulesResponse

// ClientScheduleRepository is repository interface
// for client schedule
type ClientScheduleRepository interface {
	Get(clientID string) ([]ClientSchedulesCatering, int, error)
	Update(clientID, scheduleID string, isWorking *bool, newSchedule ClientSchedule) (ClientSchedulesCatering, int, error)
}

// ClientScheduleAPI is API interface
// for client schedule
type ClientScheduleAPI interface {
	Get(c *gin.Context)
	Update(c *gin.Context)
}
