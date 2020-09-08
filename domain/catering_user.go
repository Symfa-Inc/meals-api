package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// CateringUser struct
type CateringUser struct {
	Base
	CateringID uuid.UUID `json:"cateringId"`
	UserID     uuid.UUID `json:"userId"`
}

// CateringUserAPI is CateringUser interface for API
type CateringUserAPI interface {
	Add(c *gin.Context)
	Get(c *gin.Context)
	Delete(c *gin.Context)
	Update(c *gin.Context)
}

// CateringUserRepository is CateringUser interface for repository
type CateringUserRepository interface {
	GetByKey(key, value string) (CateringUser, error)
	Add(cateringUser CateringUser) error
	//Get(cateringID string, pagination url.PaginationQuery, filters url.UserFilterQuery) ([]models.GetCateringUser, int, int, error)
	Delete(cateringID, ctxUserRole string, user User) (int, error)
	Update(user *User) (int, error)
}

// CateringUserService is CateringUser interface for service
type CateringUserService interface {
	Add(path url.PathID, user User) (UserClientCatering, User, string, error, error)
}
