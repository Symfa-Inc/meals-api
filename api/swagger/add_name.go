package swagger

// AddName request scheme
type AddName struct {
	Name string `json:"name,omitempty" example:"aisnovations" binding:"required"`
} //@name AddNameRequest
