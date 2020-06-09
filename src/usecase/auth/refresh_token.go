package auth

// @Summary Return JSON with code, expire date and new JWT
// @Produce json
// @Accept json
// @Tags auth
// @Success 200 {object} auth.RefreshToken
// @Failure 401 {object} types.Error "Error"
// @Router /refresh-token [get]
func refreshToken() {}
