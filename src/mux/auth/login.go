package auth

// Login godoc
// @Summary Returns JSON with code, expire date and JWT token (should be used for as refresh_token also)
// @Produce json
// @Accept json
// @Tags auth
// @Param body body auth.LoginUserRequest false "User Credentials"
// @Success 200 {object} auth.LoginResponse
// @Failure 404 {object} types.Error "Error"
// @Security ApiKeyAuth
// @Router /login [post]
func Login(){}

