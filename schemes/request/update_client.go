package request

import uuid "github.com/satori/go.uuid"

// UpdateClient request scheme
type UpdateClient struct {
	Name       string    `json:"name" binding:"required" example:"aisnovations"`
	CateringID uuid.UUID `json:"cateringId"`
} // @name UpdateClientRequest
