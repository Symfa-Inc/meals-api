package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/gin-gonic/gin"
)

// Catering model
type Catering struct {
	Base
	Name string `gorm:"varchar(30);not null" json:"name,omitempty" binding:"required"`
} //@name CateringsResponse

// CateringUsecase is catering interface for usecase
type CateringUsecase interface {
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
