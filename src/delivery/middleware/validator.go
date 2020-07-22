package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_api/src/utils"
	"net/http"
)

// ValidatorMiddleware used to validate users
// by their roles
type ValidatorMiddleware interface {
	ValidateRoles(roles ...string) gin.HandlerFunc
}

// Validator struct
type Validator struct{}

// NewValidator returns pointer to validator struct
// which includes all validate methods
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateRoles takes roles enums and validates each role
// for the upcoming request, aborts the request
// if role wasn't found in validRoles array
func (v *Validator) ValidateRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var validRoles []string
		claims, _ := Passport().GetClaimsFromJWT(c)

		userID := claims["id"]
		user, _ := userRepo.GetByKey("id", fmt.Sprintf("%v", userID))

		for _, role := range roles {
			if user.Role == role {
				validRoles = append(validRoles, role)
			}
		}

		if len(validRoles) == 0 {
			utils.CreateError(http.StatusForbidden, "no permissions", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
