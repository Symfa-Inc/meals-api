package api

import (
	"net/http"
	"time"

	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// ClientCategory struct
type ClientCategory struct{}

// NewClientCategory returns pointer to Category struct
// with all methods
func NewClientCategory() *ClientCategory {
	return &ClientCategory{}
}

var clientCategoryRepo = repository.NewClientCategoryRepo()

// Add add dish category in DB
// returns 200 if success and 4xx if request failed
// @Summary Returns error if exists and 200 if success
// @Produce json
// @Accept json
// @Tags catering categories
// @Param id path string true "Catering ID"
// @Param clientId path string true "Client ID"
// @Param body body swagger.AddCategory false "Category Name"
// @Success 200 {object} Category false "category object"
// @Failure 400 {object} Error "Error"
// @Router /caterings/{id}/clients/{clientId}/categories [post]
func (dc ClientCategory) Add(c *gin.Context) {
	var body models.AddCategory
	var path url.PathClient

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	cateringID, _ := uuid.FromString(path.ID)
	clientID, _ := uuid.FromString(path.ClientID)
	category := domain.Category{
		Date:       body.Date,
		Name:       body.Name,
		CateringID: cateringID,
		ClientID:   clientID,
	}

	err := clientCategoryRepo.Add(&category)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	c.JSON(http.StatusOK, category)
}

// Delete soft delete of category reading
// @Summary Soft delete
// @Tags catering categories
// @Produce json
// @Param id path string true "Catering ID"
// @Param clientId path string false "Client ID"
// @Param categoryID path string true "Category ID"
// @Success 204 "Successfully deleted"
// @Failure 404 {object} Error "Not Found"
// @Router /caterings/{id}/clients/{clientId}/categories/{categoryID} [delete]
func (dc ClientCategory) Delete(c *gin.Context) {
	var path url.PathCategory

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if code, err := clientCategoryRepo.Delete(path); err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.Status(http.StatusNoContent)
}

// Get returns list of categories or error
// @Summary Get list of categories
// @Tags catering categories
// @Produce json
// @Param id path string false "Catering ID"
// @Param clientId path string false "Client ID"
// @Param date query string false "in format YYYY-MM-DDT00:00:00Z"
// @Success 200 {array} domain.Category "array of category readings"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /caterings/{id}/clients/{clientId}/categories [get]
func (dc ClientCategory) Get(c *gin.Context) {
	var path url.PathClient
	var query url.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	if query.Date == "" {
		query.Date = time.Now().Format(time.RFC3339)
	}

	categoriesResult, code, err := clientCategoryRepo.Get(path.ID, path.ClientID, query.Date)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, categoriesResult)
}

// Update updates dish category with new value provided in body
// @Summary Returns 204 if success and 4xx error if failed
// @Produce json
// @Accept json
// @Tags catering categories
// @Param id path string true "Catering ID"
// @Param categoryID path string true "Category ID"
// @Param body body swagger.UpdateCategory false "new category name"
// @Success 204 "Successfully updated"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /caterings/{id}/clients/{clientId}/categories/{categoryID} [put]
func (dc ClientCategory) Update(c *gin.Context) {
	var path url.PathCategory
	var category domain.Category

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&category, c); err != nil {
		return
	}

	code, err := clientCategoryRepo.Update(path, &category)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.Status(http.StatusNoContent)
}
