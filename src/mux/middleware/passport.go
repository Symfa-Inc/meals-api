package middleware

import (
	"errors"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	users "go_api/src/models/user"
	requestAuth "go_api/src/schemes/request/auth"
	"go_api/src/utils"
	"time"
)

const IdentityKey = "id"

type UserID struct {
	ID string
}

//Middleware for user authentication
func Passport() *jwt.GinJWTMiddleware {
	authMiddleware, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "AIS Catering",
		Key:         []byte("jwtsecret"),
		Timeout:     time.Second * 30,
		MaxRefresh:  time.Hour * 2,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*UserID); ok {
				return jwt.MapClaims{
					IdentityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &UserID{
				ID: claims[IdentityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var body requestAuth.LoginUserRequest
			if err := c.ShouldBind(&body); err != nil {
				return "", errors.New("missing email or password")
			}

			user, err := users.GetUserByKey("email", body.Email)
			if err == nil {
				equal := utils.CheckPasswordHash(body.Password, user.Password)
				if equal {
					return &UserID{
						ID: user.ID.String(),
					}, nil
				}
			}
			return nil, errors.New("incorrect email or password")
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
	})
	return authMiddleware
}
