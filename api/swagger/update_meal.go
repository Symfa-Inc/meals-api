package swagger

// UpdateMeal request scheme
type UpdateMeal struct {
	Dishes []string `json:"dishes" binding:"required"`
} // @name UpdateMealRequest
