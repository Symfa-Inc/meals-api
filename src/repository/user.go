package repository

import (
	"net/http"

	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
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
func (ur UserRepo) GetByKey(key, value string) (domain.User, error) {
	var user domain.User
	err := config.DB.
		Unscoped().
		Where(key+" = ?", value).
		First(&user).Error

	return user, err
}

// GetAllByKey returns users by key
// and error if exist
func (ur UserRepo) GetAllByKey(key, value string) ([]domain.User, error) {
	var user []domain.User
	err := config.DB.
		Unscoped().
		Where(key+" = ?", value).
		Find(&user).Error

	return user, err
}

func (ur UserRepo) GetByID(id string) (domain.UserClientCatering, error) {
	var user domain.UserClientCatering
	err := config.DB.
		Table("users as u").
		Select("u.*, cl.id as client_id, cl.name as client_name, c.name as catering_name, c.id as catering_id, clu.floor").
		Joins("left join catering_users cu on cu.user_id = u.id").
		Joins("left join caterings c on c.id = cu.catering_id").
		Joins("left join client_users clu on clu.user_id = u.id ").
		Joins("left join clients cl on cl.id = clu.client_id ").
		Where("u.id = ?", id).
		Scan(&user).
		Error

	return user, err
}

// Get returns list of users for provided company
/*func (ur UserRepo) Get(companyID, companyType, userRole string, pagination types.PaginationQuery, filters types.UserFilterQuery) ([]domain.UserClientCatering, int, int, error) {
	var users []domain.UserClientCatering
	var total int
	page := pagination.Page
	limit := pagination.Limit
	status := filters.Status
	role := filters.Role
	querySearch := filters.Query

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	if userRole == types.UserRoleEnum.CateringAdmin {
		role = types.UserRoleEnum.ClientAdmin
	}
	if companyType == types.CompanyTypesEnum.Catering {
		if role == "" {
			role = types.UserRoleEnum.CateringAdmin
		}
		if clientName != "" {
			if err := config.DB.
				Unscoped().
				Model(&domain.User{}).
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
			Joins("left join caterings c on c.id = users.catering_id").
			Joins("full outer join clients ci on ci.id = users.client_id").
			Where("users.catering_id = ? AND (first_name || last_name) ILIKE ?"+
				" AND CAST(users.role AS text) ILIKE ?"+
				" AND (users.deleted_at > ? OR users.deleted_at IS NULL)", companyID, "%"+querySearch+"%", "%"+role+"%", time.Now()).
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
			Joins("full outer join clients ci on ci.id = users.client_id").
			Where("users.catering_id = ? AND (first_name || last_name) ILIKE ?"+
				" AND CAST(users.role AS text) ILIKE ?"+
				" AND (users.deleted_at > ? OR users.deleted_at IS NULL)", companyID, "%"+querySearch+"%", "%"+role+"%", time.Now()).
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
		Where("users.client_id = ? AND users.company_type = ? AND (first_name || last_name) ILIKE ?"+
			"AND status ILIKE ? AND CAST(users.role AS text) ILIKE ?"+
			" AND (users.deleted_at > ? OR users.deleted_at IS NULL)",
			companyID, types.CompanyTypesEnum.Client, "%"+querySearch+"%", "%"+status+"%", "%"+role+"%", time.Now()).
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
		Where("users.client_id = ? AND users.company_type = ? AND (first_name || last_name) ILIKE ?"+
			"AND status ILIKE ? AND CAST(users.role AS text) ILIKE ?"+
			" AND (users.deleted_at > ? OR users.deleted_at IS NULL)",
			companyID, types.CompanyTypesEnum.Client, "%"+querySearch+"%", "%"+status+"%", "%"+role+"%", time.Now()).
		Order("created_at DESC, first_name ASC").
		Scan(&users).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, 0, http.StatusNotFound, errors.New("user not found")
		}
		return nil, 0, http.StatusBadRequest, err
	}
	return users, total, 0, nil
}*/

// Add adds new user for certain company passed in user struct
// returns user and error
func (ur UserRepo) Add(user domain.User) (domain.User, error) {
	if err := config.DB.
		Create(&user).
		Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// Delete changes status of user to deleted
// and sets deleted_at field to 21 days from now
/*func (ur UserRepo) Delete(companyID, ctxUserRole string, user domain.User) (int, error) {
	var totalUsers int
	companyType := utils.DerefString(user.CompanyType)
	if companyType == types.CompanyTypesEnum.Catering {
		if ctxUserRole != types.UserRoleEnum.SuperAdmin {
			config.DB.
				Model(&domain.User{}).
				Where("catering_id = ? AND company_type = ? AND status != ?",
					companyID, types.CompanyTypesEnum.Catering, types.StatusTypesEnum.Deleted).
				Count(&totalUsers)

			if totalUsers == 1 {
				return http.StatusNotFound, errors.New("can't delete last user")
			}
		}

		if userExist := config.DB.
			Model(&domain.User{}).
			Where("catering_id = ?", companyID).
			Update(&user).
			RowsAffected; userExist == 0 {
			return http.StatusBadRequest, errors.New("user not found")
		}
		return 0, nil
	}

	userRole, _ := ur.GetByKey("id", user.ID.String())
	if userRole.Role == types.UserRoleEnum.ClientAdmin && ctxUserRole != types.UserRoleEnum.SuperAdmin {
		config.DB.
			Model(&domain.User{}).
			Where("client_id = ? AND company_type = ? AND status != ? AND role = ?",
				companyID, types.CompanyTypesEnum.Client, types.StatusTypesEnum.Deleted, types.UserRoleEnum.ClientAdmin).
			Count(&totalUsers)

		if totalUsers == 1 {
			return http.StatusBadRequest, errors.New("can't delete last admin")
		}
	}

	if userExist := config.DB.
		Model(&domain.User{}).
		Where("client_id = ? AND company_type = ?", companyID, types.CompanyTypesEnum.Client).
		Update(&user).
		RowsAffected; userExist == 0 {
		return http.StatusBadRequest, errors.New("user not found")
	}
	return 0, nil
}*/

// Update updates user for passed company ID
// checks if user belongs to client or catering
/*func (ur UserRepo) Update(companyID string, user domain.User) (domain.UserClientCatering, int, error) {
	companyType := utils.DerefString(user.CompanyType)
	var prevUser domain.User
	userStatus := utils.DerefString(user.Status)
	var updatedUser domain.UserClientCatering

	if userExist := config.DB.
		Where("id = ? AND email = ?", user.ID, user.Email).
		Find(&domain.User{}).
		RowsAffected; userExist == 0 {
		if emailExist := config.DB.
			Where("email = ?", user.Email).
			Find(&domain.User{}).
			RowsAffected; emailExist != 0 {
			return domain.UserClientCatering{}, http.StatusBadRequest, errors.New("user with this email already exists")
		}
	}

	if companyType == types.CompanyTypesEnum.Catering {
		config.DB.
			Unscoped().
			Model(&domain.User{}).
			Find(&prevUser).
			Where("users.catering_id = ? AND users.id = ? AND (users.deleted_at > ? OR users.deleted_at IS NULL)",
				companyID, user.ID, time.Now())

		if err := config.DB.
			Unscoped().
			Model(&domain.User{}).
			Update(&user).
			Select("users.*, c.id as catering_id, c.name as catering_name, ci.id as client_id, ci.name as client_name").
			Joins("left join caterings c on c.id = users.catering_id").
			Joins("left join clients ci on ci.id = users.client_id").
			Where("users.catering_id = ? AND (users.deleted_at > ? OR users.deleted_at IS NULL)", companyID, time.Now()).
			Scan(&updatedUser).
			Error; err != nil {

			if gorm.IsRecordNotFoundError(err) {
				return domain.UserClientCatering{}, http.StatusNotFound, errors.New("user not found")
			}
			return domain.UserClientCatering{}, http.StatusBadRequest, err
		}

		prevUserStatus := utils.DerefString(prevUser.Status)
		if userStatus == types.StatusTypesEnum.Active && prevUserStatus == types.StatusTypesEnum.Deleted {
			config.DB.
				Unscoped().
				Model(&domain.User{}).
				Update(&user).
				Update(map[string]interface{}{
					"DeletedAt": user.DeletedAt,
				}).
				Where("users.catering_id = ? AND (users.deleted_at > ? OR users.deleted_at IS NULL)", companyID, time.Now())
		}
		return updatedUser, 0, nil
	}

	config.DB.
		Unscoped().
		Model(&domain.User{}).
		Where("users.client_id = ? AND users.id = ? AND (users.deleted_at > ? OR users.deleted_at IS NULL)",
			companyID, user.ID, time.Now()).
		Find(&prevUser)

	if err := config.DB.
		Unscoped().
		Model(&domain.User{}).
		Update(&user).
		Select("users.*, c.id as catering_id, c.name as catering_name, ci.id as client_id, ci.name as client_name").
		Joins("left join caterings c on c.id = users.catering_id").
		Joins("left join clients ci on ci.id = users.client_id").
		Where("users.client_id = ? AND (users.deleted_at > ? OR users.deleted_at IS NULL)", companyID, time.Now()).
		Scan(&updatedUser).
		Error; err != nil {

		if gorm.IsRecordNotFoundError(err) {
			return domain.UserClientCatering{}, http.StatusNotFound, errors.New("user not found")
		}
		return domain.UserClientCatering{}, http.StatusBadRequest, err
	}

	prevUserStatus := utils.DerefString(prevUser.Status)
	if userStatus == types.StatusTypesEnum.Active && prevUserStatus == types.StatusTypesEnum.Deleted {
		config.DB.
			Unscoped().
			Model(&domain.User{}).
			Update(&user).
			Update(map[string]interface{}{
				"DeletedAt": user.DeletedAt,
			}).
			Where("users.client_id = ? AND (users.deleted_at > ? OR users.deleted_at IS NULL)", companyID, time.Now())
	}

	return updatedUser, 0, nil
}*/

// UpdateStatus updates status for provided userID
func (ur UserRepo) UpdateStatus(userID uuid.UUID, status string) (int, error) {
	if err := config.DB.
		Model(&domain.User{}).
		Where("id = ?", userID).
		Update("status", status).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return http.StatusNotFound, err
		}
		return http.StatusBadRequest, err
	}

	return 0, nil
}
