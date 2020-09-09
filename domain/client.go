package domain

import (
	uuid "github.com/satori/go.uuid"
)

// Client model
type Client struct {
	Base
	Name              string    `gorm:"type:varchar(30);not null" json:"name,omitempty" binding:"required"`
	CateringID        uuid.UUID `json:"cateringId" swaggerignore:"true"`
	AutoApproveOrders bool      `json:"autoApproveOrders" swaggerignore:"true"`
} //@name ClientsResponse
