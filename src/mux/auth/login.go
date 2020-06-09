package auth

// @Summary Returns a JSON with code, expire date and JWT (should be used as refresh token also)
// @Produce json
// @Accept json
// @Tags auth
// @Param body body auth.LoginUserRequest false "User Credentials"
// @Success 200 {object} models.User
// @Failure 401 {object} types.Error "Error"
// @Router /login [post]
func login(){}