package api

import (
	"github.com/Aiscom-LLC/meals-api/api/api_types"
	"github.com/Aiscom-LLC/meals-api/api/swagger"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ClientSchedule struct
type ClientSchedule struct{}

// NewClientSchedule returns pointer to client schedule
// with all methods
func NewClientSchedule() *ClientSchedule {
	return &ClientSchedule{}
}

var clientScheduleRepo = repository.NewClientScheduleRepo()

// Get return list of schedules
// @Summary Returns list of schedules
// @Tags clients schedules
// @Produce json
// @Param id path string false "Client ID"
// @Success 200 {array} domain.ClientSchedulesCatering "List of schedules"
// @Failure 400 {object} api_types.Error "Error"
// @Failure 404 {object} api_types.Error "Error"
// @Router /clients/{id}/schedules [get]
func (cs ClientSchedule) Get(c *gin.Context) {
	var path api_types.PathID

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	schedules, code, err := clientScheduleRepo.Get(path.ID)
	client, _ := clientRepo.GetByKey("id", path.ID)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list":              schedules,
		"autoApproveOrders": client.AutoApproveOrders,
	})
}

// Update updates schedule
// @Summary Returns 200 and updated model if success and 4xx error if failed
// @Produce json
// @Accept json
// @Tags clients schedules
// @Param id path string true "Client ID"
// @Param scheduleId path string true "Client Schedule ID"
// @Param body body request.UpdateSchedule false "Client Schedule model"
// @Success 200 {object} domain.ClientSchedulesCatering "Client Schedule model"
// @Failure 400 {object} api_types.Error "Error"
// @Failure 404 {object} api_types.Error "Not Found"
// @Router /clients/{id}/schedules/{scheduleId} [put]
func (cs ClientSchedule) Update(c *gin.Context) {
	var path api_types.PathSchedule
	var body swagger.UpdateSchedule

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	schedule := domain.ClientSchedule{
		Start: body.Start,
		End:   body.End,
	}

	updatedSchedule, code, err := clientScheduleRepo.Update(path.ID, path.ScheduleID, body.IsWorking, schedule)
	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	clientRepo.UpdateAutoApproveSchedules(path.ID)
	c.JSON(http.StatusOK, updatedSchedule)
}
