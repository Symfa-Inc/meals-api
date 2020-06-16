package tests

import (
	"encoding/json"
	"github.com/appleboy/gofight"
	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"go_api/src/delivery"
	"go_api/src/delivery/middleware"
	"go_api/src/repository/catering"
	"go_api/src/repository/user"
	"go_api/src/schemes/response/meal"
	"net/http"
	"testing"
)

func TestGetMeals(t *testing.T) {
	r := gofight.New()

	userResult, _ := user.GetUserByKey("email", "admin@meals.com")
	cateringResult, _ := catering.GetCateringByKey("name", "Twiist")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	// Testing validation of params
	// Should throw an error
	r.GET("/meals/"+cateringResult.ID.String()+"?startDate=2010-01-01").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Testing non-existing catering ID
	// Should throw an error
	fakeId, _ := uuid.NewV4()
	r.GET("/meals/"+fakeId.String()+"?startDate=2010-01-01").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, "record not found", errorValue)
			assert.Equal(t, http.StatusNotFound, r.Code)
		})

	// Trying to get list of meals for catering
	// Should be success
	r.GET("/meals/"+cateringResult.ID.String()+"?startDate=2010-01-01&endDate=2021-01-01&limit=5").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			var result meal.GetMealsModel
			_ = json.Unmarshal(data, &result)
			assert.Equal(t, 5, len(result.Items))
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
