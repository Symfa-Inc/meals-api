package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/types"
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
	Add(address Address) (Address, error)
	Get(id string) ([]Address, int, error)
	Delete(path types.PathAddress) error
	Update(path types.PathAddress, address Address) (Address, error)
}
