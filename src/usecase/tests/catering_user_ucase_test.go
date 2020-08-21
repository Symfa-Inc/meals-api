package tests

import (
	"encoding/json"
	"github.com/Aiscom-LLC/meals-api/src/delivery"
	"github.com/Aiscom-LLC/meals-api/src/delivery/middleware"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/Aiscom-LLC/meals-api/src/schemes/response"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAddCateringUser(t *testing.T) {
	r := gofight.New()

	var cateringRepo = repository.NewCateringRepo()
	var userRepo = repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("email", "marianafox@comcubine.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
	cateringResult, _:= cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()

	//Trying to create new Catering user
	//Must return success
	r.POST("/caterings/"+cateringResult.ID.String()+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email": "d.novikov@wellyes.ru",
			"firstName": "Dmitry",
			"lastName": "Novikov",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Create user which email is already exist
	// Should return error
	r.POST("/caterings/"+cateringID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email": "alineglenn/.omcubine.com",
			"firstName": "ExampleFName",
			"lastName": "ExampleLName",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "email is not valid", errorValue)
	})

	r.POST("/caterings/"+cateringID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email": "admin@meals.com",
			"firstName": "ExampleFName",
			"lastName": "ExampleLName",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "user with that email already exists", errorValue)
		})
}

func TestGetCateringUsers(t *testing.T) {
	r := gofight.New()

	var cateringRepo = repository.NewCateringRepo()
	var userRepo = repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("email", "marianafox@comcubine.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
	cateringResult, _:= cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()

	// Trying to get all users
	// Should be success
	r.GET("/caterings/"+cateringID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
	}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			var result response.GetCateringUser
			_ = json.Unmarshal(data, &result)
			assert.Equal(t, http.StatusOK, r.Code)
	})

	// Trying to get users on broken URI
	// Must return an error
	cateringID += "1"
	r.GET("/caterings/"+cateringID+"/users?limit=5").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "pq: invalid input syntax for type uuid: \""+cateringID+"\"", errorValue)
		})
}

func TestDeleteCateringUsers(t *testing.T) {
	r := gofight.New()

	var cateringRepo = repository.NewCateringRepo()
	var userRepo = repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("email", "marianafox@comcubine.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
	cateringResult, _:= cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()
	user, _ := userRepo.GetByKey("email", "user2@meals.com")
	userID := user.ID.String()

	// Trying to delete catering user
	// Should be success
	r.DELETE("/caterings/"+cateringID+"/users/"+userID).
		SetCookie(gofight.H{
			"jwt": jwt,
	}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
	})

	// Trying to delete itself
	// Must return an error
	r.DELETE("/clients/"+cateringID+"/users/"+userResult.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "can't delete yourself", errorValue)
	})
}

func TestUpdateClientUsers(t *testing.T) {
	r := gofight.New()

	var cateringRepo = repository.NewCateringRepo()
	var userRepo = repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("email", "marianafox@comcubine.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()
	user, _ := userRepo.GetByKey("email", "user3@meals.com")
	userID := user.ID.String()

	// Try to change name of user
	// Should be success
	r.PUT("/caterings/"+cateringID+"/users/"+userID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"firstName": "newCoolName",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	r.PUT("/caterings/"+cateringID+"/users/"+userID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email": "newCoolNamedas",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "email is not valid", errorValue)
		})
}