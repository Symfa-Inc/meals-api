package api

import (
	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/api/swagger"
	types "github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/repository/enums"
	"github.com/Aiscom-LLC/meals-api/services"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Order struct
type Order struct{}

// NewOrder return pointer to order struct
// with all methods
func NewOrder() *Order {
	return &Order{}
}

var orderService = services.NewOrderService()
var orderRepo = repository.NewOrderRepo()

// Add creates order for client user
// @Summary Returns error or 201 status code if success
// @Produce json
// @Accept json
// @Tags users orders
// @Param id path string false "User ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Param body body swagger.OrderRequest false "User order"
// @Success 201 {object} swagger.UserOrder false "Order for user"
// @Failure 400 {object} Error "Error"
// @Router /users/{id}/orders [post]
func (o Order) Add(c *gin.Context) {
	var query types.DateQuery
	var order swagger.OrderRequest
	var path types.PathID

	claims, err := middleware.Passport().GetClaimsFromJWT(c)

	if err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&order, c); err != nil {
		return
	}

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	userOrder, code, err := orderService.Add(query.Date, order, jwt.MapClaims(claims))

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusCreated, swagger.UserOrder{
		Items:   userOrder.Items,
		Status:  userOrder.Status,
		Total:   userOrder.Total,
		OrderID: userOrder.OrderID,
	})
}

// CancelOrder changes status of order to canceled
// @Summary Returns error or 204 status code if success
// @Produce json
// @Accept json
// @Tags users orders
// @Param id path string false "User ID"
// @Param orderId path string false "Order ID"
// @Success 204 "Successfully canceled"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Error"
// @Router /users/{id}/orders/{orderId} [delete]
func (o Order) CancelOrder(c *gin.Context) {
	var path types.PathOrder

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	code, err := orderRepo.CancelOrder(path.ID, path.OrderID)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetUserOrder returns orders of provided user
// @Summary returns orders of provided user
// @Tags users orders
// @Produce json
// @Param id path string true "User ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 {array} swagger.UserOrder false "Orders for user"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /users/{id}/orders [get]
func (o Order) GetUserOrder(c *gin.Context) {
	var path types.PathID
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	_, err := time.Parse(time.RFC3339, query.Date)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	userOrder, code, err := orderRepo.GetUserOrder(path.ID, query.Date)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, userOrder)
}

// GetClientOrders returns orders of provided client
// @Summary returns orders of provided client
// @Tags clients orders
// @Produce json
// @Param id path string true "Client ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 {object} swagger.SummaryOrdersResponse false "Orders for clients"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /clients/{id}/orders [get]
func (o Order) GetClientOrders(c *gin.Context) {
	var path types.PathID
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	_, err := time.Parse(time.RFC3339, query.Date)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	company := enums.CompanyTypesEnum.Client

	result, code, err := orderRepo.GetOrders("", path.ID, query.Date, company)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(code, result)
}

// GetCateringClientOrders returns orders of provided client
// @Summary returns orders of provided client
// @Tags caterings orders
// @Produce json
// @Param id path string true "Catering ID"
// @Param clientId path string true "Client ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 {object} swagger.SummaryOrdersResponse false "Orders for clients"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /caterings/{id}/clients/{clientId}/orders [get]
func (o Order) GetCateringClientOrders(c *gin.Context) {
	var path types.PathClient
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	_, err := time.Parse(time.RFC3339, query.Date)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	company := enums.CompanyTypesEnum.Catering

	result, code, err := orderRepo.GetOrders(path.ID, path.ClientID, query.Date, company)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(code, result)
}

// ApproveOrders changes status of orders for provided day
// to approved
// @Summary approve user orders
// @Tags clients orders
// @Produce json
// @Param id path string true "Client ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 "Successfully Approved orders"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /clients/{id}/orders [put]
func (o Order) ApproveOrders(c *gin.Context) {
	var path types.PathID
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	if err := orderRepo.ApproveOrders(path.ID, query.Date); err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	c.Status(http.StatusOK)
}

// GetOrderStatus returns status of order
// @Summary returns status of order
// @Tags clients orders
// @Produce json
// @Param id path string true "Client ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 {object} swagger.OrderStatus "order status"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /clients/{id}/order-status [get]
func (o Order) GetOrderStatus(c *gin.Context) {
	var path types.PathID
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	status := orderRepo.GetOrdersStatus(path.ID, query.Date)

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}
