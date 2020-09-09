package repository

import (
	"errors"
	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"net/http"
	"sort"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
)

// ClientScheduleRepo struct
type ClientScheduleRepo struct{}

// NewClientScheduleRepo returns pointer to
// client schedule repo with all methods
func NewClientScheduleRepo() *ClientScheduleRepo {
	return &ClientScheduleRepo{}
}

// Get returns array of schedules
// for provided client id
func (cs ClientScheduleRepo) Get(clientID string) ([]models.ClientSchedulesCatering, int, error) {
	var client domain.Client
	var cateringSchedules []domain.CateringSchedule
	var updatedSchedules []models.ClientSchedulesCatering
	if err := config.DB.
		Where("id = ?", clientID).
		Find(&client).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, http.StatusNotFound, errors.New("client with that id not found")
		}

		return nil, http.StatusBadRequest, err
	}

	if err := config.DB.
		Where("client_id = ?", clientID).
		Find(&domain.ClientSchedule{}).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, http.StatusNotFound, errors.New("client with that id not found")
		}

		return nil, http.StatusBadRequest, err
	}

	if err := config.DB.
		Where("catering_id = ?", client.CateringID).
		Find(&cateringSchedules).
		Order("day").
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, http.StatusNotFound, errors.New("catering with that id not found")
		}

		return nil, http.StatusBadRequest, err
	}

	for _, cateringSchedule := range cateringSchedules {
		var oldClientSchedule domain.ClientSchedule
		var updatedSchedule models.ClientSchedulesCatering
		config.DB.
			Where("client_id = ? AND day = ?", clientID, cateringSchedule.Day).
			Find(&oldClientSchedule)

		clientSchedule := domain.ClientSchedule{
			Day:       cateringSchedule.Day,
			Start:     cateringSchedule.Start,
			End:       cateringSchedule.End,
			IsWorking: oldClientSchedule.IsWorking,
			ClientID:  client.ID,
		}

		newStartTime, _ := now.Parse(cateringSchedule.Start)
		oldStartTime, _ := now.Parse(oldClientSchedule.Start)

		newEndTime, _ := now.Parse(cateringSchedule.End)
		oldEndTime, _ := now.Parse(oldClientSchedule.End)

		startTimeDifference := newStartTime.Sub(oldStartTime).Seconds()
		endTimeDifference := oldEndTime.Sub(newEndTime).Seconds()

		if endTimeDifference < 0 {
			clientSchedule.End = oldClientSchedule.End
		}

		if startTimeDifference < 0 {
			clientSchedule.Start = oldClientSchedule.Start
		}

		config.DB.
			Model(&domain.ClientSchedule{}).
			Where("client_id = ? AND day = ?", clientID, cateringSchedule.Day).
			Update(map[string]interface{}{
				"Day":       clientSchedule.Day,
				"Start":     clientSchedule.Start,
				"End":       clientSchedule.End,
				"IsWorking": clientSchedule.IsWorking,
			}).
			Scan(&updatedSchedule)

		updatedSchedule.CateringStart = cateringSchedule.Start
		updatedSchedule.CateringEnd = cateringSchedule.End
		updatedSchedules = append(updatedSchedules, updatedSchedule)
		sort.Slice(updatedSchedules, func(i, j int) bool {
			return updatedSchedules[i].Day < updatedSchedules[j].Day
		})
	}

	return updatedSchedules, 0, nil
}

// Update updates client's schedule with new values
func (cs ClientScheduleRepo) Update(clientID, scheduleID string, isWorking *bool, newSchedule domain.ClientSchedule) (models.ClientSchedulesCatering, int, error) {
	var client domain.Client
	var clientsOldSchedule domain.ClientSchedule
	var cateringSchedule domain.CateringSchedule
	var updatedSchedule models.ClientSchedulesCatering

	if err := config.DB.
		Where("id = ?", clientID).
		First(&client).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return models.ClientSchedulesCatering{}, http.StatusNotFound, errors.New("client with that id not found")
		}

		return models.ClientSchedulesCatering{}, http.StatusBadRequest, err
	}

	config.DB.Find(&clientsOldSchedule).Where("id", scheduleID)

	if err := config.DB.
		Model(&domain.ClientSchedule{}).
		Select("cs.*").
		Joins("left join clients c on client_schedules.client_id = c.id").
		Joins("left join catering_schedules cs on cs.catering_id = c.catering_id and cs.day = client_schedules.day").
		Where("client_schedules.id = ?", scheduleID).
		Scan(&cateringSchedule).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return models.ClientSchedulesCatering{}, http.StatusNotFound, errors.New("client schedule with that id not found")
		}
		return models.ClientSchedulesCatering{}, http.StatusBadRequest, err
	}

	if newSchedule.Start != "" && newSchedule.End != "" {
		startTime, err := now.Parse(newSchedule.Start)
		oldStartTime, _ := now.Parse(cateringSchedule.Start)

		if err != nil {
			return models.ClientSchedulesCatering{}, http.StatusBadRequest, err
		}

		startDifference := startTime.Sub(oldStartTime).Seconds()

		if startDifference < 0 {
			return models.ClientSchedulesCatering{}, http.StatusBadRequest, errors.New("new start time can't be earlier than catering's start time")
		}

		endTime, err := now.Parse(newSchedule.End)
		oldEndTime, _ := now.Parse(cateringSchedule.End)

		if err != nil {
			return models.ClientSchedulesCatering{}, http.StatusBadRequest, err
		}

		endDifference := oldEndTime.Sub(endTime).Seconds()

		if endDifference < 0 {
			return models.ClientSchedulesCatering{}, http.StatusBadRequest, errors.New("new end time can't be later than catering's end time")
		}

		difference := endTime.Sub(startTime).Seconds()

		if difference < 0 {
			return models.ClientSchedulesCatering{}, http.StatusBadRequest, errors.New("end date can't be earlier than start date")
		}
	}

	if newSchedule.Start != "" {
		startTime, err := now.Parse(newSchedule.Start)
		oldStartTime, _ := now.Parse(cateringSchedule.Start)

		if err != nil {
			return models.ClientSchedulesCatering{}, http.StatusBadRequest, err
		}

		startDifference := startTime.Sub(oldStartTime).Seconds()

		if startDifference < 0 {
			return models.ClientSchedulesCatering{}, http.StatusBadRequest, errors.New("new start time can't be earlier than catering's start time")
		}

		endTime, err := now.Parse(cateringSchedule.End)

		if err != nil {
			return models.ClientSchedulesCatering{}, http.StatusBadRequest, err
		}

		difference := endTime.Sub(startTime).Seconds()

		if difference < 0 {
			return models.ClientSchedulesCatering{}, http.StatusBadRequest, errors.New("end date can't be earlier than start date")
		}
	}

	if newSchedule.End != "" {
		endTime, err := now.Parse(newSchedule.End)
		oldEndTime, _ := now.Parse(cateringSchedule.End)

		if err != nil {
			return models.ClientSchedulesCatering{}, http.StatusBadRequest, err
		}

		endDifference := oldEndTime.Sub(endTime).Seconds()

		if endDifference < 0 {
			return models.ClientSchedulesCatering{}, http.StatusBadRequest, errors.New("new end time can't be later than catering's end time")
		}

		startTime, err := now.Parse(cateringSchedule.Start)

		if err != nil {
			return models.ClientSchedulesCatering{}, http.StatusBadRequest, err
		}

		difference := endTime.Sub(startTime).Seconds()

		if difference < 0 {
			return models.ClientSchedulesCatering{}, http.StatusBadRequest, errors.New("end date can't be earlier than start date")
		}
	}

	newSchedule.Day = cateringSchedule.Day
	updatedSchedule.CateringEnd = cateringSchedule.End
	updatedSchedule.CateringStart = cateringSchedule.Start

	if newSchedule.Start == "" {
		newSchedule.Start = cateringSchedule.Start
	}

	if newSchedule.End == "" {
		newSchedule.End = cateringSchedule.End
	}

	if isWorking == nil {
		newSchedule.IsWorking = clientsOldSchedule.IsWorking
	} else {
		newSchedule.IsWorking = *isWorking
	}

	config.DB.
		Model(&domain.ClientSchedule{}).
		Where("id = ?", scheduleID).
		Update(map[string]interface{}{
			"Day":       newSchedule.Day,
			"Start":     newSchedule.Start,
			"End":       newSchedule.End,
			"IsWorking": newSchedule.IsWorking,
		}).
		Scan(&updatedSchedule)

	return updatedSchedule, 0, nil
}
