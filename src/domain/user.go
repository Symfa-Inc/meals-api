package domain

import (
	"github.com/Aiscom-LLC/meals-api/src/types"
	uuid "github.com/satori/go.uuid"
)

// User model
type User struct {
	Base
	FirstName string  `gorm:"type:varchar(20)" json:"firstName"`
	LastName  string  `gorm:"type:varchar(20)" json:"lastName"`
	Email     string  `gorm:"type:varchar(30);not null" json:"email"`
	Password  string  `gorm:"type:varchar(100);not null" json:"-"`
	Role      string  `sql:"type:user_roles" json:"role"`
	Status    *string `sql:"type:status_types" gorm:"type:varchar(10);null" json:"status"`
}

// UserClientCatering struct for joined catering and clients table
type UserClientCatering struct {
	ID           uuid.UUID `gorm:"type:uuid;" json:"id"`
	FirstName    string    `gorm:"type:varchar(20)" json:"firstName"`
	LastName     string    `gorm:"type:varchar(20)" json:"lastName"`
	Email        string    `gorm:"type:varchar(30);unique;not null" json:"email"`
	Password     string    `gorm:"type:varchar(100);not null" json:"-"`
	CompanyType  *string   `sql:"type:company_types" gorm:"type:varchar(20);null" json:"companyType"`
	UserCatering `json:"catering"`
	UserClient   `json:"client"`
	Role         string  `sql:"type:user_roles" json:"role"`
	Floor        *int    `json:"floor"`
	Status       *string `sql:"type:status_types" gorm:"type:varchar(10);null" json:"status"`
} //@name UsersResponse

// UserCatering struct used in User
type UserCatering struct {
	CateringName *string    `json:"cateringName" gorm:"column:catering_name"`
	CateringID   *uuid.UUID `json:"cateringId" gorm:"column:catering_id"`
}

// UserClient struct used in User
type UserClient struct {
	ClientName *string    `json:"clientName" gorm:"column:client_name"`
	ClientID   *uuid.UUID `json:"clientId" gorm:"column:client_id"`
}

// UserRepository is user interface for repository
type UserRepository interface {
	GetByKey(key, value string) (UserClientCatering, error)
	Add(user User) (UserClientCatering, error)
	Get(companyID, companyType, userRole string, pagination types.PaginationQuery, filters types.UserFilterQuery) ([]UserClientCatering, int, int, error)
	Delete(companyID, ctxUserRole string, user User) (int, error)
	Update(companyID string, user User) (UserClientCatering, int, error)
}
