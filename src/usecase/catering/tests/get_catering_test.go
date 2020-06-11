package tests

import (
	"fmt"
	"github.com/appleboy/gofight"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"go_api/src/delivery"
	"go_api/src/delivery/middleware"
	"go_api/src/models"
	"net/http"
	"testing"
)

func TestGetCatering(t *testing.T) {
	r := gofight.New()
	user, _ := models.GetUserByKey("email", "admin@meals.com")
	catering, _ := models.GetCateringByKey("name", "Navir")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{user.ID.String()})

	// Trying to get catering by ID
	r.GET("/caterings/"+catering.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			cateringName, _ := jsonparser.GetString(data, "name")
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "Navir", cateringName)
		})

	// Deleting catering for next tests
	r.DELETE("/caterings/"+catering.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {})

	// Trying to get catering with wrong uuid format
	r.GET("/caterings/uuidtest").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			fmt.Println(r.Code)
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to get catering which already been deleted
	r.GET("/caterings/"+catering.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "record not found", errorValue)
		})
}
