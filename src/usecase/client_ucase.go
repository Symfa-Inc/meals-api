package usecase

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go_api/src/delivery/middleware"
	"go_api/src/domain"
	"go_api/src/repository"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
	"time"
)

// Client struct
type Client struct{}

// NewClient returns pointer to client struct
// with all methods
func NewClient() *Client {
	return &Client{}
}

var clientRepo = repository.NewClientRepo()

// Add creates client
// @Summary Returns error or 201 status code if success
// @Produce json
// @Accept json
// @Tags caterings clients
// @Param id path string true "Catering ID"
// @Param body body request.AddName false "Client Name"
// @Success 201 {object} domain.Client false "client object"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/clients [post]
func (cl Client) Add(c *gin.Context) {
	var client domain.Client
	var path types.PathID
	if err := utils.RequestBinderBody(&client, c); err != nil {
		return
	}

	if err := utils.RequestBinderURI(&path, c); err != nil {
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

	parsedCateringID, _ := uuid.FromString(path.ID)
	client.CateringID = parsedCateringID

	if err := clientRepo.Add(path.ID, &client, user); err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusCreated, client)
}

// GetCateringClients return list of clients
// @Summary Returns list of clients
// @Tags caterings clients
// @Produce json
// @Param id path string true "Catering ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Param limit query int false "used for pagination"
// @Param page query int false "used for pagination"
// @Success 200 {object} response.GetCateringClientsSwagger "List of clients"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/clients [get]
func (cl Client) GetCateringClients(c *gin.Context) {
	var query types.PaginationWithDateQuery
	var path types.PathID

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	_, err := time.Parse(time.RFC3339, query.Date)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	result, total, err := clientRepo.GetCateringClients(path.ID, query)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}

	c.JSON(http.StatusOK, gin.H{
		"items": result,
		"page":  query.Page,
		"total": total,
	})
}

// GetByID returns client
// @Summary Returns info about client
// @Tags clients
// @Produce json
// @Param id path string true "Client ID"
// @Success 200 {object} domain.Client "client model"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /clients/{id} [get]
func (cl Client) GetByID(c *gin.Context) {
	var path types.PathID

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	result, err := clientRepo.GetByKey("id", path.ID)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, result)
}

// Get return list of clients
// @Summary Returns list of clients
// @Tags clients
// @Produce json
// @Param limit query int false "used for pagination"
// @Param page query int false "used for pagination"
// @Success 200 {object} response.GetClients "List of clients"
// @Failure 400 {object} types.Error "Error"
// @Router /clients [get]
func (cl Client) Get(c *gin.Context) {
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

	result, total, err := clientRepo.Get(query, cateringID, user.Role)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}

	c.JSON(http.StatusOK, gin.H{
		"items": result,
		"page":  query.Page,
		"total": total,
	})
}

// Delete soft delete of client
// @Summary Soft delete
// @Tags clients
// @Produce json
// @Param id path string true "Client ID"
// @Success 204 "Successfully deleted"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /clients/{id} [delete]
func (cl Client) Delete(c *gin.Context) {
	var path types.PathID
	if err := utils.RequestBinderURI(&path, c); err != nil {
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
// @Tags clients
// @Param id path string true "Client ID"
// @Param body body request.UpdateClient false "update client model"
// @Success 204 "Successfully updated"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /clients/{id} [put]
func (cl Client) Update(c *gin.Context) {
	var path types.PathID
	var clientModal domain.Client

	if err := utils.RequestBinderBody(&clientModal, c); err != nil {
		return
	}

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if code, err := clientRepo.Update(path.ID, clientModal); err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}
