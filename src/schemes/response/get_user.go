package response

import (
	uuid "github.com/satori/go.uuid"
)

// UserResponse struct
type UserResponse struct {
	ID          uuid.UUID `gorm:"type:uuid;" json:"id"`
	FirstName   string    `gorm:"type:varchar(20)" json:"firstName"`
	LastName    string    `gorm:"type:varchar(20)" json:"lastName"`
	Email       string    `gorm:"type:varchar(30);unique;not null" json:"email"`
	Password    string    `gorm:"type:varchar(100);not null" json:"-"`
	CompanyType string    `sql:"type:company_types" gorm:"type:varchar(20);null" json:"companyType"`
	ClientID    string    `json:"clientId"`
	CateringID  string    `json:"cateringId"`
	Role        string    `sql:"type:user_roles" json:"role"`
	Floor       int       `json:"floor"`
	Status      string    `sql:"type:status_types" gorm:"type:varchar(10);null" json:"status"`
}
