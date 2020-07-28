package types

// PaginationQuery struct used for pagination binding
type PaginationQuery struct {
	Limit int `form:"limit"`
	Page  int `form:"page"`
}

// PaginationWithDateQuery struct used for pagination binding
type PaginationWithDateQuery struct {
	Limit int    `form:"limit"`
	Page  int    `form:"page"`
	Date  string `form:"date" binding:"required"`
}

// DateQuery struct used for binding date
type DateQuery struct {
	Date string `form:"date" binding:"required"`
}

// DishIDQuery struct used for binding dish
type DishIDQuery struct {
	DishID string `form:"dishId" binding:"required"`
}

// CategoryIDQuery struct used for binding category id
type CategoryIDQuery struct {
	CategoryID string `form:"categoryID" binding:"required"`
}

// UserFilterQuery used to filter and sort users in DB
type UserFilterQuery struct {
	Query      string `form:"q"`
	Role       string `form:"role"`
	ClientName string `form:"client"`
	Status     string `form:"status"`
}
