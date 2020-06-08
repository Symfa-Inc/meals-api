package middleware

import (
	"errors"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"go_api/src/models"
	requestAuth "go_api/src/schemes/request/auth"
	"go_api/src/schemes/response/auth"
	"go_api/src/utils"
	"net/http"
	"os"
	"time"
)

const IdentityKeyID = "id"

type UserID struct {
	ID string
}

//Middleware for user authentication
func Passport() *jwt.GinJWTMiddleware {
	authMiddleware, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:          "AIS Catering",
		Key:            []byte(os.Getenv("JWTSECRET")),
		Timeout:        time.Hour * 4,
		MaxRefresh:     time.Hour * 24,
		IdentityKey:    IdentityKeyID,
		SendCookie:     true,
		CookieMaxAge:   time.Hour * 24,
		CookieHTTPOnly: true,
		CookieName:     "jwt",
		TokenLookup:    "cookie:jwt",
		LoginResponse: func(c *gin.Context, i int, s string, t time.Time) {
			value, _ := Passport().ParseTokenString(s)
			id := jwt.ExtractClaimsFromToken(value)["id"]
			user, _ := models.GetUserByKey("id", id.(string))
			c.JSON(http.StatusOK, auth.IsAuthenticated{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Role:      user.Role,
			})
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*UserID); ok {
				return jwt.MapClaims{
					IdentityKeyID: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &UserID{
				ID: claims[IdentityKeyID].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var body requestAuth.LoginUserRequest
			if err := c.ShouldBind(&body); err != nil {
				return "", errors.New("missing email or password")
			}

			user, err := models.GetUserByKey("email", body.Email)
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
