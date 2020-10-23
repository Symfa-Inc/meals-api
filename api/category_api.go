package api

import (
	"net/http"

	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"github.com/Aiscom-LLC/meals-api/utils"
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
// @Param body body swagger.AddCategory false "Category Name"
// @Success 200 {object} Category false "category object"
// @Failure 400 {object} Error "Error"
// @Router /caterings/{id}/categories [post]
func (dc Category) Add(c *gin.Context) {
	var body models.AddCategory
	var path url.PathID

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	cateringID, _ := uuid.FromString(path.ID)
	category := domain.Category{
		Name:       body.Name,
		CateringID: cateringID,
	}

	err := categoryRepo.Add(&category)

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
// @Param categoryID path string true "Category ID"
// @Success 204 "Successfully deleted"
// @Failure 404 {object} Error "Not Found"
// @Router /caterings/{id}/categories/{categoryID} [delete]
func (dc Category) Delete(c *gin.Context) {
	var path url.PathCategory

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if code, err := categoryRepo.Delete(path); err != nil {
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
// @Success 200 {array} domain.Category "array of category readings"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /caterings/{id}/categories [get]
func (dc Category) Get(c *gin.Context) {
	var path url.PathID

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	categoriesResult, code, err := categoryRepo.Get(path.ID)

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
// @Router /caterings/{id}/categories/{categoryID} [put]
func (dc Category) Update(c *gin.Context) {
	var path url.PathCategory
	var category domain.Category

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&category, c); err != nil {
		return
	}

	code, err := categoryRepo.Update(path, &category)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.Status(http.StatusNoContent)
}
