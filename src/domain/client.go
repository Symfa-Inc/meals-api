package domain

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Client model
type Client struct {
	Base
	Name       string    `gorm:"type:varchar(30);not null" json:"name,omitempty" binding:"required"`
	CateringID uuid.UUID `json:"-" swaggerignore:"true"`
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
	Add(cateringID string, client Client) (Client, error)
	Update(cateringID, id string, client Client) (int, error)
	Delete(cateringID, id string) error
	GetByKey(key, value string) (Client, error)
}
