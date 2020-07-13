package request

type UpdateMeal struct {
	Dishes []string `json:"dishes" binding:"required"`
} // @name UpdateMealRequest
