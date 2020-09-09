package domain

import (
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// OrderAPI is order interface for API
type OrderAPI interface {
	Add(c *gin.Context)
	CancelOrder(c *gin.Context)
	GetUserOrder(c *gin.Context)
	GetClientOrders(c *gin.Context)
	GetCateringClientOrders(c *gin.Context)
	ApproveOrders(c *gin.Context)
	GetOrderStatus(c *gin.Context)
}

// OrderService is order interface for service
type OrderService interface {
	Add(query string, order models.OrderRequest, claims jwt.MapClaims) (models.UserOrder, int, error)
}

// OrderRepository is order interface for repository
type OrderRepository interface {
	Add(query string, order models.OrderRequest, claims jwt.MapClaims) (models.UserOrder, int, error)
}
