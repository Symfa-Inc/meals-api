package domain

import (
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"github.com/gin-gonic/gin"
)

// ClientScheduleRepository is repository interface
// for client schedule
type ClientScheduleRepository interface {
	Get(clientID string) ([]models.ClientSchedulesCatering, int, error)
	Update(clientID, scheduleID string, isWorking *bool, newSchedule domain.ClientSchedule) (models.ClientSchedulesCatering, int, error)
}

// ClientScheduleAPI is API interface
// for client schedule
type ClientScheduleAPI interface {
	Get(c *gin.Context)
	Update(c *gin.Context)
}
