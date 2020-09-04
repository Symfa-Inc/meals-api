package domain

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Client model
type Client struct {
	Base
	Name              string    `gorm:"api_types:varchar(30);not null" json:"name,omitempty" binding:"required"`
	CateringID        uuid.UUID `json:"cateringId" swaggerignore:"true"`
	AutoApproveOrders bool      `json:"autoApproveOrders" swaggerignore:"true"`
} //@name ClientsResponse

// ClientUsecase is client interface for usecase
type ClientUsecase interface {
	Get(c *gin.Context)
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// ClientRepository is client interface for repository
type ClientRepository interface {
	Add(cateringID string, client *Client, user UserClientCatering) error
	Update(id string, client Client) (int, error)
	/* TODO fix cycle imports
	GetCateringClientsOrders(cateringID string, query api_types.PaginationWithDateQuery) ([]response.ClientOrder, int, error)
	Get(query api_types.PaginationQuery, cateringID, role string) ([]response.Client, int, error)
	*/
	UpdateAutoApproveOrders(id string, status bool) (int, error)
	InitAutoApprove(id string) (int, error)
	Delete(cateringID, id string) error
	GetByKey(key, value string) (Client, error)
	GetAll() ([]Client, error)
}
