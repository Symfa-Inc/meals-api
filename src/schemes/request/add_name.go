package request

type AddName struct {
	Name string `json:"name,omitempty" example:"aisnovations" binding:"required"`
} //@name AddNameRequest
