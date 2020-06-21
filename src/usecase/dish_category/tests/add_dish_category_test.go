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

func TestAddDishCategory(t *testing.T) {
	r := gofight.New()

	userResult, _ := repository.GetUserByKey("email", "admin@meals.com")
	cateringResult, _ := repository.GetCateringByKey("name", "Gology")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	// Trying to add category to non-existing ID
	// Should throw error
	r.POST("/caterings/qwerty/dish-categories").
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
	r.POST("/caterings/"+cateringResult.ID.String()+"/dish-categories").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "закуски",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			name, _ := jsonparser.GetString(data, "name")
			assert.Equal(t, http.StatusCreated, r.Code)
			assert.Equal(t, "закуски", name)
		})

	// Trying to add already existing category
	// Should throw error
	r.POST("/caterings/"+cateringResult.ID.String()+"/dish-categories").
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
