package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/gin-gonic/gin"
)

// AddressUsecase is Address interface for usecase
type AddressUsecase interface {
	Get(c *gin.Context)
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// AddressRepository is Address interface for repository
type AddressRepository interface {
	Add(address domain.Address) (domain.Address, error)
	Get(id string) ([]domain.Address, int, error)
	Delete(path url.PathAddress) error
	Update(path url.PathAddress, address domain.Address) (domain.Address, error)
}
