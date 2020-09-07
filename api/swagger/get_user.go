package swagger

import (
	uuid "github.com/satori/go.uuid"
)

// UserResponse struct
type UserResponse struct {
	ID          uuid.UUID `gorm:"uuid;" json:"id"`
	FirstName   string    `gorm:"varchar(20)" json:"firstName"`
	LastName    string    `gorm:"varchar(20)" json:"lastName"`
	Email       string    `gorm:"varchar(30);unique;not null" json:"email"`
	Password    string    `gorm:"varchar(100);not null" json:"-"`
	CompanyType string    `sql:"company_types" gorm:"varchar(20);null" json:"companyType"`
	ClientID    string    `json:"clientId"`
	CateringID  string    `json:"cateringId"`
	Role        string    `sql:"user_roles" json:"role"`
	Floor       int       `json:"floor"`
	Status      string    `sql:"status_types" gorm:"varchar(10);null" json:"status"`
}
