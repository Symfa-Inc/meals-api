package request

type UpdateDishCategory struct {
	Name string `json:"name" binding:"required" example:"веган"`
} // @name UpdateDishCategoryRequest
