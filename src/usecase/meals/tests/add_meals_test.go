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
	"go_api/src/repository"
	"net/http"
	"testing"
	"time"
)

func TestAddMeals(t *testing.T) {
	r := gofight.New()

	userResult, _ := repository.GetUserByKey("email", "admin@meals.com")
	cateringResult, _ := repository.GetCateringByKey("name", "Pyrami")
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

	var nonValidMealArray []models.Meal
	for i := 0; i < 7; i++ {
		meal := models.Meal{
			Date:       time.Now().AddDate(0, 0, i-1).UTC().Truncate(trunc),
			CateringID: cateringResult.ID,
		}
		nonValidMealArray = append(nonValidMealArray, meal)
	}

	jsonNonValidMealArray, _ := json.Marshal(nonValidMealArray)
	jsonValidMealArray, _ := json.Marshal(validMealsArray)

	// Trying do add non-valid days
	// Should throw an error
	r.POST("/meals/"+cateringResult.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetBody(string(jsonNonValidMealArray)).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			fmt.Println(string(jsonNonValidMealArray))
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "item 1 has wrong date (can't use previous dates)", errorValue)
		})

	// Trying to add valid meals
	// Should be success
	r.POST("/meals/"+cateringResult.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetBody(string(jsonValidMealArray)).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			var result []models.Meal
			_ = json.Unmarshal(data, &result)
			assert.Equal(t, http.StatusCreated, r.Code)
			assert.Equal(t, len(validMealsArray), len(result))
		})

	r.POST("/meals/"+cateringResult.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetBody(string(jsonValidMealArray)).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "item 1 already exist", errorValue)
		})
}
