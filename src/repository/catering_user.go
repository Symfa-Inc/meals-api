package repository

import (
	"errors"
	"net/http"
	"time"

	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/schemes/response"
	"github.com/Aiscom-LLC/meals-api/src/types"
	"github.com/Aiscom-LLC/meals-api/src/utils"
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

func (cur *CateringUserRepo) Get(cateringID string, pagination types.PaginationQuery, filters types.UserFilterQuery) ([]response.GetCateringUser, int, int, error) {
	var users []response.GetCateringUser
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
		Select("u.*").
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
	if ctxUserRole != types.UserRoleEnum.SuperAdmin {
		config.DB.
			Table("users as u").
			Joins("left join catering_users cu on cu.user_id = u.id").
			Where("cu.catering_id = ? AND u.status != ?",
				cateringID, types.StatusTypesEnum.Deleted).
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
	var prevUser domain.User
	userStatus := utils.DerefString(user.Status)

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

	config.DB.
		Unscoped().
		Model(&domain.User{}).
		Find(&prevUser).
		Where("AND users.id = ? AND (users.deleted_at > ? OR users.deleted_at IS NULL)",
			user.ID, time.Now())

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

	prevUserStatus := utils.DerefString(prevUser.Status)
	if userStatus == types.StatusTypesEnum.Active && prevUserStatus == types.StatusTypesEnum.Deleted {
		config.DB.
			Unscoped().
			Model(&domain.User{}).
			Update(user).
			Update(map[string]interface{}{
				"DeletedAt": user.DeletedAt,
			}).
			Where("(users.deleted_at > ? OR users.deleted_at IS NULL)", time.Now())
	}
	return 0, nil
}
