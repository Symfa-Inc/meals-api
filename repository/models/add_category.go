package models

import "time"

// AddCategory request scheme
type AddCategory struct {
	Name string     `json:"name" example:"закуски" binding:"required"`
	Date *time.Time `json:"date" example:"2020-06-29T00:00:00Z"`
} //@name AddCategoryRequest
