package tests

import (
	"encoding/json"
	"github.com/appleboy/gofight"
	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"go_api/src/delivery"
	"go_api/src/delivery/middleware"
	"go_api/src/domain"
	"go_api/src/schemes/response"
	"net/http"
	"testing"
	"time"
)

func TestAddMeals(t *testing.T) {
	r := gofight.New()

	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	cateringResult, _ := cateringRepo.GetByKey("name", "Pyrami")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	trunc := 24 * time.Hour
	var validMealsArray []domain.Meal
	for i := 0; i < 7; i++ {
		meal := domain.Meal{
			Date:       time.Now().AddDate(0, 0, i).Truncate(trunc),
			CateringID: cateringResult.ID,
		}
		validMealsArray = append(validMealsArray, meal)
	}

	var nonValidMealArray []domain.Meal
	for i := 0; i < 7; i++ {
		meal := domain.Meal{
			Date:       time.Now().AddDate(0, 0, i-1).UTC().Truncate(trunc),
			CateringID: cateringResult.ID,
		}
		nonValidMealArray = append(nonValidMealArray, meal)
	}

	jsonNonValidMealArray, _ := json.Marshal(nonValidMealArray)
	jsonValidMealArray, _ := json.Marshal(validMealsArray)

	// Trying do add non-valid days
	// Should throw an error
	r.POST("/caterings/"+cateringResult.ID.String()+"/meals").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetBody(string(jsonNonValidMealArray)).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "item 1 has wrong date (can't use previous dates)", errorValue)
		})

	// Trying to add valid meals
	// Should be success
	r.POST("/caterings/"+cateringResult.ID.String()+"/meals").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetBody(string(jsonValidMealArray)).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			var result []domain.Meal
			_ = json.Unmarshal(data, &result)
			assert.Equal(t, http.StatusCreated, r.Code)
			assert.Equal(t, len(validMealsArray), len(result))
		})

	r.POST("/caterings/"+cateringResult.ID.String()+"/meals").
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

func TestGetMeals(t *testing.T) {
	r := gofight.New()

	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	// Testing validation of params
	// Should throw an error
	r.GET("/caterings/"+cateringResult.ID.String()+"/meals?startDate=2010-01-01").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Testing non-existing catering ID
	// Should throw an error
	fakeId, _ := uuid.NewV4()
	r.GET("/caterings/"+fakeId.String()+"/meals?startDate=2010-01-01").
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
	r.GET("/caterings/"+cateringResult.ID.String()+"/meals?startDate=2010-01-01&endDate=2021-01-01&limit=5").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			var result response.GetMealsModel
			_ = json.Unmarshal(data, &result)
			assert.Equal(t, 5, len(result.Items))
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestUpdateMeal(t *testing.T) {
	r := gofight.New()

	var resultFirst []domain.Meal

	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	cateringResult, _ := cateringRepo.GetByKey("name", "Qiao")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
	cateringId := cateringResult.ID.String()

	trunc := 24 * time.Hour
	var validMealsArray []domain.Meal
	for i := 0; i < 7; i++ {
		meal := domain.Meal{
			Date:       time.Now().AddDate(0, 0, i).Truncate(trunc),
			CateringID: cateringResult.ID,
		}
		validMealsArray = append(validMealsArray, meal)
	}
	jsonValidMealArray, _ := json.Marshal(validMealsArray)

	// Creates new meals for catering
	r.POST("/caterings/"+cateringId+"/meals").
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

	// Trying to update meal with existing date
	// Should be success
	r.PUT("/caterings/"+cateringId+"/meals/"+resultFirst[0].ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date": time.Now().AddDate(0, 0, 10).Truncate(trunc).UTC(),
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to update the same date with already existing date
	// Should throw an error
	r.PUT("/caterings/"+cateringId+"/meals/"+resultFirst[0].ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date": time.Now().AddDate(0, 0, 10).Truncate(trunc).UTC(),
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "this date already exist", errorValue)
		})

	//Trying to create meal before today
	//Should throw an error
	r.PUT("/caterings/"+cateringId+"/meals/"+resultFirst[0].ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date": time.Now().AddDate(0, 0, -1).Truncate(trunc).UTC(),
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "can't add meals to previous dates", errorValue)
		})
}
