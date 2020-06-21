package models

import uuid "github.com/satori/go.uuid"

type Dish struct {
	Base
	Name           string    `json:"name" gorm:"not null" binding:"required"`
	Weight         int       `json:"weight" gorm:"not null"binding:"required"`
	Price          int       `json:"price" gorm:"not null" binding:"required"`
	Desc           string    `json:"desc"`
	Images         string    `json:"images"`
	CateringID     uuid.UUID `json:"-"`
	DishCategoryID uuid.UUID `json:"categoryId"`
} //@name DishResponse
