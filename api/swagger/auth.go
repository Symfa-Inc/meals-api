package swagger

// LoginUserRequest user model for /login request route
type LoginUserRequest struct {
	Email    string `json:"email" example:"meals@aisnovations.com" binding:"required"`
	Password string `json:"password" example:"Password12!" binding:"required"`
} //@name LoginRequest
