package tests

import (
	"github.com/Aiscom-LLC/meals-api/src/delivery"
	"github.com/Aiscom-LLC/meals-api/src/delivery/middleware"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAddCategory(t *testing.T) {
	r := gofight.New()

	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	// Trying to add category to non-existing ID
	// Should throw error
	r.POST("/caterings/qwerty/clients/"+clientID+"/categories").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "закуски",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "this category already exist", errorValue)
		})
}

//func TestDeleteCategory(t *testing.T) {
//	r := gofight.New()
//
//	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
//	cateringResult, _ := cateringRepo.GetByKey("name", "Pheast")
//	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
//	cateringID := cateringResult.ID.String()
//
//	// Creates new dish category
//	// Should be success
//	r.POST("/caterings/"+cateringID+"/categories").
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		SetJSON(gofight.D{
//			"name": "закуски",
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			assert.Equal(t, http.StatusOK, r.Code)
//		})
//
//	createdCategory, _ := categoryRepo.GetByKey("name", "закуски", cateringID)
//	categoryID := createdCategory.ID.String()
//
//	// Trying to delete new category dish
//	// Should be success
//	r.DELETE("/caterings/"+cateringID+"/categories/"+categoryID).
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			assert.Equal(t, http.StatusNoContent, r.Code)
//		})
//
//	// Trying to delete already deleted category
//	// Should throw error
//	r.DELETE("/caterings/"+cateringID+"/categories/"+categoryID).
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			data := []byte(r.Body.String())
//			errorValue, _ := jsonparser.GetString(data, "error")
//			assert.Equal(t, http.StatusNotFound, r.Code)
//			assert.Equal(t, "category not found", errorValue)
//		})
//}
//
////
//func TestGetCategories(t *testing.T) {
//	r := gofight.New()
//
//	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
//	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
//	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
//	id := cateringResult.ID.String()
//	var fakeID, _ = uuid.NewV4()
//
//	// Trying to get categories from non-existing catering
//	// Should throw error
//	r.GET("/caterings/"+fakeID.String()+"/categories").
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			data := []byte(r.Body.String())
//			errorValue, _ := jsonparser.GetString(data, "error")
//			assert.Equal(t, http.StatusNotFound, r.Code)
//			assert.Equal(t, "catering with that ID is not found", errorValue)
//		})
//
//	// Trying to get categories from existing catering
//	// Should be success
//	r.GET("/caterings/"+id+"/categories").
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			assert.Equal(t, http.StatusOK, r.Code)
//		})
//}
//
//func TestUpdateCategory(t *testing.T) {
//	r := gofight.New()
//
//	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
//	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
//	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
//	cateringID := cateringResult.ID.String()
//
//	// Trying to update non-existing dish category
//	// Should throw an error
//	fakeID := uuid.NewV4()
//	r.PUT("/caterings/"+cateringID+"/categories/"+fakeID.String()).
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		SetJSON(gofight.D{
//			"name": "zxcvb",
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			assert.Equal(t, http.StatusNotFound, r.Code)
//		})
//
//	// Posting new dish category to update it
//	// Should be success
//	r.POST("/caterings/"+cateringID+"/categories").
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		SetJSON(gofight.D{
//			"name": "qwerty",
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			assert.Equal(t, http.StatusOK, r.Code)
//		})
//
//	createdCategory, _ := categoryRepo.GetByKey("name", "qwerty", cateringID)
//	categoryID := createdCategory.ID.String()
//
//	// Trying to update new dish category
//	// Should be success
//	r.PUT("/caterings/"+cateringID+"/categories/"+categoryID).
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		SetJSON(gofight.D{
//			"name": "zxcvb",
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			assert.Equal(t, http.StatusNoContent, r.Code)
//		})
//
//	// Trying to update the same category with already existing category
//	// Should throw an error
//	r.PUT("/caterings/"+cateringID+"/categories/"+categoryID).
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		SetJSON(gofight.D{
//			"name": "салаты",
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			data := []byte(r.Body.String())
//			errorValue, _ := jsonparser.GetString(data, "error")
//			assert.Equal(t, http.StatusBadRequest, r.Code)
//			assert.Equal(t, "category with that name already exist", errorValue)
//		})
//}
