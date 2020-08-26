package tests

import (
	"net/http"
	"testing"

	"github.com/Aiscom-LLC/meals-api/src/delivery"
	"github.com/Aiscom-LLC/meals-api/src/delivery/middleware"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/go-playground/assert/v2"
)

func TestAddAddress(t *testing.T) {
	r := gofight.New()

	var clientRepo = repository.NewClientRepo()
	var userRepo = repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()

	// Trying to create new address
	// Should be success
	r.POST("/clients/"+clientID+"/addresses").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"city":   "Sochi",
			"street": "Gagar",
			"house":  "31",
			"floor":  2,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to create new address with invalid data
	// Should return an error
	r.POST("/clients/"+clientID+"/addresses").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"city": "Sochi",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to create new address with non-valid
	// Should return an error
	r.POST("/clients/qwerty/addresses").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"city":   "Sochi",
			"street": "Gagar",
			"house":  "31",
			"floor":  2,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})
}

func TestGetAddress(t *testing.T) {
	r := gofight.New()

	var clientRepo = repository.NewClientRepo()
	var userRepo = repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()

	// Trying to get list of addresses
	// Should be success
	r.GET("/clients/"+clientID+"/addresses").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to get list of addresses with non-valid client ID
	// Should return an error
	r.GET("/clients/qwwerty/addresses").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "client with that ID is not found", errorValue)
		})
}

func TestDeleteAddress(t *testing.T) {
	r := gofight.New()

	var clientRepo = repository.NewClientRepo()
	var userRepo = repository.NewUserRepo()
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	simpleUser, _ := userRepo.GetByKey("email", "user1@meals.com")
	userJWT, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: simpleUser.ID.String()})
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	var addressID string
	var addressID2 string

	// Trying to create 2 new address
	// Should be success
	r.POST("/clients/"+clientID+"/addresses").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"city":   "Sochi",
			"street": "Gagar",
			"house":  "31",
			"floor":  2,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			addressID, _ = jsonparser.GetString(data, "id")
			assert.Equal(t, http.StatusOK, r.Code)
		})

	r.POST("/clients/"+clientID+"/addresses").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"city":   "Sochi",
			"street": "Gagar",
			"house":  "31",
			"floor":  2,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			addressID2, _ = jsonparser.GetString(data, "id")
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to delete address
	// Should be success
	r.DELETE("/clients/"+clientID+"/addresses/"+addressID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to delete non-existing address
	// Should return an error
	r.DELETE("/clients/"+clientID+"/addresses/qwewr").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "address not found", errorValue)
		})

	// Trying to delete address by user without permissions
	// Should return an error
	r.DELETE("/clients/"+clientID+"/addresses/"+addressID2).
		SetCookie(gofight.H{
			"jwt": userJWT,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusForbidden, r.Code)
			assert.Equal(t, "no permissions", errorValue)
		})
}

func TestUpdateAddress(t *testing.T) {
	r := gofight.New()

	var clientRepo = repository.NewClientRepo()
	var userRepo = repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()
	simpleUser, _ := userRepo.GetByKey("email", "user1@meals.com")
	userJWT, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: simpleUser.ID.String()})
	var addressID string

	// Trying to create new address
	// Should be success
	r.POST("/clients/"+clientID+"/addresses").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"city":   "Sochi",
			"street": "Gagara",
			"house":  "31",
			"floor":  2,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			addressID, _ = jsonparser.GetString(data, "id")
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to change data for address
	// Should be success
	r.PUT("/clients/"+clientID+"/addresses/"+addressID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"city":   "stdasg",
			"floor":  1,
			"house":  "sdsring",
			"street": "stridasng",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to change address with invalid data
	// Should return an error
	r.PUT("/clients/"+clientID+"/addresses/"+addressID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"citdy":   "stdasg",
			"flodor":  1,
			"houdse":  "sdsring",
			"stdreet": "stridasng",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to change data for address by user without permissions
	// Should return an error
	r.PUT("/clients/"+clientID+"/addresses/"+addressID).
		SetCookie(gofight.H{
			"jwt": userJWT,
		}).
		SetJSON(gofight.D{
			"citdy":   "stdasg",
			"flodor":  1,
			"houdse":  "sdsring",
			"stdreet": "stridasng",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusForbidden, r.Code)
			assert.Equal(t, "no permissions", errorValue)
		})

	// Trying to update address which doesn't exist
	// Should return an error
	r.PUT("/clients/"+clientID+"/addresses/qwerty").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"city":   "stdasg",
			"floor":  1,
			"house":  "sdsring",
			"street": "stridasng",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})
}
