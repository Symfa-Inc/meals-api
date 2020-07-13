package types

type PaginationQuery struct {
	Limit int `form:"limit"`
	Page  int `form:"page"`
}

type DateQuery struct {
	Date string `form:"date" binding:"required"`
}

type CategoryIdQuery struct {
	CategoryId string `form:"categoryId" binding:"required"`
}
