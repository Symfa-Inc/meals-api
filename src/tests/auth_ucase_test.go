package tests

import (
	"github.com/Aiscom-LLC/meals-api/src/delivery/api"
	"github.com/Aiscom-LLC/meals-api/src/delivery/middleware"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIsAuthenticated(t *testing.T) {
	r := gofight.New()

	userRepo := repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})

	// Trying to login without jwt cookie
	r.GET("/is-authenticated").
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusUnauthorized, r.Code)
			assert.Equal(t, "cookie token is empty", errorValue)
		})

	// Trying to login with jwt cookie
	r.GET("/is-authenticated").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		data := r.Body.Bytes()
		email, _ := jsonparser.GetString(data, "email")
		name, _ := jsonparser.GetString(data, "firstName")
		assert.Equal(t, http.StatusOK, r.Code)
		assert.Equal(t, "admin@meals.com", email)
		assert.Equal(t, "super", name)
	})
}

func TestValidator(t *testing.T) {
	r := gofight.New()

	userRepo := repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("role", "User")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})

	// Trying to access catering route with wrong permissions
	// Should throw an error
	r.GET("/images").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusForbidden, r.Code)
		})

	userResult2, _ := userRepo.GetByKey("role", "Catering administrator")
	jwt2, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult2.ID.String()})

	// Trying to access catering with right permissions
	// Should be success
	r.GET("/images").
		SetCookie(gofight.H{
			"jwt": jwt2,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
