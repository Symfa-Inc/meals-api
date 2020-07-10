package types

type PaginationQuery struct {
	Limit int `form:"limit"`
	Page  int `form:"page"`
}

type MealIdQuery struct {
	MealId string `form:"mealId" binding:"required"`
}

type CategoryIdQuery struct {
	CategoryId string `form:"categoryId" binding:"required"`
}
