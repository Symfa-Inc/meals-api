package domain

import (
	uuid "github.com/satori/go.uuid"
)

// User model
type User struct {
	Base
	FirstName   string     `gorm:"type:varchar(20)" json:"firstName"`
	LastName    string     `gorm:"type:varchar(20)" json:"lastName"`
	Email       string     `gorm:"type:varchar(30);unique;not null" json:"email"`
	Password    string     `gorm:"type:varchar(100);not null" json:"-"`
	CompanyType *string    `sql:"type:company_types" gorm:"type:varchar(20);null" json:"companyType"`
	CateringID  *uuid.UUID `json:"cateringId"`
	ClientID    *uuid.UUID `json:"clientId"`
	Role        string     `sql:"type:user_roles" json:"role"`
	Floor       *int       `json:"floor"`
	Status      *string    `sql:"type:status_types" gorm:"type:varchar(10);null" json:"status"`
} //@name UsersResponse

// UserRepository is user interface for repository
type UserRepository interface {
	GetByKey(key, value string) (User, error)
	Add(companyID string, user User) (User, error)
}
