package repository

import (
	"errors"
	"net/http"
	"time"

	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/types"
	"github.com/jinzhu/gorm"
)

// CateringUserRepo struct
type CateringUserRepo struct{}

// NewCateringUserRepo returns pointer to user repository
// with all methods
func NewCateringUserRepo() *CateringUserRepo {
	return &CateringUserRepo{}
}

func (cur *CateringUserRepo) GetByKey(key, value string) (domain.CateringUser, error) {
	var user domain.CateringUser
	err := config.DB.
		Unscoped().
		Where(key+" = ?", value).
		First(&user).Error
	return user, err
}

func (cur *CateringUserRepo) Add(cateringUser domain.CateringUser) error {
	err := config.DB.
		Create(&cateringUser).
		Error
	return err
}

func (cur *CateringUserRepo) Get(cateringID string, pagination types.PaginationQuery, filters types.UserFilterQuery) ([]domain.UserClientCatering, int, int, error) {
	var users []domain.UserClientCatering
	var total int
	page := pagination.Page
	limit := pagination.Limit
	role := filters.Role
	clientName := filters.ClientName
	querySearch := filters.Query

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	if clientName != "" {
		clientName = "AND cl.name ILIKE " + "'%" + clientName + "%'"
	}

	if err := config.DB.
		Debug().
		Unscoped().
		Table("users as u").
		Joins("left join catering_users cu on cu.user_id = u.id").
		Joins("left join caterings c on c.id = cu.catering_id").
		Joins("left join client_users clu on clu.user_id  = u.id ").
		Joins("left join clients cl on cl.id = clu.client_id ").
		Where("cu.catering_id = ? AND (first_name || last_name) ILIKE ?"+
			" AND CAST(u.role AS text) ILIKE ? "+clientName+
			" AND (u.deleted_at > ? OR u.deleted_at IS NULL)", cateringID, "%"+querySearch+"%", "%"+role+"%", time.Now()).
		Count(&total).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, 0, http.StatusNotFound, errors.New("user not found")
		}
		return nil, 0, http.StatusBadRequest, err
	}

	if err := config.DB.
		Debug().
		Unscoped().
		Limit(limit).
		Offset((page-1)*limit).
		Table("users as u").
		Select("u.*, c.id as catering_id, c.name as catering_name, cl.id as client_id, cl.name as client_name").
		Joins("left join catering_users cu on cu.user_id = u.id").
		Joins("left join caterings c on c.id = cu.catering_id").
		Joins("left join client_users clu on clu.user_id  = u.id ").
		Joins("left join clients cl on cl.id = clu.client_id ").
		Where("cu.catering_id = ? AND (first_name || last_name) ILIKE ?"+
			" AND CAST(u.role AS text) ILIKE ? "+clientName+
			" AND (u.deleted_at > ? OR u.deleted_at IS NULL)", cateringID, "%"+querySearch+"%", "%"+role+"%", time.Now()).
		Order("created_at DESC, first_name ASC").
		Scan(&users).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, 0, http.StatusNotFound, errors.New("user not found")
		}
		return nil, 0, http.StatusBadRequest, err
	}

	for i := range users {
		users[i].CompanyType = &types.CompanyTypesEnum.Catering
	}

	return users, total, 0, nil
}
