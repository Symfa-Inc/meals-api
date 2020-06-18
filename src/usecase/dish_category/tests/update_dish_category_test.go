package tests

import (
	"encoding/json"
	"github.com/appleboy/gofight"
	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"go_api/src/delivery"
	"go_api/src/delivery/middleware"
	"go_api/src/models"
	"go_api/src/repository"
	"net/http"
	"testing"
)

func TestUpdateDishCategory(t *testing.T) {
	r := gofight.New()

	var resultFirst models.DishCategory

	userResult, _ := repository.GetUserByKey("email", "admin@meals.com")
	cateringResult, _ := repository.GetCateringByKey("name", "Twiist")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
	cateringId := cateringResult.ID.String()
	// Trying to update non-existing dish category
	// Should throw an error
	fakeId, _ := uuid.NewV4()
	r.PUT("/caterings/"+cateringId+"/dish-category/"+fakeId.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "zxcvb",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, r.Code)
		})

	// Posting new dish category to update it
	// Should be success
	r.POST("/caterings/"+cateringId+"/dish-category").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "qwerty",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			_ = json.Unmarshal(data, &resultFirst)
			assert.Equal(t, http.StatusCreated, r.Code)
			assert.Equal(t, "qwerty", resultFirst.Name)
		})

	// Trying to update new dish category
	// Should be success
	r.PUT("/caterings/"+cateringId+"/dish-category/"+resultFirst.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "zxcvb",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			name, _ := jsonparser.GetString(data, "name")
			assert.Equal(t, http.StatusOK, r.Code)
			assert.Equal(t, "zxcvb", name)
		})

	// Trying to update the same category with already existing category
	// Should throw an error
	r.PUT("/caterings/"+cateringId+"/dish-category/"+resultFirst.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"name": "салаты",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "this category already exist", errorValue)
		})
}
