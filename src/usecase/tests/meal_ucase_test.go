package tests
//
//import (
//	"github.com/appleboy/gofight/v2"
//	"github.com/buger/jsonparser"
//	uuid "github.com/satori/go.uuid"
//	"github.com/stretchr/testify/assert"
//	"github.com/Aiscom-LLC/meals-api/src/delivery"
//	"github.com/Aiscom-LLC/meals-api/src/delivery/middleware"
//	"net/http"
//	"testing"
//	"time"
//)
//
//var dishesIDArray []string
//
//func TestAddMeals(t *testing.T) {
//	r := gofight.New()
//
//	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
//	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
//	cateringID := cateringResult.ID.String()
//	categoryResult, _ := categoryRepo.GetByKey("name", "супы", cateringResult.ID.String())
//	dishesResult, _, _ := dishRepo.Get(cateringID, categoryResult.ID.String())
//	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
//
//	trunc := 24 * time.Hour
//
//	for _, dish := range dishesResult {
//		dishID := dish.ID.String()
//		dishesIDArray = append(dishesIDArray, dishID)
//	}
//
//	// Trying to add new meal with previous date
//	// Should throw an error
//	r.POST("/caterings/"+cateringID+"/meals").
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		SetJSON(gofight.D{
//			"date":   time.Now().AddDate(0, 0, -1).UTC().Truncate(trunc),
//			"dishes": dishesIDArray,
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			data := []byte(r.Body.String())
//			errorValue, _ := jsonparser.GetString(data, "error")
//			assert.Equal(t, http.StatusBadRequest, r.Code)
//			assert.Equal(t, "item has wrong date (can't use previous dates)", errorValue)
//		})
//
//	//Trying to add valid meals
//	//Should be success
//	r.POST("/caterings/"+cateringResult.ID.String()+"/meals").
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		SetJSON(gofight.D{
//			"date":   time.Now().AddDate(0, 0, 10).UTC().Truncate(trunc),
//			"dishes": dishesIDArray,
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			assert.Equal(t, http.StatusCreated, r.Code)
//		})
//
//	// Trying to add meal with already existing date
//	// Should throw an errro
//	r.POST("/caterings/"+cateringResult.ID.String()+"/meals").
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		SetJSON(gofight.D{
//			"date":   time.Now().AddDate(0, 0, 10).UTC().Truncate(trunc),
//			"dishes": dishesIDArray,
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			data := []byte(r.Body.String())
//			errorValue, _ := jsonparser.GetString(data, "error")
//			assert.Equal(t, http.StatusBadRequest, r.Code)
//			assert.Equal(t, "item already exist", errorValue)
//		})
//}
//
//func TestGetMeals(t *testing.T) {
//	r := gofight.New()
//
//	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
//	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
//	cateringID := cateringResult.ID.String()
//	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
//
//	// Testing validation of params
//	// Should throw an error
//	r.GET("/caterings/"+cateringID+"/meals?qwerty=qwerty").
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			assert.Equal(t, http.StatusBadRequest, r.Code)
//		})
//
//	// Testing non-existing catering ID
//	// Should throw an error
//	fakeID := uuid.NewV4()
//	r.GET("/caterings/"+fakeID.String()+"/meals?mealId=qwerty").
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			data := []byte(r.Body.String())
//			errorValue, _ := jsonparser.GetString(data, "error")
//			assert.Equal(t, "record not found", errorValue)
//			assert.Equal(t, http.StatusNotFound, r.Code)
//		})
//
//	//Trying to get meal for catering
//	//Should be success
//	trunc := time.Hour * 24
//	date := time.Now().AddDate(0, 0, 10).Truncate(trunc).UTC().Format(time.RFC3339)
//	r.GET("/caterings/"+cateringResult.ID.String()+"/meals?date="+date).
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		SetJSON(gofight.D{
//			"date": time.Now().AddDate(0, 0, 10).UTC().Truncate(trunc),
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			assert.Equal(t, http.StatusOK, r.Code)
//		})
//}
//
//func TestUpdateMeal(t *testing.T) {
//	r := gofight.New()
//
//	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
//	cateringResult, _ := cateringRepo.GetByKey("name", "Qiao")
//	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})
//	cateringID := cateringResult.ID.String()
//	trunc := time.Hour * 24
//	date := time.Now().AddDate(0, 0, 10).Truncate(trunc).UTC().Format(time.RFC3339)
//	meal, _, _ := mealRepo.GetByKey("date", date)
//	mealID := meal.ID.String()
//
//	// Trying to update meal with no dishes field
//	// Should throw an error
//	r.PUT("/caterings/"+cateringID+"/meals/"+mealID).
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		SetJSON(gofight.D{
//			"test": "123",
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			assert.Equal(t, http.StatusBadRequest, r.Code)
//		})
//
//	// Trying to update meal with dishes array
//	// Should should be success
//	r.PUT("/caterings/"+cateringID+"/meals/"+mealID).
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		SetJSON(gofight.D{
//			"dishes": dishesIDArray[:2],
//		}).
//		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			assert.Equal(t, http.StatusNoContent, r.Code)
//		})
//}
