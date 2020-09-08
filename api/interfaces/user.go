package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository/models"
)

// UserRepository is user interface for repository
type UserRepository interface {
	GetByKey(key, value string) (models.UserClientCatering, error)
	Add(user domain.User) (models.UserClientCatering, error)
	Get(companyID, companyType, userRole string, pagination url.PaginationQuery, filters url.UserFilterQuery) ([]models.UserClientCatering, int, int, error)
	Delete(companyID, ctxUserRole string, user domain.User) (int, error)
	Update(companyID string, user domain.User) (models.UserClientCatering, int, error)
}
