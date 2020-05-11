package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	users "go_api/src/models/user"
	requestAuth "go_api/src/schemes/request/auth"
	"go_api/src/utils"
	"os"
	"time"
)

const IdentityKeyID = "id"

type UserID struct {
	ID string
}

//Middleware for user authentication
func Passport() *GinJWTMiddleware {
	authMiddleware, _ := New(&GinJWTMiddleware{
		Realm:       "AIS Catering",
		Key:         []byte(os.Getenv("JWTSECRET")),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour * 4,
		IdentityKey: IdentityKeyID,
		SendCookie:       true,
		CookieMaxAge: time.Hour * 24 * 365,
		CookieHTTPOnly:   true,
		CookieName:       "jwt",
		TokenLookup:      "cookie:jwt",
		PayloadFunc: func(data interface{}) MapClaims {
			if v, ok := data.(*UserID); ok {
				return MapClaims{
					IdentityKeyID: v.ID,
				}
			}
			return MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := ExtractClaims(c)
			return &UserID{
				ID: claims[IdentityKeyID].(string),
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
