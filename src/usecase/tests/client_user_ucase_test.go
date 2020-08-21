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

func TestAddClientUser(t *testing.T) {
	r := gofight.New()

	var clientRepo = repository.NewClientRepo()
	var userRepo = repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()

	// Trying to create new user
	// Should be success
	r.POST("/clients/"+clientID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     "newUserEmail@mail.ru",
			"firstName": "newFirstName",
			"floor":     5,
			"lastName":  "NewLastName",
			"role":      "User",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Create user with exist email
	// Should return BadRequest error
	r.POST("/clients/"+clientID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     "newUserEmail@mail.ru",
			"firstName": "newFirstName",
			"floor":     5,
			"lastName":  "NewLast_Name",
			"role":      "User",
		}).

	// Trying to create user with already existing email
	// Must return an error
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "user with that email already exist", errorValue)
		})

	// Trying to create user with invalid email
	// Should return BadRequest error
	r.POST("/clients/"+clientID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     "SecondNewUserEmail.mail.ru",
			"firstName": "newFirstName",
			"floor":     5,
			"lastName":  "NewLast_Name",
			"role":      "User",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "email is not valid", errorValue)
		})
}

func TestGetClientUsers(t *testing.T) {
	r := gofight.New()

	var clientRepo = repository.NewClientRepo()
	var userRepo = repository.NewUserRepo()
	result, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{result.ID.String()})
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()

	//Trying to get list of users
	//Should be success
	r.GET("/clients/"+clientID+"/users?limit=5").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			var result response.GetClientUser
			_ = json.Unmarshal(data, &result)
			assert.Equal(t, http.StatusOK, r.Code)
		})
	clientID += "1"
	r.GET("/clients/"+clientID+"/users?limit=5").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "pq: invalid input syntax for type uuid: \""+clientID+"\"", errorValue)
		})
}

func TestDeleteClientUsers(t *testing.T) {
	r := gofight.New()

	var clientRepo = repository.NewClientRepo()
	var userRepo = repository.NewUserRepo()
	result, _ := userRepo.GetByKey("email", "admin@meals.com")
	admin1, _ := userRepo.GetByKey("email", "marianafox@comcubine.com")
	admin2, _ := userRepo.GetByKey("email", "maggietodd@comcubine.com")
	clientAdmin, _ := userRepo.GetByKey("email", "melodybond@comcubine.com")

	adminJWT, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{clientAdmin.ID.String()})
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{result.ID.String()})
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()
	var userID string

	//	Trying to create new client user
	// Should be success
	r.POST("/clients/"+clientID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     "newUserEmails@mail.ru",
			"firstName": "newFirstName",
			"floor":     5,
			"lastName":  "NewLastName",
			"role":      "User",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			userID, _ = jsonparser.GetString(data, "id")
			assert.Equal(t, http.StatusCreated, r.Code)
		})

		// Delete client simple user
		// Must be success
	r.DELETE("/clients/"+clientID+"/users/"+userID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Delete client admin user
	// Must be success
	r.DELETE("/clients/"+clientID+"/users/"+admin1.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})
	r.DELETE("/clients/"+clientID+"/users/"+admin2.ID.String()).
		SetCookie(gofight.H{
			"jwt": adminJWT,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to delete SUPER USER
	// Must return error
	r.DELETE("/clients/"+clientID+"/users/"+result.ID.String()).
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

func TestUpdateClientUser(t *testing.T) {
	r := gofight.New()

	var clientRepo = repository.NewClientRepo()
	var userRepo = repository.NewUserRepo()
	result, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{result.ID.String()})
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()
	var userID string

	//	Trying to create new client user
	// Should be success
	r.POST("/clients/"+clientID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     "newUserrEmails@mail.ru",
			"firstName": "newFirstName",
			"floor":     5,
			"lastName":  "NewLastName",
			"role":      "User",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			userID, _ = jsonparser.GetString(data, "id")
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Trying to change name for user
	// Should be success
	r.PUT("/clients/"+clientID+"/users/"+userID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"lastName": "testName",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
}