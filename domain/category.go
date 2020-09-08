package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Category struct
type Category struct {
	Base
	Date       *time.Time `json:"date"`
	Name       string     `gorm:"type:varchar(30);not null" json:"name" binding:"required"`
	CateringID uuid.UUID  `json:"-"`
	ClientID   uuid.UUID  `json:"clientId"`
} //@name CategoryResponse

// CategoryAPI is category interface for API
type CategoryAPI interface {
	Get(c *gin.Context)
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// CategoryRepository is category interface for repository
type CategoryRepository interface {
	Add(category *Category) error
	Get(cateringID, clientID, date string) ([]Category, int, error)
	GetByKey(id, value, cateringID string) (Category, error)
	Delete(path url.PathCategory) (int, error)
	Update(path url.PathCategory, category *Category) (int, error)
}
