package tests

import (
	"encoding/json"
	"github.com/appleboy/gofight"
	"go_api/src/schemes/response/catering"
	"github.com/stretchr/testify/assert"
	"go_api/src/delivery"
	"go_api/src/delivery/middleware"
	"go_api/src/models"
	"net/http"
	"testing"
)

func TestGetCaterings(t *testing.T) {
	r := gofight.New()
	user, _ := models.GetUserByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{user.ID.String()})

	// Trying to get list of caterings
	r.GET("/caterings").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			var result catering.GetCaterings
			_ = json.Unmarshal(data, &result)
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, result.Size, len(result.Items))
			assert.Equal(t, 1, result.Page)
		})

	// Testing pagination params
	r.GET("/caterings?page=2&limit=25").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			var result catering.GetCaterings
			_ = json.Unmarshal(data, &result)
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, result.Size, len(result.Items))
			assert.Equal(t, 2, result.Page)
		})
}
