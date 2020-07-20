package domain

import uuid "github.com/satori/go.uuid"

type ImageDish struct {
	Base
	ImageID uuid.UUID `json:"imageId"`
	DishID  uuid.UUID `json:"dishId"`
} //@name ImageDishResponse
