package swagger

// UpdateCategory request scheme
type UpdateCategory struct {
	Name string `json:"name" binding:"required" example:"веган"`
} // @name UpdateCategoryRequest
