package tests

import (
	"encoding/json"
	"github.com/Aiscom-LLC/meals-api/src/delivery"
	"github.com/Aiscom-LLC/meals-api/src/delivery/middleware"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/Aiscom-LLC/meals-api/src/schemes/response"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAddClient(t *testing.T) {
	r := gofight.New()

	userRepo := repository.NewUserRepo()
	cateringRepo := repository.NewCateringRepo()
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()

	// Create new client
	r.POST("/caterings/"+cateringID+"/clients").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "newclientname",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			name, _ := jsonparser.GetString(data, "name")
			assert.Equal(t, http.StatusCreated, r.Code)
			assert.Equal(t, "newclientname", name)
		})

	r.POST("/caterings/"+cateringID+"/clients").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "newclientname",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "client with that name already exist", errorValue)
		})
}

func TestGetClients(t *testing.T) {
	r := gofight.New()

	userRepo := repository.NewUserRepo()
	result, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: result.ID.String()})

	// Trying to get list of clients
	r.GET("/clients?limit=5").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			var result response.GetClients
			_ = json.Unmarshal(data, &result)
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, 1, result.Page)
		})

	// Testing pagination params
	r.GET("/clients?limit=5&page=2").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			var result response.GetClients
			_ = json.Unmarshal(data, &result)
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, 2, result.Page)
		})
}

func TestDeleteClient(t *testing.T) {
	r := gofight.New()

	userRepo := repository.NewUserRepo()
	cateringRepo := repository.NewCateringRepo()
	result, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{result.ID.String()})
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()
	var clientID string

	// Create new client
	r.POST("/caterings/"+cateringID+"/clients").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "yandex",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			clientID, _ = jsonparser.GetString(data, "id")
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Deleting client
	r.DELETE("/clients/"+clientID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to delete client which already deleted
	r.DELETE("/clients/"+clientID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "client not found", errorValue)
		})
}

func TestUpdateClient(t *testing.T) {
	r := gofight.New()

	userRepo := repository.NewUserRepo()
	cateringRepo := repository.NewCateringRepo()
	result, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{result.ID.String()})
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()
	var clientID string

	// Create new client
	r.POST("/caterings/"+cateringID+"/clients").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "google",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			clientID, _ = jsonparser.GetString(data, "id")
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Trying to change name of the client
	// Should be success
	r.PUT("/clients/"+clientID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "testingupdate",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to change name of the client with name that already exist in DB
	// Should throw an error
	r.PUT("/clients/"+clientID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "Dymi",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to change a name of client with non-valid ID
	// Should throw an error
	r.PUT("/clients/qwerty").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "newclientname",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to change name of client with non-exist ID
	// Should throw an error
	fakeID := uuid.NewV4()
	r.PUT("/clients/"+fakeID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "testclient",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, r.Code)
		})
}
