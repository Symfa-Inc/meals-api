package swagger

import (
	"time"
)

// RefreshToken response model for /refresh-token response route
type RefreshToken struct {
	Code   int       `json:"code"`
	Expire time.Time `json:"expire"`
	Token  string    `json:"token"`
} //@name RefreshTokenResponse
