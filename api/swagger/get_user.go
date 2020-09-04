package swagger

import (
	uuid "github.com/satori/go.uuid"
)

// UserResponse struct
type UserResponse struct {
	ID          uuid.UUID `gorm:"api_types:uuid;" json:"id"`
	FirstName   string    `gorm:"api_types:varchar(20)" json:"firstName"`
	LastName    string    `gorm:"api_types:varchar(20)" json:"lastName"`
	Email       string    `gorm:"api_types:varchar(30);unique;not null" json:"email"`
	Password    string    `gorm:"api_types:varchar(100);not null" json:"-"`
	CompanyType string    `sql:"api_types:company_types" gorm:"api_types:varchar(20);null" json:"companyType"`
	ClientID    string    `json:"clientId"`
	CateringID  string    `json:"cateringId"`
	Role        string    `sql:"api_types:user_roles" json:"role"`
	Floor       int       `json:"floor"`
	Status      string    `sql:"api_types:status_types" gorm:"api_types:varchar(10);null" json:"status"`
}
