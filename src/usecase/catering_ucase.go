package usecase

import (
	"github.com/gin-gonic/gin"
	"go_api/src/delivery/middleware"
	"go_api/src/domain"
	"go_api/src/repository"
	"go_api/src/schemes/response"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
)

// Catering struct
type Catering struct{}

// NewCatering returns pointer to catering struct
// with all methods
func NewCatering() *Catering {
	return &Catering{}
}

var cateringRepo = repository.NewCateringRepo()

// Add creates catering
// @Summary Returns error or 201 status code if success
// @Produce json
// @Accept json
// @Tags catering
// @Param body body request.AddName false "Catering Name"
// @Success 201 {object} domain.Catering false "catering object"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings [post]
func (ca Catering) Add(c *gin.Context) {
	var catering domain.Catering

	if err := utils.RequestBinderBody(&catering, c); err != nil {
		return
	}

	err := cateringRepo.Add(&catering)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusCreated, catering)
}

// Delete soft delete of catering
// @Summary Soft delete
// @Tags catering
// @Produce json
// @Param id path string true "Catering ID"
// @Success 204 "Successfully deleted"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id} [delete]
func (ca Catering) Delete(c *gin.Context) {
	var path types.PathID
	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := cateringRepo.Delete(path.ID); err != nil {
		utils.CreateError(http.StatusNotFound, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetByID returns catering
// @Summary Returns info about catering
// @Tags catering
// @Produce json
// @Param id path string true "Catering ID"
// @Success 200 {object} domain.Catering "catering model"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id} [get]
func (ca Catering) GetByID(c *gin.Context) {
	var path types.PathID

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	result, err := cateringRepo.GetByKey("id", path.ID)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, result)
}

// Get return list of caterings
// @Summary Returns list of caterings
// @Tags catering
// @Produce json
// @Param limit query int false "used for pagination"
// @Param page query int false "used for pagination"
// @Success 200 {object} response.GetCaterings "List of caterings"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings [get]
func (ca Catering) Get(c *gin.Context) {
	var query types.PaginationQuery
	var cateringID string

	if err := utils.RequestBinderQuery(&query, c); err != nil {
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

	if user.CateringID == nil {
		cateringID = ""
	} else {
		cateringID = user.CateringID.String()
	}

	result, total, err := cateringRepo.Get(cateringID, query)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}

	c.JSON(http.StatusOK, response.GetCaterings{
		Items: result,
		Page:  query.Page,
		Total: total,
	})
}

// Update updates catering
// @Summary Returns 204 if success and 4xx error if failed
// @Produce json
// @Accept json
// @Tags catering
// @Param id path string true "Catering ID"
// @Param body body request.AddName false "Catering Name"
// @Success 204 "Successfully updated"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id} [put]
func (ca Catering) Update(c *gin.Context) {
	var path types.PathID
	var cateringModel domain.Catering

	if err := utils.RequestBinderBody(&cateringModel, c); err != nil {
		return
	}

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if code, err := cateringRepo.Update(path.ID, cateringModel); err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}
