package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"go_api/src/config"
	"go_api/src/domain"
	"net/http"
)

// CateringScheduleRepo struct
type CateringScheduleRepo struct{}

// NewCateringScheduleRepo returns pointer to
// catering schedule repo with all methods
func NewCateringScheduleRepo() *CateringScheduleRepo {
	return &CateringScheduleRepo{}
}

// Get returns array of schedules
// for provided catering id
func (cs CateringScheduleRepo) Get(cateringID string) ([]domain.CateringSchedule, int, error) {
	var schedules []domain.CateringSchedule
	if err := config.DB.
		Where("id = ?", cateringID).
		Find(&domain.Catering{}).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, http.StatusNotFound, errors.New("catering with that id not found")
		}

		return nil, http.StatusBadRequest, err
	}

	if err := config.DB.
		Model(&domain.CateringSchedule{}).
		Where("catering_schedules.catering_id = ?", cateringID).
		Order("day").
		Scan(&schedules).Error; err != nil {
		return nil, http.StatusBadRequest, err
	}

	return schedules, 0, nil
}

// Update updates schedule and returns new updated schedule
func (cs CateringScheduleRepo) Update(cateringID, scheduleID string, isWorking *bool, newSchedule domain.CateringSchedule) (domain.CateringSchedule, int, error) {
	var schedule domain.CateringSchedule

	if err := config.DB.
		Where("id = ?", cateringID).
		Find(&domain.Catering{}).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return domain.CateringSchedule{}, http.StatusNotFound, errors.New("catering with that id not found")
		}

		return domain.CateringSchedule{}, http.StatusBadRequest, err
	}

	if err := config.DB.
		Where("id = ?", scheduleID).
		Find(&schedule).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return domain.CateringSchedule{}, http.StatusNotFound, errors.New("schedule with that id not found")
		}

		return domain.CateringSchedule{}, http.StatusBadRequest, err
	}

	if newSchedule.Start != "" && newSchedule.End != "" {
		startTime, err := now.Parse(newSchedule.Start)

		if err != nil {
			return domain.CateringSchedule{}, http.StatusBadRequest, err
		}

		endTime, err := now.Parse(newSchedule.End)

		if err != nil {
			return domain.CateringSchedule{}, http.StatusBadRequest, err
		}

		difference := endTime.Sub(startTime).Seconds()

		if difference < 0 {
			return domain.CateringSchedule{}, http.StatusBadRequest, errors.New("end date can't be earlier than start date")
		}
	}

	if newSchedule.Start != "" {
		startTime, err := now.Parse(newSchedule.Start)

		if err != nil {
			return domain.CateringSchedule{}, http.StatusBadRequest, err
		}

		endTime, err := now.Parse(schedule.End)

		if err != nil {
			return domain.CateringSchedule{}, http.StatusBadRequest, err
		}

		difference := endTime.Sub(startTime).Seconds()

		if difference < 0 {
			return domain.CateringSchedule{}, http.StatusBadRequest, errors.New("end date can't be earlier than start date")
		}
	}

	if newSchedule.End != "" {
		endTime, err := now.Parse(newSchedule.End)

		if err != nil {
			return domain.CateringSchedule{}, http.StatusBadRequest, err
		}

		startTime, err := now.Parse(schedule.Start)

		if err != nil {
			return domain.CateringSchedule{}, http.StatusBadRequest, err
		}

		difference := endTime.Sub(startTime).Seconds()

		if difference < 0 {
			return domain.CateringSchedule{}, http.StatusBadRequest, errors.New("end date can't be earlier than start date")
		}
	}

	newSchedule.Day = schedule.Day
	newSchedule.ID = schedule.ID

	if newSchedule.Start == "" {
		newSchedule.Start = schedule.Start
	}

	if newSchedule.End == "" {
		newSchedule.End = schedule.End
	}

	if isWorking == nil {
		newSchedule.IsWorking = schedule.IsWorking
	} else {
		newSchedule.IsWorking = *isWorking
	}

	config.DB.
		Model(&schedule).
		Where("id = ?", schedule.ID).
		Update(map[string]interface{}{
			"Day":       newSchedule.Day,
			"Start":     newSchedule.Start,
			"End":       newSchedule.End,
			"IsWorking": newSchedule.IsWorking,
		}).
		Scan(&newSchedule)

	return newSchedule, 0, nil
}
