package models

// UpdateMeal request scheme
type UpdateMeal struct {
	Dishes []string `json:"dishes" binding:"required"`
	Status string   `json:"status"`
} // @name UpdateMealRequest
