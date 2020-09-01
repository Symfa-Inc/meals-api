package api

import (
	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/schemes/request"
	"github.com/Aiscom-LLC/meals-api/schemes/response"
	"github.com/Aiscom-LLC/meals-api/services"
	"github.com/Aiscom-LLC/meals-api/types"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Order struct
type Order struct{}

// NewOrder return pointer to order struct
// with all methods
func NewOrder() *Order {
	return &Order{}
}

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
	var query types.DateQuery
	var order request.OrderRequest
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

	userOrder, code, err := services.NewOrderService().Add(query.Date, order, jwt.MapClaims(claims))

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusCreated, response.UserOrder{
		Items:   userOrder.Items,
		Status:  userOrder.Status,
		Total:   userOrder.Total,
		OrderID: userOrder.OrderID,
	})
}
