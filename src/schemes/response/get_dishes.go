package response

import uuid "github.com/satori/go.uuid"

type GetDish struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name" gorm:"not null" binding:"required"`
	Weight     int       `json:"weight" gorm:"not null"binding:"required"`
	Price      int       `json:"price" gorm:"not null" binding:"required"`
	Desc       string    `json:"desc"`
	Images     []string  `json:"images,omitempty"`
	CateringID uuid.UUID `json:"-"`
	CategoryID uuid.UUID `json:"categoryId,omitempty"`
}
