package auth

// User model for /login request route
type LoginUserRequest struct {
	Email    string `json:"email" example:"admin@meals.com" binding:"required"`
	Password string `json:"password" example:"Password12!" binding:"required"`
} //@name LoginRequest
