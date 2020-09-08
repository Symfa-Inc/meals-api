package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/gin-gonic/gin"
)

// CategoryAPI is category interface for API
type CategoryAPI interface {
	Get(c *gin.Context)
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// CategoryRepository is category interface for repository
type CategoryRepository interface {
	Add(category *domain.Category) error
	Get(cateringID, clientID, date string) ([]domain.Category, int, error)
	GetByKey(id, value, cateringID string) (domain.Category, error)
	Delete(path url.PathCategory) (int, error)
	Update(path url.PathCategory, category *domain.Category) (int, error)
}
