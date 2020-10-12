package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"github.com/gin-gonic/gin"
)

// CateringUserAPI is CateringUser interface for API
type CateringUserAPI interface {
	Add(c *gin.Context)
	Get(c *gin.Context)
	Delete(c *gin.Context)
	Update(c *gin.Context)
}

// CateringUserRepository is CateringUser interface for repository
type CateringUserRepository interface {
	GetByKey(key, value string) (domain.CateringUser, error)
	Add(cateringUser domain.CateringUser) error
	Get(cateringID string, pagination url.PaginationQuery, filters url.UserFilterQuery) ([]models.GetCateringUser, int, int, error)
	Delete(cateringID, ctxUserRole string, user domain.User) (int, error)
	Update(user *domain.User) (int, error)
}

// CateringUserService is CateringUser interface for service
type CateringUserService interface {
	Add(path url.PathID, user domain.User) (models.UserClientCatering, domain.User, string, error, error)
}
