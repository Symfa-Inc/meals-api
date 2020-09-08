package api

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/interfaces"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CateringSchedule struct
type CateringSchedule struct{}

// NewCateringSchedule returns pointer to catering schedule
// with all methods
func NewCateringSchedule() *CateringSchedule {
	return &CateringSchedule{}
}

var cateringScheduleRepo = repository.NewCateringScheduleRepo()

// Get return list of schedules
// @Summary Returns list of schedules
// @Tags caterings schedules
// @Produce json
// @Param id path string false "Catering ID"
// @Success 200 {array} domain.CateringSchedule "List of schedules"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Error"
// @Router /caterings/{id}/schedules [get]
func (s CateringSchedule) Get(c *gin.Context) {
	var path url.PathID

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	schedules, code, err := cateringScheduleRepo.Get(path.ID)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list": schedules,
	})
}

// Update updates schedule
// @Summary Returns 200 and updated model if success and 4xx error if failed
// @Produce json
// @Accept json
// @Tags caterings schedules
// @Param id path string true "Catering ID"
// @Param scheduleId path string true "CateringSchedule ID"
// @Param body body swagger.UpdateSchedule false "CateringSchedule model"
// @Success 200 {object} domain.CateringSchedule "CateringSchedule model"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /caterings/{id}/schedules/{scheduleId} [put]
func (s CateringSchedule) Update(c *gin.Context) {
	var path url.PathSchedule
	var body models.UpdateSchedule

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	schedule := &interfaces.CateringSchedule{
		Start: body.Start,
		End:   body.End,
	}

	code, err := cateringScheduleRepo.Update(path.ID, path.ScheduleID, body.IsWorking, schedule)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, schedule)
}
