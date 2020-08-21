package usecase

import (
	"net/http"
	"strings"
	"time"

	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/Aiscom-LLC/meals-api/src/schemes/request"
	"github.com/Aiscom-LLC/meals-api/src/types"
	"github.com/Aiscom-LLC/meals-api/src/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Category struct
type Category struct{}

// NewCategory returns pointer to Category struct
// with all methods
func NewCategory() *Category {
	return &Category{}
}

var categoryRepo = repository.NewCategoryRepo()

// Add add dish category in DB
// returns 200 if success and 4xx if request failed
// @Summary Returns error if exists and 200 if success
// @Produce json
// @Accept json
// @Tags catering categories
// @Param id path string true "Catering ID"
// @Param clientId path string true "Client ID"
// @Param body body request.AddCategory false "Category Name"
// @Success 200 {object} domain.Category false "category object"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/clients/{clientId}/categories [post]
func (dc Category) Add(c *gin.Context) {
	var body request.AddCategory
	var path types.PathClient

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
		Name:       strings.ToLower(body.Name),
		CateringID: cateringID,
		ClientID:   clientID,
	}

	err := categoryRepo.Add(&category)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
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
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/clients/{clientId}/categories/{categoryID} [delete]
func (dc Category) Delete(c *gin.Context) {
	var path types.PathCategory
	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if code, err := categoryRepo.Delete(path); err != nil {
		utils.CreateError(code, err.Error(), c)
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
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/clients/{clientId}/categories [get]
func (dc Category) Get(c *gin.Context) {
	var path types.PathClient
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	if query.Date == "" {
		query.Date = time.Now().Format(time.RFC3339)
	}

	categoriesResult, code, err := categoryRepo.Get(path.ID, path.ClientID, query.Date)
	if err != nil {
		utils.CreateError(code, err.Error(), c)
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
// @Param body body request.UpdateCategory false "new category name"
// @Success 204 "Successfully updated"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/clients/{clientId}/categories/{categoryID} [put]
func (dc Category) Update(c *gin.Context) {
	var path types.PathCategory
	var category domain.Category

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&category, c); err != nil {
		return
	}

	code, err := categoryRepo.Update(path, &category)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}
