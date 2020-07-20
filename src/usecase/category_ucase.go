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

type category struct{}

func NewCategory() *category {
	return &category{}
}

var categoryRepo = repository.NewCategoryRepo()

// AddCategory add dish category in DB
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
func (dc category) Add(c *gin.Context) {
	var body request.AddCategory
	var path types.PathId

	if err := utils.RequestBinderUri(&path, c); err != nil {
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

	cateringId, _ := uuid.FromString(path.ID)
	categoryModel := domain.Category{
		DeletedAt:  body.DeletedAt,
		Name:       strings.ToLower(body.Name),
		CateringID: cateringId,
	}

	category, err := categoryRepo.Add(categoryModel)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory soft delete of category reading
// @Summary Soft delete
// @Tags catering categories
// @Produce json
// @Param id path string true "Catering ID"
// @Param categoryId path string true "Category ID"
// @Success 204 "Successfully deleted"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/categories/{categoryId} [delete]
func (dc category) Delete(c *gin.Context) {
	var path types.PathCategory
	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := categoryRepo.Delete(path); err != nil {
		utils.CreateError(http.StatusNotFound, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetCategories returns list of categories or error
// @Summary Get list of categories
// @Tags catering categories
// @Produce json
// @Param id path string false "Catering ID"
// @Success 200 {array} domain.Category "array of category readings"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/categories [get]
func (dc category) Get(c *gin.Context) {
	var path types.PathId

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	categoriesResult, err, code := categoryRepo.Get(path.ID)
	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, categoriesResult)
}

// UpdateCategory updates dish category with new value provided in body
// @Summary Returns 204 if success and 4xx error if failed
// @Produce json
// @Accept json
// @Tags catering categories
// @Param id path string true "Catering ID"
// @Param categoryId path string true "Category ID"
// @Param body body request.UpdateCategory false "new category name"
// @Success 204 "Successfully updated"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/categories/{categoryId} [put]
func (dc category) Update(c *gin.Context) {
	var path types.PathCategory
	var body domain.Category

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	err, code := categoryRepo.Update(path, body)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}
