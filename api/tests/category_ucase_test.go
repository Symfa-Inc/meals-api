package tests

import (
	"github.com/Aiscom-LLC/meals-api/api"
	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAddCategory(t *testing.T) {
	r := gofight.New()

	clientRepo := repository.NewClientRepo()
	userRepo := repository.NewUserRepo()
	cateringRepo := repository.NewCateringRepo()
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})

	// Trying to add category to non-existing ID
	// Should throw error
	r.POST("/caterings/qwerty/clients/"+clientID+"/categories").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "закуски",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to add category to existing ID
	// Should be success
	r.POST("/caterings/"+cateringResult.ID.String()+"/clients/"+clientID+"/categories").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "закуски",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to add already existing category
	// Should throw error
	r.POST("/caterings/"+cateringResult.ID.String()+"/clients/"+clientID+"/categories").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "закуски",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "this category already exist", errorValue)
		})
}

func TestDeleteCategory(t *testing.T) {
	r := gofight.New()

	clientRepo := repository.NewClientRepo()
	userRepo := repository.NewUserRepo()
	cateringRepo := repository.NewCateringRepo()
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	cateringID := cateringResult.ID.String()
	clientID := clientResult.ID.String()
	fakeID := uuid.NewV4()
	var categoryID string

	// Trying to create new dish category
	// Should be success
	r.POST("/caterings/"+cateringID+"/clients/"+clientID+"/categories").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "testFeed",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			categoryID, _ = jsonparser.GetString(data, "id")
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to delete new category dish
	// Should be success
	r.DELETE("/caterings/"+cateringID+"/clients/"+clientID+"/categories/"+categoryID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to delete category with non-existing client
	// Should return an error
	r.DELETE("/caterings/"+cateringID+"/clients/"+fakeID.String()+"/categories/"+categoryID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "category not found", errorValue)
		})

	// Trying to delete category with non-existing catering
	// Should return an error
	r.DELETE("/caterings/"+cateringID+"/clients/"+clientID+"/categories/"+fakeID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "category not found", errorValue)
		})
}

func TestGetCategories(t *testing.T) {
	r := gofight.New()

	clientRepo := repository.NewClientRepo()
	userRepo := repository.NewUserRepo()
	cateringRepo := repository.NewCateringRepo()
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	id := cateringResult.ID.String()
	clientID := clientResult.ID.String()
	fakeID := uuid.NewV4()

	// Trying to get categories from non-existing catering
	// Should throw error
	r.GET("/caterings/"+fakeID.String()+"/clients/"+clientID+"/categories?date=2020-08-25%2018%3A13%3A55").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "catering with that ID is not found", errorValue)
		})

	// Trying to get categories from existing catering
	// Should be success
	r.GET("/caterings/"+id+"/clients/"+clientID+"/categories?date=2020-08-25%2018%3A13%3A55").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestUpdateCategory(t *testing.T) {
	r := gofight.New()

	clientRepo := repository.NewClientRepo()
	userRepo := repository.NewUserRepo()
	cateringRepo := repository.NewCateringRepo()
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	cateringID := cateringResult.ID.String()
	clientID := clientResult.ID.String()
	var categoryID string

	// Trying to update non-existing dish category
	// Should throw an error
	fakeID := uuid.NewV4()
	r.PUT("/caterings/"+cateringID+"/clients/"+clientID+"/categories/"+fakeID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "zxcvb",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, r.Code)
		})

	// Trying to update non-existing client
	// Should throw an error
	r.PUT("/caterings/"+cateringID+"/clients/"+fakeID.String()+"/categories/"+fakeID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "zxcvb",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, r.Code)
		})

	// Trying to update non-existing catering
	// Should throw an error
	r.PUT("/caterings/"+fakeID.String()+"/clients/"+clientID+"/categories/"+fakeID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "zxcvb",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, r.Code)
		})

	// Posting new dish category to update it
	// Should be success
	r.POST("/caterings/"+cateringID+"/clients/"+clientID+"/categories").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "qwerty",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			categoryID, _ = jsonparser.GetString(data, "id")
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to update new dish category
	// Should be success
	r.PUT("/caterings/"+cateringID+"/clients/"+clientID+"/categories/"+categoryID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "zxcvb",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to update the same category with already existing category
	// Should throw an error
	r.PUT("/caterings/"+cateringID+"/clients/"+clientID+"/categories/"+categoryID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "салаты",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "category with that name already exist", errorValue)
		})
}
