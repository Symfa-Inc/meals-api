package api

import (
	"errors"
	"github.com/Aiscom-LLC/meals-api/api/api_types"
	"github.com/Aiscom-LLC/meals-api/api/swagger"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/mailer"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"github.com/Aiscom-LLC/meals-api/services"
	"github.com/Aiscom-LLC/meals-api/types"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

// CateringUser struct
type CateringUser struct{}

// NewCateringUser returns pointer to catering
// user struct with all methods
func NewCateringUser() *CateringUser {
	return &CateringUser{}
}

var cateringUserService = services.NewCateringUserService()
var cateringUserRepo = repository.NewCateringUserRepo()

// Add creates user for catering
// @Summary Returns error or 201 status code if success
// @Produce json
// @Accept json
// @Tags caterings users
// @Param id path string false "Catering ID"
// @Param body body request.CateringUser false "Catering user"
// @Success 201 {object} response.UserResponse false "Catering user"
// @Failure 400 {object} api_types.Error "Error"
// @Router /caterings/{id}/users [post]
func (cu *CateringUser) Add(c *gin.Context) {
	var path api_types.PathID
	var body models.CateringUser
	var user domain.User

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	if err := copier.Copy(&user, &body); err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	if ok := utils.IsEmailValid(user.Email); !ok {
		utils.CreateError(http.StatusBadRequest, errors.New("email is not valid"), c)
		return
	}

	userClientCatering, userCreated, password, err, userErr := cateringUserService.Add(path, user)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	if userErr != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	// nolint:errcheck
	go mailer.SendEmail(userCreated, password)
	c.JSON(http.StatusCreated, userClientCatering)
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
// @Failure 400 {object} api_types.Error "Error"
// @Failure 404 {object} api_types.Error "Error"
// @Router /caterings/{id}/users [get]
func (cu *CateringUser) Get(c *gin.Context) { //nolint:dupl
	var path api_types.PathID
	var paginationQuery api_types.PaginationQuery
	var filterQuery api_types.UserFilterQuery

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
		utils.CreateError(code, err, c)
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
// @Failure 400 {object} api_types.Error "Error"
// @Failure 404 {object} api_types.Error "Error"
// @Router /caterings/{id}/users/{userId} [delete]
func (cu *CateringUser) Delete(c *gin.Context) {
	var path api_types.PathUser
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
		utils.CreateError(http.StatusBadRequest, errors.New("can't delete yourself"), c)
		return
	}

	code, err := cateringUserRepo.Delete(path.ID, ctxUserRole, user)

	if err != nil {
		utils.CreateError(code, err, c)
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
// @Failure 400 {object} api_types.Error "Error"
// @Failure 404 {object} api_types.Error "Error"
// @Router /caterings/{id}/users/{userId} [put]
func (cu *CateringUser) Update(c *gin.Context) { //nolint:dupl
	var path api_types.PathUser
	var body swagger.CateringUserUpdate
	var user domain.User

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	if body.Email != "" {
		if ok := utils.IsEmailValid(body.Email); !ok {
			utils.CreateError(http.StatusBadRequest, errors.New("email is not valid"), c)
			return
		}
	}

	if err := copier.Copy(&user, &body); err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	parsedUserID, _ := uuid.FromString(path.UserID)
	user.CompanyType = &types.CompanyTypesEnum.Catering
	user.ID = parsedUserID

	code, err := cateringUserRepo.Update(&user)
	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	updatedUser, _ := userRepo.GetByID(path.UserID)
	c.JSON(http.StatusOK, updatedUser)
}
