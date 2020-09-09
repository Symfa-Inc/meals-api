package repository

import (
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"net/http"

	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository/enums"
	"github.com/Aiscom-LLC/meals-api/utils"
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

func (ur UserRepo) GetByID(id string) (models.UserClientCatering, error) {
	var user models.UserClientCatering
	err := config.DB.
		Table("users as u").
		Select("u.*").
		Where("u.id = ?", id).
		Scan(&user).
		Error

	company := utils.DerefString(user.CompanyType)

	if company == enums.CompanyTypesEnum.Catering {
		config.DB.
			Table("users as u").
			Select("c.id as catering_id, c.name as catering_name").
			Joins("left join catering_users cu on cu.user_id = u.id").
			Joins("left join caterings c on c.id = cu.catering_id").
			Where("u.id = ?", id).
			Scan(&user)
	} else if company == enums.CompanyTypesEnum.Client {
		config.DB.
			Table("users as u").
			Select("cl.id as client_id, cl.name as client_name, c.name as catering_name, c.id as catering_id, clu.floor").
			Joins("left join client_users clu on clu.user_id = u.id ").
			Joins("left join clients cl on cl.id = clu.client_id ").
			Joins("left join caterings c on c.id = cl.catering_id").
			Where("u.id = ?", id).
			Scan(&user)
	}
	return user, err
}

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
