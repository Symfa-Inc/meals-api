package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
	"time"
)

// UserRepo struct
type UserRepo struct{}

// NewUserRepo returns pointer to user repository
// with all methods
func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

// GetByKey returns user by key
// and error if exist
func (ur UserRepo) GetByKey(key, value string) (domain.UserClientCatering, error) {
	var user domain.UserClientCatering
	err := config.DB.
		Model(&domain.User{}).
		Select("users.*, c.id as catering_id, c.name as catering_name, ci.id as client_id, ci.name as client_name").
		Joins("left join caterings c on c.id = users.catering_id").
		Joins("left join clients ci on ci.id = users.client_id").
		Where("users."+key+" = ?", value).
		Scan(&user).Error

	return user, err
}

// Get returns list of users for provided company
func (ur UserRepo) Get(companyID, companyType string, pagination types.PaginationQuery, filters types.UserFilterQuery) ([]domain.UserClientCatering, int, int, error) {
	var users []domain.UserClientCatering
	var total int
	page := pagination.Page
	limit := pagination.Limit
	status := filters.Status
	role := filters.Role
	clientName := filters.ClientName
	querySearch := filters.Query

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	if companyType == types.CompanyTypesEnum.Catering {
		if role == "" {
			role = types.UserRoleEnum.CateringAdmin
		}

		if err := config.DB.
			Unscoped().
			Model(&domain.User{}).
			Select("users.*, c.id as catering_id, c.name as catering_name, ci.id as client_id, ci.name as client_name").
			Joins("left join caterings c on c.id = users.catering_id").
			Joins("left join clients ci on ci.id = users.client_id").
			Where("users.catering_id = ? AND (first_name || last_name) ILIKE ?"+
				" AND ci.name ILIKE ? AND CAST(users.role AS text) ILIKE ?"+
				" AND (users.deleted_at > ? OR users.deleted_at IS NULL)", companyID, "%"+querySearch+"%", "%"+clientName+"%", "%"+role+"%", time.Now()).
			Count(&total).
			Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return nil, 0, http.StatusNotFound, errors.New("user not found")
			}
			return nil, 0, http.StatusBadRequest, err
		}

		if err := config.DB.
			Unscoped().
			Limit(limit).
			Offset((page-1)*limit).
			Model(&domain.User{}).
			Select("users.*, c.id as catering_id, c.name as catering_name, ci.id as client_id, ci.name as client_name").
			Joins("left join caterings c on c.id = users.catering_id").
			Joins("left join clients ci on ci.id = users.client_id").
			Where("users.catering_id = ? AND (first_name || last_name) ILIKE ?"+
				" AND ci.name ILIKE ? AND CAST(users.role AS text) ILIKE ?"+
				" AND (users.deleted_at > ? OR users.deleted_at IS NULL)", companyID, "%"+querySearch+"%", "%"+clientName+"%", "%"+role+"%", time.Now()).
			Order("created_at DESC, first_name ASC").
			Scan(&users).
			Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return nil, 0, http.StatusNotFound, errors.New("user not found")
			}
			return nil, 0, http.StatusBadRequest, err
		}
		return users, total, 0, nil
	}

	if err := config.DB.
		Unscoped().
		Model(&domain.User{}).
		Select("users.*, c.id as catering_id, c.name as catering_name, ci.id as client_id, ci.name as client_name").
		Joins("left join caterings c on c.id = users.catering_id").
		Joins("left join clients ci on ci.id = users.client_id").
		Where("users.client_id = ? AND (first_name || last_name) ILIKE ?"+
			"AND status ILIKE ? AND CAST(users.role AS text) ILIKE ?"+
			" AND (users.deleted_at > ? OR users.deleted_at IS NULL)", companyID, "%"+querySearch+"%", "%"+status+"%", "%"+role+"%", time.Now()).
		Count(&total).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, 0, http.StatusNotFound, errors.New("user not found")
		}
		return nil, 0, http.StatusBadRequest, err
	}

	if err := config.DB.
		Unscoped().
		Limit(limit).
		Offset((page-1)*limit).
		Model(&domain.User{}).
		Select("users.*, c.id as catering_id, c.name as catering_name, ci.id as client_id, ci.name as client_name").
		Joins("left join caterings c on c.id = users.catering_id").
		Joins("left join clients ci on ci.id = users.client_id").
		Where("users.client_id = ? AND (first_name || last_name) ILIKE ?"+
			"AND status ILIKE ? AND CAST(users.role AS text) ILIKE ?"+
			" AND (users.deleted_at > ? OR users.deleted_at IS NULL)", companyID, "%"+querySearch+"%", "%"+status+"%", "%"+role+"%", time.Now()).
		Order("created_at DESC, first_name ASC").
		Scan(&users).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, 0, http.StatusNotFound, errors.New("user not found")
		}
		return nil, 0, http.StatusBadRequest, err
	}
	return users, total, 0, nil
}

// Add adds new user for certain company passed in user struct
// returns user and error
func (ur UserRepo) Add(user domain.User) (domain.UserClientCatering, error) {
	var createdUser domain.UserClientCatering
	if err := config.DB.
		Model(&domain.User{}).
		Create(&user).
		Select("users.*, c.id as catering_id, c.name as catering_name, ci.id as client_id, ci.name as client_name").
		Joins("left join caterings c on c.id = users.catering_id").
		Joins("left join clients ci on ci.id = users.client_id").
		Where("users.id = ?", user.ID).
		Scan(&createdUser).
		Error; err != nil {
		return domain.UserClientCatering{}, err
	}
	return createdUser, nil
}

// Delete changes status of user to deleted
// and sets deleted_at field to 21 days from now
func (ur UserRepo) Delete(companyID string, user domain.User) (int, error) {
	companyType := utils.DerefString(user.CompanyType)
	if companyType == types.CompanyTypesEnum.Catering {
		if userExist := config.DB.
			Model(&domain.User{}).
			Where("catering_id = ?", companyID).
			Update(&user).
			RowsAffected; userExist == 0 {
			return http.StatusNotFound, errors.New("user not found")
		}
		return 0, nil
	}

	if userExist := config.DB.
		Model(&domain.User{}).
		Where("client_id = ?", companyID).
		Update(&user).
		RowsAffected; userExist == 0 {
		return http.StatusBadRequest, errors.New("user not found")
	}
	return 0, nil
}

// Update updates user for passed company ID
// checks if user belongs to client or catering
func (ur UserRepo) Update(companyID string, user domain.User) (domain.UserClientCatering, int, error) {
	companyType := utils.DerefString(user.CompanyType)
	var updatedUser domain.UserClientCatering

	if companyType == types.CompanyTypesEnum.Catering {
		if err := config.DB.
			Model(&domain.User{}).
			Update(&user).
			Select("users.*, c.id as catering_id, c.name as catering_name, ci.id as client_id, ci.name as client_name").
			Joins("left join caterings c on c.id = users.catering_id").
			Joins("left join clients ci on ci.id = users.client_id").
			Where("users.catering_id = ?", companyID).
			Scan(&updatedUser).
			Error; err != nil {

			if gorm.IsRecordNotFoundError(err) {
				return domain.UserClientCatering{}, http.StatusNotFound, errors.New("user not found")
			}
			return domain.UserClientCatering{}, http.StatusBadRequest, err
		}
		return updatedUser, 0, nil
	}

	if err := config.DB.
		Model(&domain.User{}).
		Update(&user).
		Select("users.*, c.id as catering_id, c.name as catering_name, ci.id as client_id, ci.name as client_name").
		Joins("left join caterings c on c.id = users.catering_id").
		Joins("left join clients ci on ci.id = users.client_id").
		Where("users.client_id = ?", companyID).
		Scan(&updatedUser).
		Error; err != nil {

		if gorm.IsRecordNotFoundError(err) {
			return domain.UserClientCatering{}, http.StatusNotFound, errors.New("user not found")
		}
		return domain.UserClientCatering{}, http.StatusBadRequest, err
	}
	return updatedUser, 0, nil
}
