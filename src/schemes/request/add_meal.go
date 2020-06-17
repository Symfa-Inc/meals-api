package request

import "time"

type AddMealRequestList []AddMealRequest // @name AddMealListResponse

type AddMealRequest struct {
	Date time.Time `json:"date" binding:"required" example:"2020-06-20T00:00:00Z"`
} // @name AddMealResponse
