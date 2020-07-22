package request

// AddCatering request scheme
type AddCatering struct {
	Name string `json:"name,omitempty" example:"aisnovations" binding:"required"`
} //@name AddCateringRequest
