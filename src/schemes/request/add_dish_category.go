package request

type AddDishCategory struct {
	Name string `json:"name" example:"закуски" binding:"required"`
} //@name AddDishCategoryResponse
