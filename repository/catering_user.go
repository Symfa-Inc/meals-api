package repository

import (
	"errors"
	types "github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"net/http"
	"time"

	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository/enums"
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

func (cur *CateringUserRepo) Get(cateringID string, pagination types.PaginationQuery, filters types.UserFilterQuery) ([]models.GetCateringUser, int, int, error) {
	var users []models.GetCateringUser
	var total int
	page := pagination.Page
	limit := pagination.Limit
	role := filters.Role
	querySearch := filters.Query

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	if err := config.DB.
		Unscoped().
		Table("users as u").
		Joins("left join catering_users cu on cu.user_id = u.id").
		Joins("left join caterings c on c.id = cu.catering_id").
		Where("cu.catering_id = ? AND (first_name || last_name) ILIKE ?"+
			" AND CAST(u.role AS text) ILIKE ?"+
			" AND (u.deleted_at > ? OR u.deleted_at IS NULL)", cateringID, "%"+querySearch+"%", "%"+role+"%", time.Now()).
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
		Table("users as u").
		Select("u.*, c.id as catering_id, c.name as catering_name").
		Joins("left join catering_users cu on cu.user_id = u.id").
		Joins("left join caterings c on c.id = cu.catering_id").
		Where("cu.catering_id = ? AND (first_name || last_name) ILIKE ?"+
			" AND CAST(u.role AS text) ILIKE ? "+
			" AND (u.deleted_at > ? OR u.deleted_at IS NULL)", cateringID, "%"+querySearch+"%", "%"+role+"%", time.Now()).
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

func (cur *CateringUserRepo) Delete(cateringID, ctxUserRole string, user domain.User) (int, error) {
	var totalUsers int
	if ctxUserRole != enums.UserRoleEnum.SuperAdmin {
		config.DB.
			Table("users as u").
			Joins("left join catering_users cu on cu.user_id = u.id").
			Where("cu.catering_id = ? AND u.status != ?",
				cateringID, enums.StatusTypesEnum.Deleted).
			Count(&totalUsers)

		if totalUsers == 1 {
			return http.StatusBadRequest, errors.New("can't delete last user")
		}
	}

	if userExist := config.DB.
		Table("users as u").
		Where("u.id = ?", user.ID).
		Update(&user).
		RowsAffected; userExist == 0 {
		return http.StatusNotFound, errors.New("user not found")
	}
	return 0, nil
}

func (cur *CateringUserRepo) Update(user *domain.User) (int, error) {
	if userExist := config.DB.
		Where("id = ? AND email = ?", user.ID, user.Email).
		Find(&domain.User{}).
		RowsAffected; userExist == 0 {
		if emailExist := config.DB.
			Where("email = ?", user.Email).
			Find(&domain.User{}).
			RowsAffected; emailExist != 0 {
			return http.StatusBadRequest, errors.New("user with this email already exists")
		}
	}

	if err := config.DB.
		Unscoped().
		Model(&domain.User{}).
		Update(user).
		Error; err != nil {

		if gorm.IsRecordNotFoundError(err) {
			return http.StatusNotFound, errors.New("user not found")
		}
		return http.StatusBadRequest, err
	}
	return 0, nil
}
