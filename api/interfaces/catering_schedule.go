package domain

import (
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/gin-gonic/gin"
)

// CateringScheduleRepository is repository interface
// for catering schedule
type CateringScheduleRepository interface {
	Get(cateringID string) ([]domain.CateringSchedule, int, error)
	Update(cateringID, scheduleID string, isWorking *bool, newSchedule *domain.CateringSchedule) (int, error)
}

// CateringScheduleAPI is API interface
// for catering schedule
type CateringScheduleAPI interface {
	Get(c *gin.Context)
	Update(c *gin.Context)
}
