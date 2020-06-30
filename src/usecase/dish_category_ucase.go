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

type dishCategory struct{}

func NewDishCategory() *dishCategory {
	return &dishCategory{}
}

var dishCategoryRepo = repository.NewDishCategoryRepo()

// AddDishCategory add dish category in DB
// returns 204 if success and 4xx if request failed
// @Summary Returns error if exists and 204 if success
// @Produce json
// @Accept json
// @Tags catering dish-categories
// @Param id path string true "Catering ID"
// @Param body body request.AddDishCategory false "Category Name"
// @Success 204 "Successfully created"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/dish-categories [post]
func (dc dishCategory) Add(c *gin.Context) {
	var body request.AddDishCategory
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
	categoryModel := domain.DishCategory{
		Base:       domain.Base{DeletedAt: body.DeletedAt},
		Name:       strings.ToLower(body.Name),
		CateringID: cateringId,
	}

	_, err := dishCategoryRepo.Add(categoryModel)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteDishCategory soft delete of category reading
// @Summary Soft delete
// @Tags catering dish-categories
// @Produce json
// @Param id path string true "Catering ID"
// @Param categoryId path string true "Category ID"
// @Success 204 "Successfully deleted"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/dish-categories/{categoryId} [delete]
func (dc dishCategory) Delete(c *gin.Context) {
	var path types.PathDishCategory
	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := dishCategoryRepo.Delete(path); err != nil {
		utils.CreateError(http.StatusNotFound, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetDishCategories returns list of categories or error
// @Summary Get list of categories
// @Tags catering dish-categories
// @Produce json
// @Param id path string false "Catering ID"
// @Success 200 {array} domain.DishCategory "array of category readings"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/dish-categories [get]
func (dc dishCategory) Get(c *gin.Context) {
	var path types.PathId

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	categoriesResult, err, code := dishCategoryRepo.Get(path.ID)
	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, categoriesResult)
}

// UpdateDishCategory updates dish category with new value provided in body
// @Summary Returns 204 if success and 4xx error if failed
// @Produce json
// @Accept json
// @Tags catering dish-categories
// @Param id path string true "Catering ID"
// @Param categoryId path string true "Category ID"
// @Param body body request.UpdateDishCategory false "new category name"
// @Success 204 "Successfully updated"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/dish-categories/{categoryId} [put]
func (dc dishCategory) Update(c *gin.Context) {
	var path types.PathDishCategory
	var body domain.DishCategory

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	err, code := dishCategoryRepo.Update(path, body)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}
