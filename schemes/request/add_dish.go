package request

import uuid "github.com/satori/go.uuid"

// AddDish request scheme
type AddDish struct {
	Name       string    `json:"name" gorm:"not null" binding:"required" example:"грибной суп"`
	Weight     float32   `json:"weight" gorm:"not null" binding:"required" example:"250"`
	Price      float32   `json:"price" gorm:"not null" binding:"required" example:"120"`
	Desc       string    `json:"desc" example:"Очень вкусный"`
	CategoryID uuid.UUID `json:"categoryId" binding:"required"`
} // @name AddDishRequest
