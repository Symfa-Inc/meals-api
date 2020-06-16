package meal

import "time"

type UpdateMealRequest struct {
	Date       time.Time `json:"date" binding:"required" example:"2020-06-25T00:00:00Z"`
	CateringID string    `json:"cateringId" binding:"required"`
} // @name UpdateMealResponse
