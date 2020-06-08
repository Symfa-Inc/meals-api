package catering

type AddCateringScheme struct {
	Name string `json:"name,omitempty" example:"aisnovations" binding:"required"`
}
