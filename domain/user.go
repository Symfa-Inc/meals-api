package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	uuid "github.com/satori/go.uuid"
)

// User model
type User struct {
	Base
	FirstName   string  `gorm:"varchar(20)" json:"firstName"`
	LastName    string  `gorm:"varchar(20)" json:"lastName"`
	Email       string  `gorm:"varchar(30);not null" json:"email"`
	Password    string  `gorm:"varchar(100);not null" json:"-"`
	Role        string  `sql:"user_roles" json:"role"`
	CompanyType *string `sql:"company_types" gorm:"varchar(20);null" json:"companyType"`
	Status      *string `sql:"status_types" gorm:"varchar(10);null" json:"status"`
}

// UserClientCatering struct for joined catering and clients table
type UserClientCatering struct {
	ID           uuid.UUID `json:"id"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	CompanyType  *string   `json:"companyType"`
	UserCatering `json:"catering"`
	UserClient   `json:"client"`
	Role         string  `json:"role"`
	Floor        *int    `json:"floor"`
	Status       *string `json:"status"`
} //@name UsersResponse

type UserCatering struct {
	CateringName *string `json:"cateringName" gorm:"column:catering_name"`
	CateringID   *string `json:"cateringId" gorm:"column:catering_id"`
}

type UserClient struct {
	ClientName *string `json:"clientName" gorm:"column:client_name"`
	ClientID   *string `json:"clientId" gorm:"column:client_id"`
}

// UserRepository is user interface for repository
type UserRepository interface {
	GetByKey(key, value string) (UserClientCatering, error)
	Add(user User) (UserClientCatering, error)
	Get(companyID, companyType, userRole string, pagination url.PaginationQuery, filters url.UserFilterQuery) ([]UserClientCatering, int, int, error)
	Delete(companyID, ctxUserRole string, user User) (int, error)
	Update(companyID string, user User) (UserClientCatering, int, error)
}
