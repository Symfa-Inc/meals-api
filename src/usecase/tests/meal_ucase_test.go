package tests

import (
	"github.com/appleboy/gofight"
	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"go_api/src/delivery"
	"go_api/src/delivery/middleware"
	"net/http"
	"testing"
	"time"
)

func TestAddMeals(t *testing.T) {
	r := gofight.New()

	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	cateringResult, _ := cateringRepo.GetByKey("name", "Pyrami")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	trunc := 24 * time.Hour

	r.POST("/caterings/"+cateringResult.ID.String()+"/meals").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date": time.Now().AddDate(0, 0, -1).UTC().Truncate(trunc),
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "item has wrong date (can't use previous dates)", errorValue)
		})

	// Trying to add valid meals
	// Should be success
	r.POST("/caterings/"+cateringResult.ID.String()+"/meals").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date": time.Now().AddDate(0, 0, 1).UTC().Truncate(trunc),
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
		})
	// Trying to add meal with already existing date
	// Should throw an errro
	r.POST("/caterings/"+cateringResult.ID.String()+"/meals").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date": time.Now().AddDate(0, 0, 1).UTC().Truncate(trunc),
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "item already exist", errorValue)
		})
}

func TestGetMeals(t *testing.T) {
	r := gofight.New()

	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
	var mealId string

	// Testing validation of params
	// Should throw an error
	r.GET("/caterings/"+cateringResult.ID.String()+"/meals?qwerty=qwerty").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Testing non-existing catering ID
	// Should throw an error
	fakeId, _ := uuid.NewV4()
	r.GET("/caterings/"+fakeId.String()+"/meals?mealId=qwerty").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, "record not found", errorValue)
			assert.Equal(t, http.StatusNotFound, r.Code)
		})

	// Creating new meal to retrieve its ID
	trunc := 24 * time.Hour
	r.POST("/caterings/"+cateringResult.ID.String()+"/meals").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date": time.Now().AddDate(0, 0, 10).UTC().Truncate(trunc),
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			mealId, _ = jsonparser.GetString(data, "id")
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	//Trying to get meal for catering
	//Should be success
	r.GET("/caterings/"+cateringResult.ID.String()+"/meals?mealId="+mealId).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestUpdateMeal(t *testing.T) {
	r := gofight.New()

	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	cateringResult, _ := cateringRepo.GetByKey("name", "Qiao")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
	cateringId := cateringResult.ID.String()
	var mealId string

	// Creating new meal to retrieve its ID
	trunc := 24 * time.Hour
	r.POST("/caterings/"+cateringResult.ID.String()+"/meals").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date": time.Now().AddDate(0, 0, 10).UTC().Truncate(trunc),
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			mealId, _ = jsonparser.GetString(data, "id")
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Trying to update meal with existing date
	// Should be success
	r.PUT("/caterings/"+cateringId+"/meals/"+mealId).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date": time.Now().AddDate(0, 0, 11).Truncate(trunc).UTC(),
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to update the same date with already existing date
	// Should throw an error
	r.PUT("/caterings/"+cateringId+"/meals/"+mealId).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date": time.Now().AddDate(0, 0, 11).Truncate(trunc).UTC(),
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "this date already exist", errorValue)
		})

	//Trying to create meal before today
	//Should throw an error
	r.PUT("/caterings/"+cateringId+"/meals/"+mealId).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date": time.Now().AddDate(0, 0, -1).Truncate(trunc).UTC(),
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "can't add meals to previous dates", errorValue)
		})
}
