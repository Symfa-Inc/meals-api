package domain

import (
	"github.com/Aiscom-LLC/meals-api/src/types"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

// CategoryBase struct without deletedAt
type CategoryBase struct {
	ID        uuid.UUID `gorm:"type:uuid;" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// BeforeCreate func which generates uuid v4 for each inserted row
func (base *CategoryBase) BeforeCreate(scope *gorm.Scope) error {
	uuidv4 := uuid.NewV4()
	return scope.SetColumn("ID", uuidv4)
}

// Category struct
type Category struct {
	CategoryBase
	DeletedAt  *time.Time `sql:"index" json:"deletedAt"`
	Name       string     `gorm:"type:varchar(30);not null" json:"name" binding:"required"`
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
	Delete(path types.PathCategory) (int, error)
	Update(path types.PathCategory, category *Category) (int, error)
}
