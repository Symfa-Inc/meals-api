package domain

import (
	uuid "github.com/satori/go.uuid"
	"go_api/src/types"
)

// User model
type User struct {
	Base
	FirstName   string              `gorm:"type:varchar(20); unique_index" json:"firstName,omitempty"`
	LastName    string              `gorm:"type:varchar(20)" json:"lastName,omitempty"`
	Email       string              `gorm:"type:varchar(30);unique;not null" json:"email,omitempty"`
	Password    string              `gorm:"type:varchar(100);not null" json:"-"`
	CompanyType *types.CompanyTypes `gorm:"type:varchar(20);null" json:"companyType,omitempty"`
	CateringId  *uuid.UUID          `json:"cateringId,omitempty"`
	ClientId    *uuid.UUID          `json:"clientId,omitempty"`
	Role        string              `sql:"type:user_roles" json:"role,omitempty"`
} //@name UsersResponse

type UserRepository interface {
	GetByKey(key, value string) (User, error)
}
