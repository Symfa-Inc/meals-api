package domain

import (
	"go_api/src/types"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Address struct
type Address struct {
	Base
	City     string    `json:"city" gorm:"not null" binding:"required"`
	Street   string    `json:"street" gorm:"not null" binding:"required"`
	House    string    `json:"house" gorm:"not null" binding:"required"`
	Floor    int       `json:"floor" gorm:"not null" binding:"required"`
	ClientID uuid.UUID `json:"-"`
} //@name AddressResponse

// AddressUsecase is Address interface for usecase
type AddressUsecase interface {
	Get(c *gin.Context)
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// AddressRepository is Address interface for repository
type AddressRepository interface {
	Add(Address Address) (Address, error)
	Get(id string) ([]Address, int, error)
	Delete(path types.PathID) error
	Update(path types.PathID, Address Address) (int, error)
}
