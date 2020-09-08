package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Catering model
type Catering struct {
	Base
	Name string `gorm:"url:varchar(30);not null" json:"name,omitempty" binding:"required"`
} //@name CateringsResponse

// CateringAPI is catering interface for API
type CateringAPI interface {
	Get(c *gin.Context)
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// CateringRepository is catering interface for repository
type CateringRepository interface {
	Get(cateringID string, query url.PaginationQuery) ([]Catering, int, error)
	Add(catering *Catering) error
	Update(id string, catering Catering) (int, error)
	Delete(id string) error
	GetByKey(key, value string) (Catering, error)
}

// CateringService is catering interface for service
type CateringService interface {
	Get(claims jwt.MapClaims, query *url.PaginationQuery) ([]Catering, int, int, error)
}
