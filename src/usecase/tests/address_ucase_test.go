package tests

import (
	"github.com/Aiscom-LLC/meals-api/src/delivery"
	"github.com/Aiscom-LLC/meals-api/src/delivery/middleware"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/go-playground/assert/v2"
	"net/http"
	"testing"
)

func TestAddAddress(t *testing.T) {
	r := gofight.New()

	var clientRepo = repository.NewClientRepo()
	var userRepo = repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
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

	// Trying to create new address with invalid or not full data
	// Should be either invalid data or something like "with not all fields filled"
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

	// Trying to create new address without authorize
	// Should return an error
	r.POST("/clients/qwerty/addresses").
		SetJSON(gofight.D{
			"city":   "Sochi",
			"street": "Gagar",
			"house":  "31",
			"floor":  2,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusUnauthorized, r.Code)
		})
}

func TestGetAddress(t *testing.T) {
	r := gofight.New()

	var clientRepo = repository.NewClientRepo()
	var userRepo = repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
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
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "client with that ID is not found", errorValue)
		})

	// Trying to get list of addresses without authorize
	// Should return an error
	r.GET("/clients/qwwerty/addresses").
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusUnauthorized, r.Code)
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
	simpleJwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{simpleUser.ID.String()})
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
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
			data := []byte(r.Body.String())
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
			data := []byte(r.Body.String())
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

	// Trying to delete not found address
	// Should return an error
	r.DELETE("/clients/"+clientID+"/addresses/qwewr").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "address not found", errorValue)
		})

	// Trying to delete address by simple user
	// Should return an error
	r.DELETE("/clients/"+clientID+"/addresses/"+addressID2).
		SetCookie(gofight.H{
			"jwt": simpleJwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
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
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()
	simpleUser, _ := userRepo.GetByKey("email", "user1@meals.com")
	simpleJwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{simpleUser.ID.String()})
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
			data := []byte(r.Body.String())
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

	// Trying to change invalid data for address
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

	// Trying to change data for address by simple user
	// Should return an error
	r.PUT("/clients/"+clientID+"/addresses/"+addressID).
		SetCookie(gofight.H{
			"jwt": simpleJwt,
		}).
		SetJSON(gofight.D{
			"citdy":   "stdasg",
			"flodor":  1,
			"houdse":  "sdsring",
			"stdreet": "stridasng",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
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

	// Trying to change data for address without authorization
	// Should return an error
	r.PUT("/clients/"+clientID+"/addresses/"+addressID).
		SetJSON(gofight.D{
			"citdy":   "stdasg",
			"flodor":  1,
			"houdse":  "sdsring",
			"stdreet": "stridasng",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "message")
			assert.Equal(t, http.StatusUnauthorized, r.Code)
			assert.Equal(t, "token contains an invalid number of segments", errorValue)
		})
}
