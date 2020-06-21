package request

import uuid "github.com/satori/go.uuid"

type AddDish struct {
	Name       string    `json:"name" gorm:"not null" binding:"required" example:"грибной суп"`
	Weight     int       `json:"weight" gorm:"not null"binding:"required" example:"250"`
	Price      int       `json:"price" gorm:"not null" binding:"required" example:"120"`
	Desc       string    `json:"desc" example:"Очень вкусный"`
	Images     string    `json:"images"`
	CategoryID uuid.UUID `json:"categoryId" binding:"required"`
} // @name AddDishRequest
