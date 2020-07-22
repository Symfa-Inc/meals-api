package usecase

import (
	"github.com/gin-gonic/gin"
	"go_api/src/domain"
	"go_api/src/repository"
	"go_api/src/schemes/request"
	"go_api/src/types"
	"go_api/src/utils"
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
// @Success 200 {array} response.GetClientSchedules "List of schedules"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /clients/{id}/schedules [get]
func (cs ClientSchedule) Get(c *gin.Context) {
	var path types.PathID

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	schedules, code, err := clientScheduleRepo.Get(path.ID)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, schedules)
}

// Update updates schedule
// @Summary Returns 200 and updated model if success and 4xx error if failed
// @Produce json
// @Accept json
// @Tags clients schedules
// @Param id path string true "Client ID"
// @Param scheduleId path string true "Client Schedule ID"
// @Param body body request.UpdateSchedule false "Client Schedule model"
// @Success 200 {object} response.GetClientSchedules "Client Schedule model"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /clients/{id}/schedules/{scheduleId} [put]
func (cs ClientSchedule) Update(c *gin.Context) {
	var path types.PathSchedule
	var body request.UpdateSchedule

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
		utils.CreateError(code, err.Error(), c)
		return
	}
	c.JSON(http.StatusOK, updatedSchedule)
}
