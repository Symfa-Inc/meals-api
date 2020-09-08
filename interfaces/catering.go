package interfaces

import (
	"github.com/gin-gonic/gin"
)

// Catering model
type Catering struct {
	Base
	Name string `gorm:"type:varchar(30);not null" json:"name,omitempty" binding:"required"`
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
	Add(catering *Catering) error
	Update(id string, catering Catering) (int, error)
	Delete(id string) error
	GetByKey(key, value string) (Catering, error)
}
