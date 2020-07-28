package usecase

import (
	"go_api/src/repository"
	"go_api/src/schemes/request"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Order struct
type Order struct{}

// NewOrder returns pointer to client struct
// with all methods
func NewOrder() *Order {
	return &Order{}
}

var orderRepo = repository.NewOrderRepo()

// Add creates order for client user
// @Summary Returns error or 201 status code if success
// @Produce json
// @Accept json
// @Tags users orders
// @Param id path string false "User ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Param body body request.OrderRequest false "User order"
// @Success 201 {object} response.UserOrder false "Order for user"
// @Failure 400 {object} types.Error "Error"
// @Router /users/{id}/orders [post]
func (o Order) Add(c *gin.Context) {
	var path types.PathID
	var order request.OrderRequest
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&order, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	for i, dish := range order.Items {
		if dish.Amount == 0 {
			utils.CreateError(http.StatusBadRequest, "can't add dish with 0 amount", c)
			return
		}
		for j := i + 1; j < len(order.Items); j++ {
			if dish.DishID == order.Items[j].DishID {
				utils.CreateError(http.StatusBadRequest, "can't add 2 same dishes, please increment amount field instead", c)
				return
			}
		}
	}

	date, err := time.Parse(time.RFC3339, query.Date)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	difference := date.Sub(time.Now().Truncate(time.Hour * 24)).Hours()

	if difference < 0 {
		utils.CreateError(http.StatusBadRequest, "can't add order to previous date", c)
		return
	}

	userOrder, err := orderRepo.Add(path.ID, date, order)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusCreated, userOrder)
}

// CancelOrder changes status of order to canceled
// @Summary Returns error or 204 status code if success
// @Produce json
// @Accept json
// @Tags users orders
// @Param id path string false "User ID"
// @Param orderId path string false "Order ID"
// @Success 204 "Successfully canceled"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /users/{id}/orders/{orderId} [delete]
func (o Order) CancelOrder(c *gin.Context) {
	var path types.PathOrder

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	code, err := orderRepo.CancelOrder(path.ID, path.OrderID)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
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
// @Success 200 {array} response.UserOrder false "Orders for user"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
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
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	userOrders, code, err := orderRepo.GetUserOrder(path.ID, query.Date)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, userOrders)
}

// GetClientOrders returns orders of provided client
// @Summary returns orders of provided client
// @Tags clients orders
// @Produce json
// @Param id path string true "Client ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 {object} response.SummaryOrdersResponse false "Orders for clients"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
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
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	client := types.CompanyTypesEnum.Client
	result, code, err := orderRepo.GetOrders("", path.ID, query.Date, client)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetCateringClientOrders returns orders of provided client
// @Summary returns orders of provided client
// @Tags caterings orders
// @Produce json
// @Param id path string true "Catering ID"
// @Param clientId path string true "Client ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 {object} response.SummaryOrdersResponse false "Orders for clients"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
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
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	client := types.CompanyTypesEnum.Catering
	result, code, err := orderRepo.GetOrders(path.ID, path.ClientID, query.Date, client)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, result)
}

// ApproveOrders changes status of orders for provided day
// to approved
// @Summary approve user orders
// @Tags clients orders
// @Produce json
// @Param id path string true "Client ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 "Successfully Approved orders"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
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
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
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
// @Success 200 {object} response.OrderStatus "order status"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
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
