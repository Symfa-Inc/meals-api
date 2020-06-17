package tests

import (
	"github.com/appleboy/gofight"
	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"go_api/src/delivery"
	"go_api/src/delivery/middleware"
	"go_api/src/repository"
	"net/http"
	"testing"
)

func TestUpdateCatering(t *testing.T) {
	r := gofight.New()

	userResult, _ := repository.GetUserByKey("email", "admin@meals.com")
	result, _ := repository.GetCateringByKey("name", "Telpod")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	// Trying to change name of the catering
	r.PUT("/caterings/"+result.ID.String()).
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

	// Trying to change name of the catering with name that already exist in DB
	r.PUT("/caterings/"+result.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "Twiist",
		}).Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, http.StatusBadRequest, r.Code)
	})

	// Trying to change a name of catering with non-valid ID
	r.PUT("/caterings/qwerty").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "newcateringname",
		}).Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, http.StatusBadRequest, r.Code)
	})

	// Trying to change name of catering with non-exist ID
	fakeId, _ := uuid.NewV4()
	r.PUT("/caterings/"+fakeId.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "newcateringname",
		}).Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, http.StatusNotFound, r.Code)
	})
}
