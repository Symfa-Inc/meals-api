package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"go_api/src/domain"
	"go_api/src/mailer"
	"go_api/src/repository"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
	"time"
)

// User struct
type User struct{}

// NewUser returns pointer to user struct
// with all methods
func NewUser() *User {
	return &User{}
}

var userRepo = repository.NewUserRepo()

// AddCateringUser creates user for catering
// @Summary Returns error or 201 status code if success
// @Produce json
// @Accept json
// @Tags caterings users
// @Param id path string false "Catering ID"
// @Param body body request.CateringUser false "Catering user"
// @Success 201 {object} domain.User false "Catering user"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/users [post]
func (u User) AddCateringUser(c *gin.Context) {
	var path types.PathID
	var user domain.User

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&user, c); err != nil {
		return
	}

	if ok := utils.IsEmailValid(user.Email); !ok {
		utils.CreateError(http.StatusBadRequest, "email is not valid", c)
		return
	}

	parsedID, _ := uuid.FromString(path.ID)
	user.CateringID = &parsedID
	user.CompanyType = &types.CompanyTypesEnum.Catering
	user.Role = types.UserRoleEnum.CateringAdmin
	user.Status = &types.StatusTypesEnum.Invited

	password := utils.GenerateString(10)
	user.Password = utils.HashString(password)

	_, err := userRepo.GetByKey("email", user.Email)

	if gorm.IsRecordNotFoundError(err) {
		if err := mailer.SendEmail(user, password); err != nil {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}

		user, userErr := userRepo.Add(user)

		if userErr != nil {
			utils.CreateError(http.StatusBadRequest, userErr.Error(), c)
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

// DeleteCateringUser deletes user of catering
// @Summary Returns error or 204 status code if success
// @Produce json
// @Accept json
// @Tags caterings users
// @Param id path string false "Catering ID"
// @Param userId path string false "User ID"
// @Success 204 "Successfully deleted"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /caterings/{id}/users/{userId} [delete]
func (u User) DeleteCateringUser(c *gin.Context) {
	var path types.PathUser
	var user domain.User

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	parsedUserID, _ := uuid.FromString(path.UserID)
	user.ID = parsedUserID
	user.CompanyType = &types.CompanyTypesEnum.Catering
	user.Status = &types.StatusTypesEnum.Deleted
	deletedAt := time.Now().AddDate(0, 0, 21).Truncate(time.Hour * 24)
	user.DeletedAt = &deletedAt

	code, err := userRepo.Delete(path.ID, user)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateCateringUser updates user of catering
// @Summary Returns error or 200 status code if success
// @Produce json
// @Accept json
// @Tags caterings users
// @Param id path string false "Catering ID"
// @Param userId path string false "User ID"
// @Param body body request.CateringUser false "Catering user"
// @Success 201 {object} domain.User false "Catering user"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /caterings/{id}/users/{userId} [put]
func (u User) UpdateCateringUser(c *gin.Context) {
	var path types.PathUser
	var user domain.User

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&user, c); err != nil {
		return
	}

	if ok := utils.IsEmailValid(user.Email); !ok {
		utils.CreateError(http.StatusBadRequest, "email is not valid", c)
		return
	}

	parsedUserID, _ := uuid.FromString(path.UserID)
	user.CompanyType = &types.CompanyTypesEnum.Catering
	user.ID = parsedUserID

	user, code, err := userRepo.Update(path.ID, user)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, user)
}
