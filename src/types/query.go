package types

type PaginationQuery struct {
	Limit int `form:"limit"`
	Page  int `form:"page"`
}

type StartEndDateQuery struct {
	StartDate string `form:"startDate"`
	EndDate string `form:"endDate"`
}
