package usecase

import (
	"net/http"
	"time"

	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/mailer"
	"github.com/Aiscom-LLC/meals-api/src/schemes/request"
	"github.com/Aiscom-LLC/meals-api/src/types"
	"github.com/Aiscom-LLC/meals-api/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// CateringUser struct
type CateringUser struct{}

// NewCateringUser returns pointer to caterign
// user struct with all methods
func NewCateringUser() *CateringUser {
	return &CateringUser{}
}

// Add creates user for catering
// @Summary Returns error or 201 status code if success
// @Produce json
// @Accept json
// @Tags caterings users
// @Param id path string false "Catering ID"
// @Param body body request.CateringUser false "Catering user"
// @Success 201 {object} response.UserResponse false "Catering user"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/users [post]
func (cu *CateringUser) Add(c *gin.Context) {
	var path types.PathID
	var body request.CateringUser
	var user domain.User

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	copier.Copy(&user, &body)

	if ok := utils.IsEmailValid(user.Email); !ok {
		utils.CreateError(http.StatusBadRequest, "email is not valid", c)
		return
	}

	parsedID, err := uuid.FromString(path.ID)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	user.Role = types.UserRoleEnum.CateringAdmin
	user.Status = &types.StatusTypesEnum.Invited

	password := utils.GenerateString(10)
	user.Password = utils.HashString(password)

	existingUser, err := userRepo.GetByKey("email", user.Email)

	if gorm.IsRecordNotFoundError(err) {
		userResult, userErr := userRepo.Add(user)

		cateringUser := domain.CateringUser{
			UserID:     userResult.ID,
			CateringID: parsedID,
		}

		if err := cateringUserRepo.Add(cateringUser); err != nil {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}

		userClientCatering, err := userRepo.GetByID(userResult.ID.String())

		if err != nil {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}

		if userErr != nil {
			utils.CreateError(http.StatusBadRequest, userErr.Error(), c)
			return
		}

		go mailer.SendEmail(user, password)
		c.JSON(http.StatusCreated, userClientCatering)
		return
	}

	if existingUser.ID != uuid.Nil {
		utils.CreateError(http.StatusBadRequest, "user with that email already exists", c)
		return
	}
}

// Get return list of catering users
// @Summary Returns list of catering users
// @Tags caterings users
// @Produce json
// @Param id path string false "Catering ID"
// @Param limit query int false "used for pagination"
// @Param page query int false "used for pagination"
// @Param q query string false "used query search"
// @Param role query string false "used for role sort"
// @Success 200 {array} response.UserResponse false "Catering user"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /caterings/{id}/users [get]
func (cu *CateringUser) Get(c *gin.Context) { //nolint:dupl
	var path types.PathID
	var paginationQuery types.PaginationQuery
	var filterQuery types.UserFilterQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&paginationQuery, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&filterQuery, c); err != nil {
		return
	}

	users, total, code, err := cateringUserRepo.Get(path.ID, paginationQuery, filterQuery)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	if paginationQuery.Page == 0 {
		paginationQuery.Page = 1
	}

	c.JSON(http.StatusOK, gin.H{
		"items": users,
		"total": total,
		"page":  paginationQuery.Page,
	})
}

// Delete deletes user of catering
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
func (cu *CateringUser) Delete(c *gin.Context) {
	var path types.PathUser
	var user domain.User

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	parsedUserID, _ := uuid.FromString(path.UserID)
	user.ID = parsedUserID
	user.Status = &types.StatusTypesEnum.Deleted
	deletedAt := time.Now().AddDate(0, 0, 21).Truncate(time.Hour * 24)
	user.DeletedAt = &deletedAt

	ctxUser, _ := c.Get("user")
	ctxUserRole := ctxUser.(domain.User).Role

	if user.ID == ctxUser.(domain.User).ID {
		utils.CreateError(http.StatusBadRequest, "can't delete yourself", c)
		return
	}

	code, err := cateringUserRepo.Delete(path.ID, ctxUserRole, user)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}

// Update updates user of catering
// @Summary Returns error or 200 status code if success
// @Produce json
// @Accept json
// @Tags caterings users
// @Param id path string false "Catering ID"
// @Param userId path string false "User ID"
// @Param body body request.CateringUserUpdate false "Catering user"
// @Success 200 {object} response.UserResponse false "Catering user"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /caterings/{id}/users/{userId} [put]
func (cu *CateringUser) Update(c *gin.Context) { //nolint:dupl
	var path types.PathUser
	var user domain.User

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&user, c); err != nil {
		return
	}

	if user.Email != "" {
		if ok := utils.IsEmailValid(user.Email); !ok {
			utils.CreateError(http.StatusBadRequest, "email is not valid", c)
			return
		}
	}

	parsedUserID, _ := uuid.FromString(path.UserID)
	user.CompanyType = &types.CompanyTypesEnum.Catering
	user.ID = parsedUserID

	code, err := cateringUserRepo.Update(&user)
	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	updatedUser, err := userRepo.GetByID(path.UserID)
	c.JSON(http.StatusOK, updatedUser)
}
