package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Aiscom-LLC/meals-api/api"
	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/go-playground/assert/v2"
	uuid "github.com/satori/go.uuid"
)

func TestGetClientSchedules(t *testing.T) {
	r := gofight.New()

	var userRepo = repository.NewUserRepo()
	var cateringRepo = repository.NewClientRepo()
	userResult, _ := userRepo.GetByKey("email", "meals@aisnovations.com")
	clientResult, _ := cateringRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})

	// Trying to get list of schedules
	// Should be success
	r.GET("/clients/"+clientID+"/schedules").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to get list of schedules with non-valid client id
	// Should throw an error
	r.GET("/clients/+clientID/schedules").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to get list of schedules with non-existing catering id
	// Should throw an error
	fakeID := uuid.NewV4()
	r.GET("/clients/"+fakeID.String()+"/schedules").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, r.Code)
		})
}

func TestUpdateClientSchedule(t *testing.T) {
	r := gofight.New()

	var userRepo = repository.NewUserRepo()
	var cateringRepo = repository.NewClientRepo()
	userResult, _ := userRepo.GetByKey("email", "meals@aisnovations.com")
	clientResult, _ := cateringRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	var result []domain.ClientSchedule
	var scheduleID string

	// Trying to get list of schedules
	// Should be success
	r.GET("/clients/"+clientID+"/schedules").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			list, _, _, _ := jsonparser.Get(data, "list")
			_ = json.Unmarshal(list, &result)
			scheduleID = result[0].ID.String()
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to update existing schedule
	// Should be success
	r.PUT("/clients/"+clientID+"/schedules/"+scheduleID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"isWorking": false,
			"start":     "12:00",
			"end":       "14:00",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to update existing schedule with wrong start date
	// Should throw an error
	r.PUT("/clients/"+clientID+"/schedules/"+scheduleID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"start": "19:00",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "end date can't be earlier than start date", errorValue)
		})

	// Trying to update existing schedule with wrong end date
	// Should throw an error
	r.PUT("/clients/"+clientID+"/schedules/"+scheduleID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"end": "19:00",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "new end time can't be later than catering's end time", errorValue)
		})

	// Trying to update schedule with non-valid client id
	// Should throw an error
	r.PUT("/clients/+clientID+/schedules/"+scheduleID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"isWorking": false,
			"start":     "12:00",
			"end":       "14:00",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to update schedule with non-existing catering id
	// Should throw an error
	fakeID := uuid.NewV4()
	r.PUT("/clients/"+fakeID.String()+"/schedules/"+scheduleID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"isWorking": false,
			"start":     "12:00",
			"end":       "14:00",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "client with that id not found", errorValue)
		})

	// Trying to update schedule with non-existing schedule id
	// Should throw an error
	r.PUT("/clients/"+clientID+"/schedules/+scheduleID").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"isWorking": false,
			"start":     "12:00",
			"end":       "14:00",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	r.PUT("/clients/"+clientID+"/schedules/"+fakeID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"isWorking": false,
			"start":     "12:00",
			"end":       "14:00",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "client schedule with that id not found", errorValue)
		})
}
