package services

import (
	"errors"
	"github.com/Aiscom-LLC/meals-api/api/api_types"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/types"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// CateringUser struct
type CateringUserService struct{}

// NewCateringUserService returns pointer to catering
// user struct with all methods
func NewCateringUserService() *CateringUserService {
	return &CateringUserService{}
}

func (cu *CateringUserService) Add(path api_types.PathID,  user domain.User) (domain.UserClientCatering, domain.User, string, error, error) {
	parsedID, err := uuid.FromString(path.ID)
	if err != nil {
		return domain.UserClientCatering{}, user, "", err, nil
	}

	user.Role = types.UserRoleEnum.CateringAdmin
	user.Status = &types.StatusTypesEnum.Invited
	user.CompanyType = &types.CompanyTypesEnum.Catering

	password := utils.GenerateString(10)
	user.Password = utils.HashString(password)

	existingUsers, err := userRepo.GetAllByKey("email", user.Email)

	if gorm.IsRecordNotFoundError(err) {
		userResult, userErr := userRepo.Add(user)

		cateringUser := domain.CateringUser{
			UserID:     userResult.ID,
			CateringID: parsedID,
		}

		if err := cateringUserRepo.Add(cateringUser); err != nil {
			return domain.UserClientCatering{}, user, password, err, nil
		}

		userClientCatering, err := userRepo.GetByID(userResult.ID.String())

		return userClientCatering, user, password, userErr, err
	}

	for i := range existingUsers {
		if *existingUsers[i].Status != types.StatusTypesEnum.Deleted {
			return domain.UserClientCatering{}, user, password, errors.New("user with that email already exist"), nil
		}
	}

	user, userErr := userRepo.Add(user)

	cateringUser := domain.CateringUser{
		UserID:     user.ID,
		CateringID: parsedID,
	}

	if err := cateringUserRepo.Add(cateringUser); err != nil {
		return domain.UserClientCatering{}, user, password, err, nil
	}

	userClientCatering, err := userRepo.GetByID(user.ID.String())

	return userClientCatering, user, password, err, userErr
}
