package repository

import (
	"errors"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/types"
	"net/http"
)

type clientRepo struct{}

func NewClientRepo() *clientRepo {
	return &clientRepo{}
}

// Add adds client in DB
// returns error if that client name already exists
func (c clientRepo) Add(client domain.Client) (domain.Client, error) {
	if exist := config.DB.Where("name = ?", client.Name).
		Find(&client).RowsAffected; exist != 0 {
		return domain.Client{}, errors.New("client with that name already exist")
	}
	err := config.DB.Create(&client).Error
	return client, err
}

// Returns list of clients
func (c clientRepo) Get(query types.PaginationQuery) ([]domain.Client, int, error) {
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

	config.DB.Find(&clients).Count(&total)

	err := config.DB.
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&clients).
		Error

	return clients, total, err
}

// Soft delete of client
func (c clientRepo) Delete(id string) error {
	if result := config.DB.Where("id = ?", id).
		Delete(&domain.Client{}).RowsAffected; result == 0 {
		return errors.New("client not found")
	}

	return nil
}

// Update updates client with passed args
// returns error and status code
func (c clientRepo) Update(id string, client domain.Client) (error, int) {
	if clientExist := config.DB.Where("id = ?", id).
		Find(&domain.Client{}).RowsAffected; clientExist == 0 {
		return errors.New("client not found"), http.StatusNotFound
	}

	if nameExist := config.DB.Where("name = ?", client.Name).
		Find(&client).RowsAffected; nameExist != 0 {
		return errors.New("client with that name already exist"), http.StatusBadRequest
	}
	return config.DB.Model(&client).Where("id = ?", id).Update(&client).Error, 0
}

// Get client by provided key value arguments
// Returns client, error
func (c clientRepo) GetByKey(key, value string) (domain.Client, error) {
	var client domain.Client
	err := config.DB.Where(key+" = ?", value).First(&client).Error
	return client, err
}
