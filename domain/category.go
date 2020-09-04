package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/api_types"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Category struct
type Category struct {
	Base
	Date       *time.Time `json:"date"`
	Name       string     `gorm:"api_types:varchar(30);not null" json:"name" binding:"required"`
	CateringID uuid.UUID  `json:"-"`
	ClientID   uuid.UUID  `json:"clientId"`
} //@name CategoryResponse

// CategoryUsecase is category interface for usecase
type CategoryUsecase interface {
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
	Delete(path api_types.PathCategory) (int, error)
	Update(path api_types.PathCategory, category *Category) (int, error)
}
