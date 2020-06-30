package types

type PaginationQuery struct {
	Limit int `form:"limit"`
	Page  int `form:"page"`
}

type StartEndDateQuery struct {
	StartDate string `form:"startDate" binding:"required"`
	EndDate   string `form:"endDate" binding:"required"`
}

type CategoryIdQuery struct {
	CategoryId string `form:"categoryId" binding:"required"`
}
