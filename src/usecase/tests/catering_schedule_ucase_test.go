package tests

import (
	"encoding/json"
	"github.com/Aiscom-LLC/meals-api/src/delivery"
	"github.com/Aiscom-LLC/meals-api/src/delivery/middleware"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/go-playground/assert/v2"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"testing"
)

func TestGetCateringSchedules(t *testing.T) {
	r := gofight.New()

	var userRepo = repository.NewUserRepo()
	var cateringRepo = repository.NewCateringRepo()
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})

	// Trying to get list of schedules
	// Should be success
	r.GET("/caterings/"+cateringID+"/schedules").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to get list of schedules with non-valid catering id
	// Should throw an error
	r.GET("/caterings/qwerty/schedules").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to get list of schedules with non-existing catering id
	// Should throw an error
	fakeID := uuid.NewV4()
	r.GET("/caterings/"+fakeID.String()+"/schedules").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, r.Code)
		})
}

func TestUpdateCateringSchedule(t *testing.T) {
	r := gofight.New()

	var userRepo = repository.NewUserRepo()
	var cateringRepo = repository.NewCateringRepo()
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	var result []domain.CateringSchedule
	var scheduleID string

	// Trying to get list of schedules
	// Should be success
	r.GET("/caterings/"+cateringID+"/schedules").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := (r.Body.Bytes())
			list, _, _, _ := jsonparser.Get(data, "list")
			_ = json.Unmarshal(list, &result)
			scheduleID = result[0].ID.String()
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to update existing schedule
	// Should be success
	r.PUT("/caterings/"+cateringID+"/schedules/"+scheduleID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"isWorking": false,
			"start":     "12:00",
			"end":       "14:00",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to update existing schedule with wrong start date
	// Should throw an error
	r.PUT("/caterings/"+cateringID+"/schedules/"+scheduleID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"start": "16:00",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := (r.Body.Bytes())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "end date can't be earlier than start date", errorValue)
		})

	// Trying to update existing schedule with wrong end date
	// Should throw an error
	r.PUT("/caterings/"+cateringID+"/schedules/"+scheduleID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"end": "11:00",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := (r.Body.Bytes())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "end date can't be earlier than start date", errorValue)
		})

	// Trying to update schedule with non-valid catering id
	// Should throw an error
	r.PUT("/caterings/qwerty/schedules/"+scheduleID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"end": "12:00",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to update schedule with non-existing catering id
	// Should throw an error
	fakeID := uuid.NewV4()
	r.PUT("/caterings/"+fakeID.String()+"/schedules/"+scheduleID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"end": "12:00",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := (r.Body.Bytes())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "catering with that id not found", errorValue)
		})

	// Trying to update schedule with non-existing schedule id
	// Should throw an error
	r.PUT("/caterings/"+cateringID+"/schedules/"+fakeID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"end": "12:00",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := (r.Body.Bytes())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "schedule with that id not found", errorValue)
		})

	// Trying to update schedule with non-valid schedule id
	// Should throw an error
	r.PUT("/caterings/"+cateringID+"/schedules/qwerty").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"end": "12:00",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})
}
