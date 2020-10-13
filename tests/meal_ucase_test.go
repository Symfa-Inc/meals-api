package tests

import (
	"net/http"
	"testing"

	"github.com/Aiscom-LLC/meals-api/api"
	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/go-playground/assert/v2"
	uuid "github.com/satori/go.uuid"
)

func TestAddMeal(t *testing.T) {
	r := gofight.New()

	dishRepo := repository.NewDishRepo()
	userRepo := repository.NewUserRepo()
	categoryRepo := repository.NewCategoryRepo()
	cateringRepo := repository.NewCateringRepo()
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()
	categoryResult, _ := categoryRepo.GetByKey("name", "гарнир", cateringID)
	categoryID := categoryResult.ID.String()
	userResult, _ := userRepo.GetByKey("email", "meals@aisnovations.com")
	dishResult, _, _ := dishRepo.GetByKey("name", "доширак", cateringID, categoryID)
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	var dishIDs []string
	dishIDs = append(dishIDs, dishResult.ID.String())

	// Trying to create meal
	// Should be success
	r.POST("/caterings/"+cateringID+"/clients/"+categoryID+"/meals").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date":   "2120-06-20T00:00:00Z",
			"dishes": dishIDs,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Trying to create meal with non-valid count of dishes
	// Should return an error
	r.POST("/caterings/"+cateringID+"/clients/"+categoryID+"/meals").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date":   "2120-06-20T00:00:00Z",
			"dishes": "",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to create meal with non-valid date
	// Should return an error
	r.POST("/caterings/"+cateringID+"/clients/"+categoryID+"/meals").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date":   "1120-06-20T00:00:00Z",
			"dishes": dishIDs,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "item has wrong date (can't use previous dates)", errorValue)
		})

	// Trying to create meal with non-existing dishes
	// Should return an error
	var fakeDishID [1]string
	fakeDishID[0] = "441477c2-d17f-47f3-b20c-0b22626385ce"
	r.POST("/caterings/"+cateringID+"/clients/"+categoryID+"/meals").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"date":   "2120-06-20T00:00:00Z",
			"dishes": fakeDishID,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, r.Code)
		})
}

func TestGetMeal(t *testing.T) {
	r := gofight.New()

	userRepo := repository.NewUserRepo()
	categoryRepo := repository.NewCategoryRepo()
	cateringRepo := repository.NewCateringRepo()
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()
	categoryResult, _ := categoryRepo.GetByKey("name", "гарнир", cateringID)
	categoryID := categoryResult.ID.String()
	userResult, _ := userRepo.GetByKey("email", "meals@aisnovations.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})

	// Trying to get list of meal
	// Should be success
	r.GET("/caterings/"+cateringID+"/clients/"+categoryID+"/meals?date=2120-06-20T00%3A00%3A00Z").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to get list of meal with non-valid catering
	// Should return an error
	fakeID := uuid.NewV4()
	r.GET("/caterings/"+fakeID.String()+"/clients/"+categoryID+"/meals?date=2120-06-20T00%3A00%3A00Z").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "record not found", errorValue)
		})

	// Trying to get list of meal with non-valid datetime
	// Should return an error
	r.GET("/caterings/"+cateringID+"/clients/"+categoryID+"/meals?date=11120-06-20T00%3A00%3A00Z").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "can't parse the date", errorValue)
		})
}
