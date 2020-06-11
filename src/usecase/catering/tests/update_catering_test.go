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

func TestUpdateCatering(t *testing.T) {
	r := gofight.New()
	user, _ := models.GetUserByKey("email", "admin@meals.com")
	catering, _ := models.GetCateringByKey("name", "Telpod")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{user.ID.String()})

	// Trying to change name of the catering
	r.PUT("/caterings/"+catering.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "newcateringname",
		}).Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		data := []byte(r.Body.String())
		name, _ := jsonparser.GetString(data, "name")
		assert.Equal(t, http.StatusOK, r.Code)
		assert.Equal(t, "newcateringname", name)
	})

	// Trying to change a name of catering with wrong ID
	r.PUT("/caterings/qwerty").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "newcateringname",
		}).Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		data := []byte(r.Body.String())
		fmt.Println(r.Body)
		errorValue, _ := jsonparser.GetString(data, "error")
		assert.Equal(t, http.StatusNotFound, r.Code)
		assert.Equal(t, "catering not found", errorValue)
	})
}
