package request

type UpdateDishCategoryRequest struct {
	Name       string `json:"name" binding:"required" example:"веган"`
	CateringID string `json:"cateringId" binding:"required"`
} // @name UpdateDishCategoryResponse
