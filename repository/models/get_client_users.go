package models

import (
	uuid "github.com/satori/go.uuid"
)

type GetClientUser struct {
	ID         uuid.UUID `json:"id"`
	UserClient `json:"client"`
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	Email      string  `json:"email"`
	Role       string  `json:"role"`
	Status     *string `json:"status"`
	Floor      int     `json:"floor"`
}

// Client response struct
type Client struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	AutoApproveOrders bool   `json:"autoApproveOrders"`
	UserCatering      `json:"catering"`
}
