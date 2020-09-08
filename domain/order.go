package domain

import (
	"github.com/gin-gonic/gin"
	"time"
)

// Order struct for db
type Order struct {
	Base
	Total   *int
	Status  *string `sql:"type:order_status_types"`
	Comment *string
	Date    time.Time
}

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
	// TODO cycle Add(query string, order models.OrderRequest, claims jwt.MapClaims) (models.UserOrder, int, error)
}

// OrderRepository is order interface for repository
type OrderRepository interface {
	//TODO cycle Add(query string, order models.OrderRequest, claims jwt.MapClaims) (models.UserOrder, int, error)
}
