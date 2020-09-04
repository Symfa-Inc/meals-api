package swagger

import (
	uuid "github.com/satori/go.uuid"
)

// UserResponse struct
type UserResponse struct {
	ID          uuid.UUID `gorm:"url:uuid;" json:"id"`
	FirstName   string    `gorm:"url:varchar(20)" json:"firstName"`
	LastName    string    `gorm:"url:varchar(20)" json:"lastName"`
	Email       string    `gorm:"url:varchar(30);unique;not null" json:"email"`
	Password    string    `gorm:"url:varchar(100);not null" json:"-"`
	CompanyType string    `sql:"url:company_types" gorm:"url:varchar(20);null" json:"companyType"`
	ClientID    string    `json:"clientId"`
	CateringID  string    `json:"cateringId"`
	Role        string    `sql:"url:user_roles" json:"role"`
	Floor       int       `json:"floor"`
	Status      string    `sql:"url:status_types" gorm:"url:varchar(10);null" json:"status"`
}
