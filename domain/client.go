package domain

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Client model
type Client struct {
	Base
	Name              string    `gorm:"type:varchar(30);not null" json:"name,omitempty" binding:"required"`
	CateringID        uuid.UUID `json:"cateringId" swaggerignore:"true"`
	AutoApproveOrders bool      `json:"autoApproveOrders" swaggerignore:"true"`
} //@name ClientsResponse

// ClientAPI is client interface for API
type ClientAPI interface {
	Get(c *gin.Context)
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// ClientRepository is client interface for repository
type ClientRepository interface {
	Add(cateringID string, client *Client) error
	//GetCateringClientsOrders(cateringID string, query url.PaginationWithDateQuery) ([]models.ClientOrder, int, error)
	// TODO cycle Get(query url.PaginationQuery, cateringID, role string) ([]models.Client, int, error)
	Delete(id string) error
	Update(id string, client Client) (int, error)
	UpdateAutoApproveOrders(id string, status bool) (int, error)
	InitAutoApprove(id string) (int, error)
	UpdateAutoApproveSchedules(id string)
	GetByKey(key, value string) (Client, error)
	GetAll() ([]Client, error)
}

// ClientService is client interface for service
type ClientService interface {
	// TODO cycle Get(query url.PaginationQuery, claims jwt.MapClaims) ([]models.Client, int, url.PaginationQuery, int, error)
}
