package auth


// refreshToken godoc
// @Summary Returns JSON with code, expire date and JWT token
// @Produce json
// @Accept json
// @Tags auth
// @Success 200 {object} auth.LoginResponse
// @Failure 404 {object} types.Error "Error"
// @Security ApiKeyAuth
// @Router /refresh_token [get]
func refreshToken() {}
