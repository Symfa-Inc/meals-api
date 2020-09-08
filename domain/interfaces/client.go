package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// ClientAPI is client interface for API
type ClientAPI interface {
	Get(c *gin.Context)
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// ClientRepository is client interface for repository
type ClientRepository interface {
	Add(cateringID string, client *domain.Client) error
	GetCateringClientsOrders(cateringID string, query url.PaginationWithDateQuery) ([]models.ClientOrder, int, error)
	Get(query url.PaginationQuery, cateringID, role string) ([]models.Client, int, error)
	Delete(id string) error
	Update(id string, client domain.Client) (int, error)
	UpdateAutoApproveOrders(id string, status bool) (int, error)
	InitAutoApprove(id string) (int, error)
	UpdateAutoApproveSchedules(id string)
	GetByKey(key, value string) (domain.Client, error)
	GetAll() ([]domain.Client, error)
}

// ClientService is client interface for service
type ClientService interface {
	Get(query url.PaginationQuery, claims jwt.MapClaims) ([]models.Client, int, url.PaginationQuery, int, error)
}
