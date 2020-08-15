package tests

import (
	"encoding/json"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/Aiscom-LLC/meals-api/src/delivery"
	"github.com/Aiscom-LLC/meals-api/src/delivery/middleware"
	"github.com/Aiscom-LLC/meals-api/src/schemes/response"
	"net/http"
	"testing"
)

func TestAddCatering(t *testing.T) {
	r := gofight.New()

	cateringName := "newcatering"
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	// Create new catering
	r.POST("/caterings").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": cateringName,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Trying to create new catering with same name
	r.POST("/caterings").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": cateringName,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "catering with that name already exist", errorValue)
		})
}

func TestDeleteCatering(t *testing.T) {
	r := gofight.New()

	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	cateringResult, _ := cateringRepo.GetByKey("name", "Gink")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	// Deleting catering
	r.DELETE("/caterings/"+cateringResult.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to delete catering which already deleted
	r.DELETE("/caterings/"+cateringResult.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "catering not found", errorValue)
		})
}

func TestGetCaterings(t *testing.T) {
	r := gofight.New()

	result, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{result.ID.String()})

	// Trying to get list of caterings
	r.GET("/caterings?limit=10").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			var result response.GetCaterings
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
			var result response.GetCaterings
			_ = json.Unmarshal(data, &result)
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, 15, len(result.Items))
			assert.Equal(t, 2, result.Page)
		})
}

func TestUpdateCatering(t *testing.T) {
	r := gofight.New()

	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	result, _ := cateringRepo.GetByKey("name", "Telpod")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	// Trying to change name of the catering
	// Should be success
	r.PUT("/caterings/"+result.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "newcateringname",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to change name of the catering with name that already exist in DB
	// Should throw an error
	r.PUT("/caterings/"+result.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "Twiist",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to change a name of catering with non-valid ID
	// Should throw an error
	r.PUT("/caterings/qwerty").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "newcateringname",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to change name of catering with non-exist ID
	// Should throw an error
	fakeID := uuid.NewV4()
	r.PUT("/caterings/"+fakeID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "testcatering",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, r.Code)
		})
}
