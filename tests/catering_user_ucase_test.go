package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Aiscom-LLC/meals-api/api"
	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
)

func TestAddCateringUser(t *testing.T) {
	r := gofight.New()

	var cateringRepo = repository.NewCateringRepo()
	var userRepo = repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("email", "marianafox@comcubine.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()
	email := "test@mail.ru"
	var newUserID string

	// Trying to create new Catering user
	// Should be success
	r.POST("/caterings/"+cateringResult.ID.String()+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     email,
			"firstName": "Dmitry",
			"lastName":  "Novikov",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			newUserID, _ = jsonparser.GetString(data, "id")
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Trying to create second new Catering user
	// Should be success
	r.POST("/caterings/"+cateringResult.ID.String()+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     "d.novikov@wellyes.ru",
			"firstName": "Dmitry",
			"lastName":  "Novikov",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Trying to create user with non-valid email
	r.POST("/caterings/"+cateringID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     "alineglenn/.omcubine.com",
			"firstName": "ExampleFName",
			"lastName":  "ExampleLName",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "email is not valid", errorValue)
		})

	// Trying to create user which email is already exist
	// Should be success
	r.POST("/caterings/"+cateringID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     "meals@aisnovations.com",
			"firstName": "ExampleFName",
			"lastName":  "ExampleLName",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "user with that email already exist", errorValue)
		})

	// Trying to delete user
	// Should be success
	r.DELETE("/caterings/"+cateringID+"/users/"+newUserID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to create new user with email which already exist but have status "deleted"
	// Should be success
	r.POST("/caterings/"+cateringResult.ID.String()+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     email,
			"firstName": "Dmitry",
			"lastName":  "Novikov",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
		})
}

func TestGetCateringUsers(t *testing.T) {
	r := gofight.New()

	var cateringRepo = repository.NewCateringRepo()
	var userRepo = repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("email", "marianafox@comcubine.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()

	// Trying to get all users
	// Should be success
	r.GET("/caterings/"+cateringID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			var result models.GetCateringUser
			_ = json.Unmarshal(data, &result)
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to get user with non-valid catering ID
	// Should return an error
	r.GET("/caterings/qwerty/users?limit=5").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})
}

func TestDeleteCateringUsers(t *testing.T) {
	r := gofight.New()

	var cateringRepo = repository.NewCateringRepo()
	var userRepo = repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("email", "marianafox@comcubine.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()
	user, _ := userRepo.GetByKey("email", "user2@meals.com")
	userID := user.ID.String()

	// Trying to delete catering user
	// Should be success
	r.DELETE("/caterings/"+cateringID+"/users/"+userID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to delete itself
	// Should return an error
	r.DELETE("/clients/"+cateringID+"/users/"+userResult.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
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
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()
	user, _ := userRepo.GetByKey("email", "user3@meals.com")
	userID := user.ID.String()

	// Trying to change name of user
	// Should be success
	r.PUT("/caterings/"+cateringID+"/users/"+userID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"firstName": "newCoolName",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to change email to invalid
	// Should return an error
	r.PUT("/caterings/"+cateringID+"/users/"+userID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email": "newCoolNamedas",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "email is not valid", errorValue)
		})
}
