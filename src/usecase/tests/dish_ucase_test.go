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
)

func TestAddDish(t *testing.T) {
	r := gofight.New()

	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringId := cateringResult.ID.String()

	categoryResult, _ := dishCategoryRepo.GetByKey("name", "гарнир", cateringId)

	// Trying to add dish to non-existing catering
	// Should throw an error
	fakeId, _ := uuid.NewV4()
	r.POST("/caterings/"+fakeId.String()+"/dishes").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"desc":   "Очень вкусный",
			"name":   "тест",
			"price":  120,
			"weight": 250,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to add dish to non-existing dish category
	// Should throw an error
	r.POST("/caterings/"+cateringId+"/dishes").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"desc":   "Очень вкусный",
			"name":   "тест",
			"price":  120,
			"weight": 250,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to create new dish
	// Should be success
	r.POST("/caterings/"+cateringId+"/dishes").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"categoryId": categoryResult.ID,
			"desc":       "Очень вкусный",
			"name":       "тест",
			"price":      120,
			"weight":     250,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to create same dish in same category
	// Should throw an error
	r.POST("/caterings/"+cateringId+"/dishes").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"categoryId": categoryResult.ID,
			"desc":       "Очень вкусный",
			"name":       "тест",
			"price":      120,
			"weight":     250,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "this dish already exist in that category", errorValue)
		})
}

func TestDeleteDish(t *testing.T) {
	r := gofight.New()

	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringId := cateringResult.ID.String()

	dishCategoryResult, _ := dishCategoryRepo.GetByKey("name", "супы", cateringId)
	dishCategoryId := dishCategoryResult.ID.String()

	dishResult, _ := dishRepo.GetByKey("name", "доширак", cateringId, dishCategoryId)

	fakeId, _ := uuid.NewV4()

	// Trying to delete non-existing dish
	// Should throw an error
	r.DELETE("/caterings/"+cateringId+"/dishes/"+fakeId.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "dish not found", errorValue)
		})

	// Trying to delete existing dish
	// Should be success
	r.DELETE("/caterings/"+cateringId+"/dishes/"+dishResult.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to delete soft deleted dish
	// Should throw an error
	r.DELETE("/caterings/"+cateringId+"/dishes/"+dishResult.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "dish not found", errorValue)
		})
}

func TestGetDishes(t *testing.T) {
	r := gofight.New()

	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringId := cateringResult.ID.String()

	dishCategoryResult, _ := dishCategoryRepo.GetByKey("name", "супы", cateringId)
	dishCategoryId := dishCategoryResult.ID.String()

	fakeId, _ := uuid.NewV4()

	// Trying to get dishes with non-existing catering ID
	// Should throw an error
	r.GET("/caterings/"+fakeId.String()+"/dishes?categoryId="+dishCategoryId).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "catering with that ID doesn't exist", errorValue)
		})

	// Trying to get dishes with non-existing category ID
	// Should throw an error
	r.GET("/caterings/"+cateringId+"/dishes?categoryId="+fakeId.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "category with that ID doesn't exist", errorValue)
		})

	// Trying to get dishes with all valid values
	// Should be success
	r.GET("/caterings/"+cateringId+"/dishes?categoryId="+dishCategoryResult.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestUpdateDish(t *testing.T) {
	r := gofight.New()

	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringId := cateringResult.ID.String()

	dishCategoryResult, _ := dishCategoryRepo.GetByKey("name", "супы", cateringId)
	dishCategoryId := dishCategoryResult.ID.String()

	dishResult, _ := dishRepo.GetByKey("name", "борщ", cateringId, dishCategoryId)
	dishId := dishResult.ID.String()

	fakeId, _ := uuid.NewV4()

	// Trying to update dish for non-existing catering
	// Should throw an error
	r.PUT("/caterings/"+fakeId.String()+"/dishes/"+dishId).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"categoryId": dishCategoryId,
			"desc":       "Очень острый",
			"name":       "супер доширак",
			"price":      120,
			"weight":     250,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "catering not found", errorValue)
		})

	// Trying to update dish with non-existing dish category id
	// Should throw an error
	r.PUT("/caterings/"+cateringId+"/dishes/"+dishId).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"categoryId": fakeId,
			"desc":       "Очень острый",
			"name":       "супер доширак",
			"price":      120,
			"weight":     250,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "dish category not found", errorValue)
		})

	// Trying to update dish with non-existing dish id
	// Should throw an error
	r.PUT("/caterings/"+cateringId+"/dishes/"+fakeId.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"categoryId": dishCategoryId,
			"desc":       "Очень острый",
			"name":       "супер доширак",
			"price":      120,
			"weight":     250,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "dish not found", errorValue)
		})

	// Trying to update dish with all valid values
	// Should be success
	r.PUT("/caterings/"+cateringId+"/dishes/"+dishId).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"categoryId": dishCategoryId,
			"desc":       "Очень острый",
			"name":       "супер доширак",
			"price":      120,
			"weight":     250,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})
}
