package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/types"
	"net/http"
)

// ClientRepo struct
type ClientRepo struct{}

// NewClientRepo returns pointer to client repository
// with all methods
func NewClientRepo() *ClientRepo {
	return &ClientRepo{}
}

// Add adds client in DB
// returns error if that client name already exists
func (c ClientRepo) Add(cateringID string, client *domain.Client) error {
	if err := config.DB.
		Where("id = ?", cateringID).
		Find(&domain.Catering{}).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return err
		}

		return err
	}

	if exist := config.DB.
		Where("name = ?", client.Name).
		Find(client).RowsAffected; exist != 0 {
		return errors.New("client with that name already exist")
	}

	err := config.DB.Create(client).Error

	var cateringSchedules []domain.CateringSchedule

	config.DB.
		Where("catering_id = ?", client.CateringID).
		Find(&cateringSchedules)

	for _, schedule := range cateringSchedules {
		clientSchedule := domain.ClientSchedule{
			Day:       schedule.Day,
			Start:     schedule.Start,
			End:       schedule.End,
			IsWorking: schedule.IsWorking,
			ClientID:  client.ID,
		}
		config.DB.Create(&clientSchedule)
	}

	return err
}

// Get returns list of clients
func (c ClientRepo) Get(cateringID string, query types.PaginationQuery) ([]domain.Client, int, error) {
	var clients []domain.Client
	var total int

	page := query.Page
	limit := query.Limit

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	config.DB.
		Where("catering_id = ?", cateringID).
		Find(&clients).
		Count(&total)

	err := config.DB.
		Limit(limit).
		Offset((page-1)*limit).
		Where("catering_id = ?", cateringID).
		Find(&clients).
		Error

	return clients, total, err
}

// Delete soft delete of client
func (c ClientRepo) Delete(cateringID, id string) error {
	if result := config.DB.Where("id = ? AND catering_id = ?", id, cateringID).
		Delete(&domain.Client{}).RowsAffected; result == 0 {
		return errors.New("client not found")
	}

	return nil
}

// Update updates client with passed args
// returns error and status code
func (c ClientRepo) Update(cateringID, id string, client domain.Client) (int, error) {
	if nameExist := config.DB.Where("name = ?", client.Name).
		Find(&client).RowsAffected; nameExist != 0 {
		return http.StatusBadRequest, errors.New("client with that name already exist")
	}

	if clientExist := config.DB.Model(&client).Where("id = ? AND catering_id = ?", id, cateringID).
		Update(&client).RowsAffected; clientExist == 0 {
		return http.StatusNotFound, errors.New("client not found")
	}

	return 0, nil
}

// GetByKey client by provided key value arguments
// Returns client, error
func (c ClientRepo) GetByKey(key, value string) (domain.Client, error) {
	var client domain.Client
	err := config.DB.Where(key+" = ?", value).First(&client).Error
	return client, err
}
