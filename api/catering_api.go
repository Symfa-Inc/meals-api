package api

import (
	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/api/swagger"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/services"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Catering struct
type Catering struct{}

// NewCatering returns pointer to catering struct
// with all methods
func NewCatering() *Catering {
	return &Catering{}
}

var cateringService = services.NewCateringService()
var cateringRepo = repository.NewCateringRepo()

// Get return list of caterings
// @Summary Returns list of caterings
// @Tags catering
// @Produce json
// @Param limit query int false "used for pagination"
// @Param page query int false "used for pagination"
// @Success 200 {object} swagger.GetCaterings "List of caterings"
// @Failure 400 {object} Error "Error"
// @Router /caterings [get]
func (ca Catering) Get(c *gin.Context) {
	var query PaginationQuery

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	claims, err := middleware.Passport().GetClaimsFromJWT(c)

	if err != nil {
		return
	}

	caterings, total, code, err := cateringService.Get(jwt.MapClaims(claims), &query)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, swagger.GetCaterings{
		Items: caterings,
		Page:  query.Page,
		Total: total,
	})
}

// Add creates catering
// @Summary Returns error or 201 status code if success
// @Produce json
// @Accept json
// @Tags catering
// @Param body body swagger.AddName false "Catering Name"
// @Success 201 {object} domain.Catering false "catering object"
// @Failure 400 {object} Error "Error"
// @Router /caterings [post]
func (ca Catering) Add(c *gin.Context) {
	var catering domain.Catering

	if err := utils.RequestBinderBody(&catering, c); err != nil {
		return
	}

	err := cateringRepo.Add(&catering)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
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
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /caterings/{id} [delete]
func (ca Catering) Delete(c *gin.Context) {
	var path PathID

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := cateringRepo.Delete(path.ID); err != nil {
		utils.CreateError(http.StatusNotFound, err, c)
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
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /caterings/{id} [get]
func (ca Catering) GetByID(c *gin.Context) {
	var path PathID

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	result, err := cateringRepo.GetByKey("id", path.ID)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	c.JSON(http.StatusOK, result)
}

// Update updates catering
// @Summary Returns 204 if success and 4xx error if failed
// @Produce json
// @Accept json
// @Tags catering
// @Param id path string true "Catering ID"
// @Param body body swagger.AddName false "Catering Name"
// @Success 204 "Successfully updated"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /caterings/{id} [put]
func (ca Catering) Update(c *gin.Context) {
	var path PathID
	var cateringModel domain.Catering

	if err := utils.RequestBinderBody(&cateringModel, c); err != nil {
		return
	}

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if code, err := cateringRepo.Update(path.ID, cateringModel); err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.Status(http.StatusNoContent)
}
