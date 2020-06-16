package tests

import (
	"encoding/json"
	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"go_api/src/delivery"
	"go_api/src/delivery/middleware"
	"go_api/src/repository/user"
	"go_api/src/schemes/response/catering"
	"net/http"
	"testing"
)

func TestGetCaterings(t *testing.T) {
	r := gofight.New()

	result, _ := user.GetUserByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{result.ID.String()})

	// Trying to get list of caterings
	r.GET("/caterings?limit=10").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			var result catering.GetCaterings
			_ = json.Unmarshal(data, &result)
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, 10, len(result.Items))
			assert.Equal(t, 1, result.Page)
		})

	// Testing pagination params
	r.GET("/caterings?limit=15&page=2").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			var result catering.GetCaterings
			_ = json.Unmarshal(data, &result)
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, 15, len(result.Items))
			assert.Equal(t, 2, result.Page)
		})
}
