package middleware

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/Aiscom-LLC/meals-api/repository/models"

	"github.com/Aiscom-LLC/meals-api/repository/enums"

	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/utils"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// IdentityKeyID is used to tell
// by what field we will identify user
const IdentityKeyID = "id"

// UserID struct
type UserID struct {
	ID string
}

var userRepo = repository.NewUserRepo()

// Passport is middleware for user authentication
func Passport() *jwt.GinJWTMiddleware {
	authMiddleware, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:          "TastyOffice",
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
			result, err := userRepo.GetByID(id.(string))

			if err != nil {
				utils.CreateError(http.StatusUnauthorized, err, c)
				return
			}

			status := utils.DerefString(result.Status)

			if status == enums.StatusTypesEnum.Deleted {
				utils.CreateError(http.StatusForbidden, errors.New("user was deleted"), c)
				return
			}

			if status == enums.StatusTypesEnum.Invited {
				code, err := userRepo.UpdateStatus(result.ID, enums.StatusTypesEnum.Active)
				if err != nil {
					utils.CreateError(code, err, c)
					return
				}
			}

			c.JSON(http.StatusOK, result)
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
			var body models.LoginUserRequest
			if err := c.ShouldBind(&body); err != nil {
				return "", errors.New("missing email or password")
			}

			result, err := userRepo.GetAllByKey("email", body.Email)
			if err == nil {
				for i := range result {
					status := utils.DerefString(result[i].Status)
					if status == enums.StatusTypesEnum.Invited || status == enums.StatusTypesEnum.Active {
						equal := utils.CheckPasswordHash(body.Password, result[i].Password)
						if equal {
							return &UserID{
								ID: result[i].ID.String(),
							}, nil
						}
					}
				}
				return nil, errors.New("user was deleted")
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
