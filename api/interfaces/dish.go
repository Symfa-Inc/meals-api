package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/gin-gonic/gin"
)

// DishAPI is dish interface for API
type DishAPI interface {
	Add(c *gin.Context)
	Delete(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
}

// DishRepository is dish interface for repository
type DishRepository interface {
	Add(cateringID string, dish *domain.Dish) error
	Delete(path url.PathDish) error
	Get(cateringID, categoryID string) ([]domain.Dish, int, error)
	FindByID(cateringID, id string) (domain.Dish, int, error)
	GetByKey(key, value, cateringID, categoryID string) (domain.Dish, int, error)
	Update(path url.PathDish, dish domain.Dish) (int, error)
}
