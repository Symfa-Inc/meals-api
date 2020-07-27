package usecase

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go_api/src/domain"
	"go_api/src/repository"
	"go_api/src/schemes/request"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
	"strings"
	"time"
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
// @Param body body request.AddCategory false "Category Name"
// @Success 200 {object} domain.Category false "category object"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/categories [post]
func (dc Category) Add(c *gin.Context) {
	var body request.AddCategory
	var path types.PathID

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	if body.DeletedAt != nil {
		if body.DeletedAt.Sub(time.Now()).Hours() < 0 {
			utils.CreateError(http.StatusBadRequest, "can't create dish category with already passed deletedAt date", c)
			return
		}
	}

	cateringID, _ := uuid.FromString(path.ID)
	category := domain.Category{
		DeletedAt:  body.DeletedAt,
		Name:       strings.ToLower(body.Name),
		CateringID: cateringID,
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
// @Param categoryID path string true "Category ID"
// @Success 204 "Successfully deleted"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/categories/{categoryID} [delete]
func (dc Category) Delete(c *gin.Context) {
	var path types.PathCategory
	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := categoryRepo.Delete(path); err != nil {
		utils.CreateError(http.StatusNotFound, err.Error(), c)
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
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/categories [get]
func (dc Category) Get(c *gin.Context) {
	var path types.PathID

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	categoriesResult, code, err := categoryRepo.Get(path.ID)
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
// @Router /caterings/{id}/categories/{categoryID} [put]
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
