package usecase

import (
	"github.com/gin-gonic/gin"
	"go_api/src/domain"
	"go_api/src/repository"
	"go_api/src/schemes/response"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
)

type catering struct{}

func NewCatering() *catering {
	return &catering{}
}

var cateringRepo = repository.NewCateringRepo()

// AddCatering creates catering
// @Summary Returns error or 204 status code if success
// @Produce json
// @Accept json
// @Tags catering
// @Param body body request.AddCatering false "Catering Name"
// @Success 204 "Successfully created"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings [post]
func (ca catering) Add(c *gin.Context) {
	var cateringModel domain.Catering
	if err := utils.RequestBinderBody(&cateringModel, c); err != nil {
		return
	}

	err := cateringRepo.Add(cateringModel)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteCatering soft delete of catering
// @Summary Soft delete
// @Tags catering
// @Produce json
// @Param id path string true "Catering ID"
// @Success 204 "Successfully deleted"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id} [delete]
func (ca catering) Delete(c *gin.Context) {
	var path types.PathId
	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := cateringRepo.Delete(path.ID); err != nil {
		utils.CreateError(http.StatusNotFound, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetCatering returns catering
// @Summary Returns info about catering
// @Tags catering
// @Produce json
// @Param id path string true "Catering ID"
// @Success 200 {object} domain.Catering "catering model"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id} [get]
func (ca catering) GetById(c *gin.Context) {
	var path types.PathId

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	result, err := cateringRepo.GetByKey("id", path.ID)

	if err != nil {
		if err.Error() == "record not found" {
			utils.CreateError(http.StatusNotFound, err.Error(), c)
			return
		} else {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

// GetCaterings return list of caterings
// @Summary Returns list of caterings
// @Tags catering
// @Produce json
// @Param limit query int false "used for pagination"
// @Param page query int false "used for pagination"
// @Success 200 {object} response.GetCaterings "List of caterings"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings [get]
func (ca catering) Get(c *gin.Context) {
	var query types.PaginationQuery

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	result, total, err := cateringRepo.Get(query)

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

// UpdateCatering updates catering
// @Summary Returns 204 if success and 4xx error if failed
// @Produce json
// @Accept json
// @Tags catering
// @Param id path string true "Catering ID"
// @Param body body request.AddCatering false "Catering Name"
// @Success 204 "Successfully updated"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id} [put]
func (ca catering) Update(c *gin.Context) {
	var path types.PathId
	var cateringModel domain.Catering

	if err := utils.RequestBinderBody(&cateringModel, c); err != nil {
		return
	}

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err, code := cateringRepo.Update(path.ID, cateringModel); err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}
