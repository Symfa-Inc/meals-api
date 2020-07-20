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

type client struct{}

func NewClient() *client {
	return &client{}
}

var clientRepo = repository.NewClientRepo()

// Add creates client
// @Summary Returns error or 201 status code if success
// @Produce json
// @Accept json
// @Tags client
// @Param body body request.AddName false "Client Name"
// @Success 200 {object} domain.Client false "client object"
// @Failure 400 {object} types.Error "Error"
// @Router /clients [post]
func (cl client) Add(c *gin.Context) {
	var body domain.Client
	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	client, err := clientRepo.Add(body)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusCreated, client)
}

// Get return list of clients
// @Summary Returns list of clients
// @Tags client
// @Produce json
// @Param limit query int false "used for pagination"
// @Param page query int false "used for pagination"
// @Success 200 {object} response.GetClients "List of clients"
// @Failure 400 {object} types.Error "Error"
// @Router /clients [get]
func (cl client) Get(c *gin.Context) {
	var query types.PaginationQuery

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	result, total, err := clientRepo.Get(query)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}

	c.JSON(http.StatusOK, response.GetClients{
		Items: result,
		Page:  query.Page,
		Total: total,
	})
}

// Delete soft delete of client
// @Summary Soft delete
// @Tags client
// @Produce json
// @Param id path string true "Client ID"
// @Success 204 "Successfully deleted"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /clients/{id} [delete]
func (cl client) Delete(c *gin.Context) {
	var path types.PathId
	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := clientRepo.Delete(path.ID); err != nil {
		utils.CreateError(http.StatusNotFound, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}

// Update updates client
// @Summary Returns 204 if success and 4xx error if failed
// @Produce json
// @Accept json
// @Tags client
// @Param id path string true "Client ID"
// @Param body body request.AddName false "Client Name"
// @Success 204 "Successfully updated"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /clients/{id} [put]
func (cl client) Update(c *gin.Context) {
	var path types.PathId
	var clientModal domain.Client

	if err := utils.RequestBinderBody(&clientModal, c); err != nil {
		return
	}

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err, code := clientRepo.Update(path.ID, clientModal); err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}
