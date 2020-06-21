package tests

import (
	"github.com/appleboy/gofight"
	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"go_api/src/delivery"
	"go_api/src/delivery/middleware"
	"go_api/src/repository"
	"net/http"
	"testing"
)

func TestAddDish(t *testing.T) {
	r := gofight.New()

	userResult, _ := repository.GetUserByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	cateringResult, _ := repository.GetCateringByKey("name", "Twiist")
	cateringId := cateringResult.ID.String()

	categoryResult, _ := repository.GetDishCategoryByKey("name", "гарнир", cateringId)

	// Trying to add dish to non-existing catering
	// Should throw an error
	fakeId, _ := uuid.NewV4()
	r.POST("/caterings/"+fakeId.String()+"/dish").
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
	r.POST("/caterings/"+cateringId+"/dish").
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
	r.POST("/caterings/"+cateringId+"/dish").
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
			name, _ := jsonparser.GetString(data, "name")
			assert.Equal(t, http.StatusCreated, r.Code)
			assert.Equal(t, "тест", name)
		})

	// Trying to create same dish in same category
	// Should throw an error
	r.POST("/caterings/"+cateringId+"/dish").
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
