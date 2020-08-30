package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Aiscom-LLC/meals-api/types"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
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
		status := utils.DerefString(user.Status)

		if status == types.StatusTypesEnum.Deleted {
			utils.CreateError(http.StatusForbidden, errors.New("user was deleted"), c)
			c.Abort()
			return
		}

		for _, role := range roles {
			if user.Role == role {
				validRoles = append(validRoles, role)
			}
		}

		if len(validRoles) == 0 {
			utils.CreateError(http.StatusForbidden, errors.New("no permissions"), c)
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
