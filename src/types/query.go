package types

type PaginationQuery struct {
	Limit int `form:"limit"`
	Page  int `form:"page"`
}
