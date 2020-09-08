package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/domain"
)

// UserRepository is user interface for repository
type UserRepository interface {
	GetByKey(key, value string) (domain.UserClientCatering, error)
	Add(user domain.User) (domain.UserClientCatering, error)
	Get(companyID, companyType, userRole string, pagination url.PaginationQuery, filters url.UserFilterQuery) ([]domain.UserClientCatering, int, int, error)
	Delete(companyID, ctxUserRole string, user domain.User) (int, error)
	Update(companyID string, user domain.User) (domain.UserClientCatering, int, error)
}
