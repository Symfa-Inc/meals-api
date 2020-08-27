package request

// UpdateAutoApprove request scheme
type UpdateAutoApprove struct {
	Status *bool `json:"status" binding:"required" example:"true"`
} // @name UpdateAutoApprove
