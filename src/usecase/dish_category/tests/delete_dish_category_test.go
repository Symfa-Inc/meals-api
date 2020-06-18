package tests

import (
	"github.com/appleboy/gofight"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"go_api/src/delivery"
	"go_api/src/delivery/middleware"
	"go_api/src/repository"
	"net/http"
	"testing"
)

func TestDeleteDishCategory(t *testing.T) {
	r := gofight.New()

	userResult, _ := repository.GetUserByKey("email", "admin@meals.com")
	cateringResult, _ := repository.GetCateringByKey("name", "Pheast")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
	cateringId := cateringResult.ID.String()
	var id string

	// Creates new dish category
	// Should be success
	r.POST("/caterings/"+cateringId+"/dish-category").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "закуски",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			id, _ = jsonparser.GetString(data, "id")
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Trying to delete new category dish
	// Should be success
	r.DELETE("/caterings/"+cateringId+"/dish-category/"+id).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to delete already deleted category
	// Should throw error
	r.DELETE("/caterings/"+cateringId+"/dish-category/"+id).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "category not found", errorValue)
		})
}
