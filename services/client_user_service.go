package services

import (
	"errors"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/schemes/request"
	"github.com/Aiscom-LLC/meals-api/types"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
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

func (cu *ClientUser) Add(path types.PathID, body request.ClientUser, user domain.User) (domain.UserClientCatering, int, error, error, string) {

	parsedID, _ := uuid.FromString(path.ID)
	user.CompanyType = &types.CompanyTypesEnum.Client
	user.Status = &types.StatusTypesEnum.Invited

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
			return domain.UserClientCatering{}, http.StatusBadRequest, err, nil, password
		}

		userClientCatering, err := userRepo.GetByID(user.ID.String())

		return userClientCatering, 0, err, userErr, password
	}

	for i := range existingUser {
		if *existingUser[i].Status != types.StatusTypesEnum.Deleted {
			return domain.UserClientCatering{}, http.StatusBadRequest, errors.New("user with that email already exist"), nil, password
		}
	}
	user, userErr := userRepo.Add(user)

	clientUser := domain.ClientUser{
		UserID:   user.ID,
		ClientID: parsedID,
		Floor:    body.Floor,
	}

	if err := clientUserRepo.Add(clientUser); err != nil {
		return domain.UserClientCatering{}, http.StatusBadRequest, err, userErr, password
	}

	userClientCatering, err := userRepo.GetByID(user.ID.String())

	return userClientCatering, 0, err, userErr, password
}

func (cu *ClientUser) Delete(path types.PathUser, user domain.User, delUser interface{}) (int, error) {
	parsedUserID, _ := uuid.FromString(path.UserID)
	user.ID = parsedUserID
	user.Status = &types.StatusTypesEnum.Deleted
	deleteAt := time.Now().AddDate(0, 0, 21).Truncate(time.Hour * 24)
	user.DeletedAt = &deleteAt

	userRole := delUser.(domain.User).Role
	if user.ID == delUser.(domain.User).ID {
		return http.StatusBadRequest, errors.New("can't delete yourself")
	}

	code, err := clientUserRepo.Delete(path.ID, userRole, user)

	return code, err
}

func (cu *ClientUser) Update(path types.PathUser, body request.ClientUserUpdate, user domain.User) (int, error) {

	if body.Email != "" {
		if ok := utils.IsEmailValid(body.Email); !ok {
			return http.StatusBadRequest, errors.New("email is not valid")
		}
	}

	if err := copier.Copy(&user, &body); err != nil {
		return http.StatusBadRequest, err

	}

	parsedID, _ := uuid.FromString(path.UserID)
	user.CompanyType = &types.CompanyTypesEnum.Client
	user.ID = parsedID

	code, err := clientUserRepo.Update(&user, body.Floor)

	return code, err
}
