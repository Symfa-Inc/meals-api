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
	"net/http"
	"testing"
)

func TestGetCateringSchedules(t *testing.T) {
	r := gofight.New()

	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

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
	fakeID, _ := uuid.NewV4()
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

	var result []domain.CateringSchedule

	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{userResult.ID.String()})

	// Trying to get list of schedules
	// Should be success
	r.GET("/caterings/"+cateringID+"/schedules").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			_ = json.Unmarshal(data, &result)
			assert.Equal(t, http.StatusOK, r.Code)
		})

	schedule := result[0]
	scheduleID := schedule.ID.String()

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

	// Trying to update existing schedule
	// Should be success
	r.PUT("/caterings/"+cateringID+"/schedules/"+scheduleID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"isWorking": false,
			"start":     "16:00",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "end date can't be earlier than start date", errorValue)
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
			data := []byte(r.Body.String())
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
			data := []byte(r.Body.String())
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
	fakeID, _ := uuid.NewV4()
	r.PUT("/caterings/"+fakeID.String()+"/schedules/"+scheduleID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"end": "12:00",
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
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
			data := []byte(r.Body.String())
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
