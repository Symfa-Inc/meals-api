package response

import uuid "github.com/satori/go.uuid"

// User struct for get users response
type User struct {
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

type UserResponse struct {
	ID          uuid.UUID    `gorm:"type:uuid;" json:"id"`
	FirstName   string       `gorm:"type:varchar(20)" json:"firstName"`
	LastName    string       `gorm:"type:varchar(20)" json:"lastName"`
	Email       string       `gorm:"type:varchar(30);unique;not null" json:"email"`
	Password    string       `gorm:"type:varchar(100);not null" json:"-"`
	CompanyType string       `sql:"type:company_types" gorm:"type:varchar(20);null" json:"companyType"`
	Catering    UserCatering `json:"catering"`
	Client      UserClient   `json:"client"`
	Role        string       `sql:"type:user_roles" json:"role"`
	Floor       int          `json:"floor"`
	Status      string       `sql:"type:status_types" gorm:"type:varchar(10);null" json:"status"`
}
