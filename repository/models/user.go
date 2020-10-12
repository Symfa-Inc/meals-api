package models

import uuid "github.com/satori/go.uuid"

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
	CateringName *string `json:"cateringName" gorm:"type:column:catering_name"`
	CateringID   *string `json:"cateringId" gorm:"column:catering_id"`
}

type UserClient struct {
	ClientName *string `json:"clientName" gorm:"column:client_name"`
	ClientID   *string `json:"clientId" gorm:"column:client_id"`
}
