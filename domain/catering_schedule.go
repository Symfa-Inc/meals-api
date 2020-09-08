package domain

import (
	"github.com/gin-gonic/gin"
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

// CateringScheduleRepository is repository interface
// for catering schedule
type CateringScheduleRepository interface {
	Get(cateringID string) ([]CateringSchedule, int, error)
	Update(cateringID, scheduleID string, isWorking *bool, newSchedule *CateringSchedule) (int, error)
}

// CateringScheduleAPI is API interface
// for catering schedule
type CateringScheduleAPI interface {
	Get(c *gin.Context)
	Update(c *gin.Context)
}


