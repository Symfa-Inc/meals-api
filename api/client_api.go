package api

import (
	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/api/swagger"
	"github.com/Aiscom-LLC/meals-api/api/api_types"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/services"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
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
var clientService = services.NewClient()

// Add creates client
// @Summary Returns error or 201 status code if success
// @Produce json
// @Accept json
// @Tags caterings clients
// @Param id path string true "Catering ID"
// @Param body body swagger.AddName false "Client Name"
// @Success 201 {object} domain.Client false "client object"
// @Failure 400 {object} Error "Error"
// @Router /caterings/{id}/clients [post]
func (cl Client) Add(c *gin.Context) {
	var client domain.Client
	var path api_types.PathID
	if err := utils.RequestBinderBody(&client, c); err != nil {
		return
	}

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	client.CateringID, _ = uuid.FromString(path.ID)

	if err := clientRepo.Add(path.ID, &client); err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	c.JSON(http.StatusCreated, client)
}

// GetCateringClientsOrders return list of clients orders
// @Summary Returns list of clients orders
// @Tags caterings clients
// @Produce json
// @Param id path string true "Catering ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Param limit query int false "used for pagination"
// @Param page query int false "used for pagination"
// @Success 200 {object} swagger.GetCateringClientsSwagger "List of clients orders"
// @Failure 400 {object} Error "Error"
// @Router /caterings/{id}/clients-orders [get]
func (cl Client) GetCateringClientsOrders(c *gin.Context) {
	var query api_types.PaginationWithDateQuery
	var path api_types.PathID

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	_, err := time.Parse(time.RFC3339, query.Date)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	result, total, err := clientRepo.GetCateringClientsOrders(path.ID, query)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
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
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /clients/{id} [get]
func (cl Client) GetByID(c *gin.Context) {
	var path api_types.PathID

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	result, err := clientRepo.GetByKey("id", path.ID)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
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
// @Success 200 {object} swagger.GetClients "List of clients"
// @Failure 400 {object} Error "Error"
// @Router /clients [get]
func (cl Client) Get(c *gin.Context) {
	var query api_types.PaginationQuery

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	claims, err := middleware.Passport().GetClaimsFromJWT(c)

	if err != nil {
		utils.CreateError(http.StatusUnauthorized, err, c)
		return
	}

	result, total, query, code, err := clientService.Get(query, jwt.MapClaims(claims))

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": result,
		"page":  query.Page,
		"total": total,
	})
}

// GetByCateringID return list of clients
// @Summary Returns list of clients
// @Tags clients
// @Produce json
// @Param id path string false "Catering ID"
// @Param limit query int false "used for pagination"
// @Param page query int false "used for pagination"
// @Success 200 {object} swagger.GetClients "List of clients"
// @Failure 400 {object} Error "Error"
// @Router /caterings/{id}/clients [get]
func (cl Client) GetByCateringID(c *gin.Context) {
	var path api_types.PathID
	var query api_types.PaginationQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	result, total, err := clientRepo.Get(query, path.ID, "")

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
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
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /clients/{id} [delete]
func (cl Client) Delete(c *gin.Context) {
	var path api_types.PathID
	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := clientRepo.Delete(path.ID); err != nil {
		utils.CreateError(http.StatusNotFound, err, c)
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
// @Param body body swagger.UpdateClient false "update client model"
// @Success 204 "Successfully updated"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /clients/{id} [put]
func (cl Client) Update(c *gin.Context) {
	var path api_types.PathID
	var clientModal domain.Client

	if err := utils.RequestBinderBody(&clientModal, c); err != nil {
		return
	}

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if code, err := clientRepo.Update(path.ID, clientModal); err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateAutoApprove updates client
// @Summary Returns 204 if success and 4xx error if failed
// @Produce json
// @Accept json
// @Tags clients
// @Param id path string true "Client ID"
// @Param body body swagger.UpdateAutoApprove false "update auto approve"
// @Success 204 "Successfully updated"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /clients/{id}/auto-approve [put]
func (cl Client) UpdateAutoApprove(c *gin.Context) {
	var path api_types.PathID
	var body swagger.UpdateAutoApprove

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if code, err := clientRepo.UpdateAutoApproveOrders(path.ID, *body.Status); err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.Status(http.StatusNoContent)
}
