package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_api/src/utils"
	"net/http"
)

type ValidatorMiddleware interface {
	ValidateRoles(roles ...string) gin.HandlerFunc
}

type validator struct{}

func NewValidator() *validator {
	return &validator{}
}

func (v *validator) ValidateRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var validRoles []string
		claims, _ := Passport().GetClaimsFromJWT(c)

		userId := claims["id"]
		user, _ := userRepo.GetByKey("id", fmt.Sprintf("%v", userId))

		for _, role := range roles {
			if user.Role == role {
				validRoles = append(validRoles, role)
			}
		}

		if len(validRoles) <= 0 {
			utils.CreateError(http.StatusForbidden, "no permissions", c)
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
