package request

type AddCateringRequest struct {
	Name string `json:"name,omitempty" example:"aisnovations" binding:"required"`
} //@name AddCateringResponse
