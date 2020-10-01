package repository

import (
	"errors"
	"net/http"
	"time"

	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/schemes/response"
	"github.com/Aiscom-LLC/meals-api/src/types"
	"github.com/jinzhu/gorm"
)

// ClientUserRepo struct
type ClientUserRepo struct{}

// NewClientUserRepo returns pointer to user repository
// with all methods
func NewClientUserRepo() *ClientUserRepo {
	return &ClientUserRepo{}
}

func (cur *ClientUserRepo) Add(clientUser domain.ClientUser) error {
	err := config.DB.
		Create(&clientUser).
		Error
	return err
}

func (cur *ClientUserRepo) Get(clientID, userRole string, pagination types.PaginationQuery, filters types.UserFilterQuery) ([]response.GetClientUser, int, int, error) {
	var users []response.GetClientUser
	var total int
	page := pagination.Page
	limit := pagination.Limit
	role := filters.Role
	querySearch := filters.Query
	userConditional := ""

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	if userRole == types.UserRoleEnum.CateringAdmin {
		userConditional = " AND u.role = " + "'" + types.UserRoleEnum.ClientAdmin + "'"
	}

	if err := config.DB.
		Unscoped().
		Table("users as u").
		Joins("left join client_users cu on cu.user_id = u.id").
		Joins("left join clients c on c.id = cu.client_id").
		Where("cu.client_id = ?"+userConditional+" AND (first_name || last_name) ILIKE ?"+
			" AND CAST(u.role AS text) ILIKE ?"+
			" AND (u.deleted_at > ? OR u.deleted_at IS NULL)", clientID, "%"+querySearch+"%", "%"+role+"%", time.Now()).
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
		Select("u.*, cu.floor, c.id as client_id, c.name as client_name").
		Joins("left join client_users cu on cu.user_id = u.id").
		Joins("left join clients c on c.id = cu.client_id").
		Where("cu.client_id = ?"+userConditional+" AND (first_name || last_name) ILIKE ?"+
			" AND CAST(u.role AS text) ILIKE ? "+
			" AND (u.deleted_at > ? OR u.deleted_at IS NULL)", clientID, "%"+querySearch+"%", "%"+role+"%", time.Now()).
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

func (cur *ClientUserRepo) Delete(clientID, ctxUserRole string, user domain.User) (int, error) {
	var totalUsers int
	if ctxUserRole == types.UserRoleEnum.CateringAdmin || ctxUserRole == types.UserRoleEnum.ClientAdmin {
		config.DB.
			Table("users as u").
			Select("u.*, cu.floor").
			Joins("left join client_users cu on cu.user_id = u.id").
			Joins("left join clients c on c.id = cu.client_id").
			Where("cu.client_id = ? AND u.company_type = ? AND u.status != ? AND u.role = ?",
				clientID, types.CompanyTypesEnum.Client, types.StatusTypesEnum.Deleted, types.UserRoleEnum.ClientAdmin).
			Count(&totalUsers)

		if totalUsers == 1 {
			return http.StatusBadRequest, errors.New("can't delete last admin")
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

func (cur *ClientUserRepo) Update(user *domain.User, floor *int) (int, error) {
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
	if floor != nil {
		config.DB.
			Debug().
			Unscoped().
			Model(&domain.ClientUser{}).
			Where("user_id = ? AND (deleted_at > ? OR deleted_at IS NULL)",
				user.ID, time.Now()).
			Update("floor", *floor)
	}

	return 0, nil
}
