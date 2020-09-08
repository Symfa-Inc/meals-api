package interfaces

import uuid "github.com/satori/go.uuid"

// ImageDish struct for DB
type ImageDish struct {
	Base
	ImageID uuid.UUID `json:"imageId"`
	DishID  uuid.UUID `json:"dishId"`
} //@name ImageDishResponse
