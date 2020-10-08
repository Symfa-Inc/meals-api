package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/repository/enums"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// ClientUser struct
type ClientUser struct{}

// NewClientUser return pointer to ClientUser struct
// with all methods
func NewClientUser() *ClientUser {
	return &ClientUser{}
}

var userRepo = repository.NewUserRepo()
var clientUserRepo = repository.NewClientUserRepo()

func (cu *ClientUser) Add(path url.PathID, body models.ClientUser, user domain.User) (models.UserClientCatering, int, string, error, error) {
	parsedID, _ := uuid.FromString(path.ID)
	user.CompanyType = &enums.CompanyTypesEnum.Client
	user.Status = &enums.StatusTypesEnum.Invited

	password := utils.GenerateString(10)
	user.Password = utils.HashString(password)

	existingUser, err := userRepo.GetAllByKey("email", user.Email)
	if gorm.IsRecordNotFoundError(err) {
		user, userErr := userRepo.Add(user)

		clientUser := domain.ClientUser{
			UserID:   user.ID,
			ClientID: parsedID,
			Floor:    body.Floor,
		}

		if err := clientUserRepo.Add(clientUser); err != nil {
			return models.UserClientCatering{}, http.StatusBadRequest, password, err, nil
		}

		userClientCatering, err := userRepo.GetByID(user.ID.String())

		return userClientCatering, 0, password, err, userErr
	}

	for i := range existingUser {
		if *existingUser[i].Status != enums.StatusTypesEnum.Deleted {
			return models.UserClientCatering{}, http.StatusBadRequest, password, errors.New("user with that email already exist"), nil
		}
	}

	user, userErr := userRepo.Add(user)

	clientUser := domain.ClientUser{
		UserID:   user.ID,
		ClientID: parsedID,
		Floor:    body.Floor,
	}

	if err := clientUserRepo.Add(clientUser); err != nil {
		return models.UserClientCatering{}, http.StatusBadRequest, password, err, userErr
	}

	userClientCatering, err := userRepo.GetByID(user.ID.String())

	return userClientCatering, 0, password, err, userErr
}

func (cu *ClientUser) Delete(path url.PathUser, user domain.User, userRole string, userID string) (int, error) {
	parsedUserID, _ := uuid.FromString(path.UserID)
	user.ID = parsedUserID
	user.Status = &enums.StatusTypesEnum.Deleted
	deleteAt := time.Now().AddDate(0, 0, 21).Truncate(time.Hour * 24)
	user.DeletedAt = &deleteAt

	if user.ID.String() == userID {
		return http.StatusBadRequest, errors.New("can't delete yourself")
	}

	code, err := clientUserRepo.Delete(path.ID, userRole, user)

	return code, err
}

func (cu *ClientUser) Update(path url.PathUser, body models.ClientUserUpdate, user domain.User) (int, error) {
	if body.Email != "" {
		if ok := utils.IsEmailValid(body.Email); !ok {
			return http.StatusBadRequest, errors.New("email is not valid")
		}
	}

	if err := copier.Copy(&user, &body); err != nil {
		return http.StatusBadRequest, err
	}

	parsedID, _ := uuid.FromString(path.UserID)
	user.CompanyType = &enums.CompanyTypesEnum.Client
	user.ID = parsedID

	code, err := clientUserRepo.Update(&user, body.Floor)

	return code, err
}
