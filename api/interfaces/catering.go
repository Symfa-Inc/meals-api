package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// CateringAPI is catering interface for API
type CateringAPI interface {
	Get(c *gin.Context)
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// CateringRepository is catering interface for repository
type CateringRepository interface {
	Get(cateringID string, query url.PaginationQuery) ([]domain.Catering, int, error)
	Add(catering *domain.Catering) error
	Update(id string, catering domain.Catering) (int, error)
	Delete(id string) error
	GetByKey(key, value string) (domain.Catering, error)
}

// CateringService is catering interface for service
type CateringService interface {
	Get(claims jwt.MapClaims, query *url.PaginationQuery) ([]domain.Catering, int, int, error)
}
