package tests

import (
	"encoding/json"
	"fmt"
	"github.com/appleboy/gofight"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"go_api/src/delivery"
	"go_api/src/delivery/middleware"
	"go_api/src/models"
	"go_api/src/repository/catering"
	"go_api/src/repository/user"
	"net/http"
	"testing"
	"time"
)

func TestUpdateMeal(t *testing.T) {
	r := gofight.New()

	var resultFirst []models.Meal

	userResult, _ := user.GetUserByKey("email", "admin@meals.com")
	cateringResult, _ := catering.GetCateringByKey("name", "Qiao")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	trunc := 24 * time.Hour
	var validMealsArray []models.Meal
	for i := 0; i < 7; i++ {
		meal := models.Meal{
			Date:       time.Now().AddDate(0, 0, i).Truncate(trunc),
			CateringID: cateringResult.ID,
		}
		validMealsArray = append(validMealsArray, meal)
	}
	jsonValidMealArray, _ := json.Marshal(validMealsArray)

	// Creates new meals for catering
	r.POST("/meals/"+cateringResult.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetBody(string(jsonValidMealArray)).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			_ = json.Unmarshal(data, &resultFirst)
			assert.Equal(t, http.StatusCreated, r.Code)
			assert.Equal(t, len(validMealsArray), len(resultFirst))
		})

	// Trying to update meal with non existing date
	// Should be success
	r.PUT("/meals/"+resultFirst[0].ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date":       time.Now().AddDate(0, 0, 10).Truncate(trunc).UTC(),
			"cateringId": cateringResult.ID.String(),
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			fmt.Println("=============TEST=================")
			fmt.Println(cateringResult.ID)
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to update the same date with already existing date
	// Should throw an error
	r.PUT("/meals/"+resultFirst[0].ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date":       time.Now().AddDate(0, 0, 10).Truncate(trunc).UTC(),
			"cateringId": cateringResult.ID.String(),
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "this date already exist", errorValue)
		})

	//Trying to create meal before today
	//Should throw an error
	r.PUT("/meals/"+resultFirst[0].ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date":       time.Now().AddDate(0, 0, -1).Truncate(trunc).UTC(),
			"cateringId": cateringResult.ID,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "can't add meals to previous dates", errorValue)
		})
}
