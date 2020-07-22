package domain

import (
	"github.com/gin-gonic/gin"
)

// Client model
type Client struct {
	Base
	Name string `gorm:"type:varchar(30);not null" json:"name,omitempty" binding:"required"`
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
	Add(client Client) (Client, error)
	Update(id string, client Client) (int, error)
	Delete(id string) error
	GetByKey(key, value string) (Client, error)
}
