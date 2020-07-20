package domain

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"go_api/src/types"
	"time"
)

type CategoryBase struct {
	ID        uuid.UUID `gorm:"type:uuid;" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (base *CategoryBase) BeforeCreate(scope *gorm.Scope) error {
	uuidv4, err := uuid.NewV4()
	if err != nil {
		return err
	}
	return scope.SetColumn("ID", uuidv4)
}

type Category struct {
	CategoryBase
	DeletedAt  *time.Time `sql:"index" json:"deletedAt"`
	Name       string     `gorm:"type:varchar(30);not null" json:"name" binding:"required"`
	CateringID uuid.UUID  `json:"-"`
} //@name CategoryResponse

type CategoryUsecase interface {
	Get(c *gin.Context)
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type CategoryRepository interface {
	Add(category Category) (Category, error)
	Get(id string) ([]Category, error, int)
	GetByKey(id, value, cateringId string) (Category, error)
	Delete(path types.PathCategory) error
	Update(path types.PathCategory, category Category) (error, int)
}
