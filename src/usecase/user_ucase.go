package usecase

import (
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/Aiscom-LLC/meals-api/src/schemes/request"
	"github.com/Aiscom-LLC/meals-api/src/types"
	"github.com/Aiscom-LLC/meals-api/src/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

// User struct
type User struct{}

// NewUser returns pointer to user struct
// with all methods
func NewUser() *User {
	return &User{}
}

var userRepo = repository.NewUserRepo()

// UpdateCateringUser updates user of catering
// @Summary Returns error or 200 status code if success
// @Produce json
// @Accept json
// @Tags Users
// @Param id path string false "User ID"
// @Param body body request.UserPasswordUpdate false "User"
// @Success 200 {object} response.UserResponse false "User"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /users/{id} [put]
func (u User) ChangePassword(c *gin.Context) { //nolint:dupl
	var path types.PathID
	var body request.UserPasswordUpdate

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	if len(body.NewPassword) < 10 {
		utils.CreateError(http.StatusBadRequest, "Password must contain at least 10 characters", c)
	}
	newPassword := utils.HashString(body.NewPassword)
	parsedUserID, _ := uuid.FromString(path.ID)

	if body.OldPassword == body.NewPassword {
		utils.CreateError(http.StatusBadRequest, "Passwords are the same", c)
		return
	}
	user, err := userRepo.GetByID(parsedUserID.String())

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	if ok := utils.CheckPasswordHash(body.OldPassword, user.Password); !ok {
		utils.CreateError(http.StatusBadRequest, "Wrong password", c)
		return
	}

	code, err := userRepo.UpdatePassword(user.ID, newPassword)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, "Password updated")
}


/*// AddCateringUser creates user for catering
// @Summary Returns error or 201 status code if success
// @Produce json
// @Accept json
// @Tags caterings users
// @Param id path string false "Catering ID"
// @Param body body request.CateringUser false "Catering user"
// @Success 201 {object} response.UserResponse false "Catering user"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/users [post]
func (u User) AddCateringUser(c *gin.Context) { //nolint:dupl
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
	//user.CateringID = &parsedID
	//user.CompanyType = &types.CompanyTypesEnum.Catering
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

		if userErr != nil {
			utils.CreateError(http.StatusBadRequest, userErr.Error(), c)
			return
		}

		go mailer.SendEmail(user, password)
		c.JSON(http.StatusCreated, userResult)
		return
	}

	if existingUser.ID != uuid.Nil {
		utils.CreateError(http.StatusBadRequest, "user with this email already exists", c)
		return
	}
}

// GetCateringUsers return list of catering users
// @Summary Returns list of catering users
// @Tags caterings users
// @Produce json
// @Param id path string false "Catering ID"
// @Param limit query int false "used for pagination"
// @Param page query int false "used for pagination"
// @Param q query string false "used query search"
// @Param role query string false "used for role sort"
// @Param client query string false "used for client sort"
// @Success 200 {array} response.UserResponse false "Catering user"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /caterings/{id}/users [get]
func (u User) GetCateringUsers(c *gin.Context) { //nolint:dupl
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

	catering := types.CompanyTypesEnum.Catering
	users, total, code, err := userRepo.Get(path.ID, catering, "", paginationQuery, filterQuery)

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
	//user.CompanyType = &types.CompanyTypesEnum.Catering
	user.Status = &types.StatusTypesEnum.Deleted
	deletedAt := time.Now().AddDate(0, 0, 21).Truncate(time.Hour * 24)
	user.DeletedAt = &deletedAt

	ctxUser, _ := c.Get("user")
	ctxUserRole := ctxUser.(domain.UserClientCatering).Role

	if user.ID == ctxUser.(domain.UserClientCatering).ID {
		utils.CreateError(http.StatusBadRequest, "can't delete yourself", c)
		return
	}

	code, err := userRepo.Delete(path.ID, ctxUserRole, user)

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
// @Param body body request.CateringUserUpdate false "Catering user"
// @Success 200 {object} response.UserResponse false "Catering user"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /caterings/{id}/users/{userId} [put]
func (u User) UpdateCateringUser(c *gin.Context) { //nolint:dupl
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

	updatedUser, code, err := userRepo.Update(path.ID, user)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// GetClientUsers return list of client users
// @Summary Returns list of client users
// @Tags clients users
// @Produce json
// @Param id path string false "Client ID"
// @Param limit query int false "used for pagination"
// @Param page query int false "used for pagination"
// @Param q query string false "used query search"
// @Param role query string false "used for role sort"
// @Param status query string false "used for status sort"
// @Success 200 {array} response.UserResponse "List of client users"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /clients/{id}/users [get]
func (u User) GetClientUsers(c *gin.Context) { //nolint:dupl
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

	claims, err := middleware.Passport().GetClaimsFromJWT(c)

	if err != nil {
		utils.CreateError(http.StatusUnauthorized, err.Error(), c)
		return
	}

	id := claims["id"].(string)

	user, err := userRepo.GetByKey("id", id)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	client := types.CompanyTypesEnum.Client
	users, total, code, err := userRepo.Get(path.ID, client, user.Role, paginationQuery, filterQuery)

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

// AddClientUser creates user for client
// @Summary Returns error or 201 status code if success
// @Produce json
// @Accept json
// @Tags clients users
// @Param id path string false "Client ID"
// @Param body body request.ClientUser false "Client user"
// @Success 201 {object} response.UserResponse false "Client user"
// @Failure 400 {object} types.Error "Error"
// @Router /clients/{id}/users [post]
func (u User) AddClientUser(c *gin.Context) { //nolint:dupl
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
	user.ClientID = &parsedID
	user.CompanyType = &types.CompanyTypesEnum.Client
	user.Status = &types.StatusTypesEnum.Invited

	password := utils.GenerateString(10)
	user.Password = utils.HashString(password)

	existingUser, err := userRepo.GetByKey("email", user.Email)

	if gorm.IsRecordNotFoundError(err) {
		if err := mailer.SendEmail(user, password); err != nil {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}

		client, _ := clientRepo.GetByKey("id", path.ID)
		user.CateringID = &client.CateringID
		user, userErr := userRepo.Add(user)

		if userErr != nil {
			utils.CreateError(http.StatusBadRequest, userErr.Error(), c)
			return
		}

		c.JSON(http.StatusCreated, user)
	}

	if existingUser.ID != uuid.Nil {
		utils.CreateError(http.StatusBadRequest, "user with that email already exist", c)
		return
	}
}

// DeleteClientUser deletes user of client
// @Summary Returns error or 204 status code if success
// @Produce json
// @Accept json
// @Tags clients users
// @Param id path string false "Client ID"
// @Param userId path string false "User ID"
// @Success 204 "Successfully deleted"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /clients/{id}/users/{userId} [delete]
func (u User) DeleteClientUser(c *gin.Context) {
	var path types.PathUser
	var user domain.User

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	parsedUserID, _ := uuid.FromString(path.UserID)
	user.ID = parsedUserID
	user.CompanyType = &types.CompanyTypesEnum.Client
	user.Status = &types.StatusTypesEnum.Deleted
	deletedAt := time.Now().AddDate(0, 0, 21).Truncate(time.Hour * 24)
	user.DeletedAt = &deletedAt

	ctxUser, _ := c.Get("user")
	ctxUserRole := ctxUser.(domain.UserClientCatering).Role

	if user.ID == ctxUser.(domain.UserClientCatering).ID {
		utils.CreateError(http.StatusBadRequest, "can't delete yourself", c)
		return
	}

	code, err := userRepo.Delete(path.ID, ctxUserRole, user)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateClientUser updates user of client
// @Summary Returns error or 200 status code if success
// @Produce json
// @Accept json
// @Tags clients users
// @Param id path string false "Client ID"
// @Param userId path string false "User ID"
// @Param body body request.ClientUserUpdate false "Client user"
// @Success 200 {object} response.UserResponse false "Client user"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /clients/{id}/users/{userId} [put]
func (u User) UpdateClientUser(c *gin.Context) { //nolint:dupl
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
	user.CompanyType = &types.CompanyTypesEnum.Client
	user.ID = parsedUserID

	updatedUser, code, err := userRepo.Update(path.ID, user)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}*/
