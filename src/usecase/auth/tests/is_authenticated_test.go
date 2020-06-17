package tests

import (
	"github.com/appleboy/gofight"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"go_api/src/delivery"
	"go_api/src/delivery/middleware"
	"go_api/src/repository"
	"net/http"
	"testing"
)

func TestIsAuthenticated(t *testing.T) {
	r := gofight.New()

	userResult, _ := repository.GetUserByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	// Trying to login without jwt cookie
	r.GET("/is-authenticated").
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "message")
			assert.Equal(t, http.StatusUnauthorized, r.Code)
			assert.Equal(t, "cookie token is empty", errorValue)
		})

	// Trying to login with jwt cookie
	r.GET("/is-authenticated").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		data := []byte(r.Body.String())
		email, _ := jsonparser.GetString(data, "email")
		name, _ := jsonparser.GetString(data, "firstName")
		assert.Equal(t, http.StatusOK, r.Code)
		assert.Equal(t, "admin@meals.com", email)
		assert.Equal(t, "super", name)
	})
}
