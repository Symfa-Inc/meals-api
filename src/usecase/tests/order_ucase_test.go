package tests

import (
	"github.com/Aiscom-LLC/meals-api/src/delivery"
	"github.com/Aiscom-LLC/meals-api/src/delivery/middleware"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/go-playground/assert/v2"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"testing"
)

func TestAddOrder(t *testing.T) {
	r := gofight.New()

	type Dish struct {
		Amount int    `json:"amount"`
		DishId string `json:"dishId"`
	}
	type Order struct {
		Comment string `json:"comment"`
		Items   []Dish `json:"items"`
	}

	userRepo := repository.NewUserRepo()
	cateringRepo := repository.NewCateringRepo()
	categoryRepo := repository.NewCategoryRepo()
	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()
	categoryResult, _ := categoryRepo.GetByKey("name", "гарнир", cateringID)
	categoryID := categoryResult.ID.String()
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	dishRepo := repository.NewDishRepo()
	dishResult, _, _ := dishRepo.GetByKey("name", "доширак", cateringID, categoryID)
	dishID := dishResult.ID.String()
	newUser, _ := userRepo.GetByKey("email", "user1@meals.com")
	userID := newUser.ID.String()
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	var test Order
	var test1 Dish
	test1.Amount = 1
	test1.DishId = dishID
	test.Comment = "cool comment"
	test.Items = append(test.Items, test1)

	// Trying to create new order
	// Should be success
	r.POST("/users/"+userID+"/orders?date=2120-06-20T00%3A00%3A00Z").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSONInterface(test).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Trying to create new order with same data
	// Should return an error
	r.POST("/users/"+userID+"/orders?date=2120-06-20T00%3A00%3A00Z").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSONInterface(test).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "order for current day already created", errorValue)
		})

	// Trying to create new order with non-valid data
	// Should return an error
	var test2 Dish
	fakeID := uuid.NewV4()
	test1.Amount = 1
	test2.DishId = fakeID.String()
	test.Items = append(test.Items, test1)

	r.POST("/users/"+userID+"/orders?date=2220-06-20T00%3A00%3A00Z\n").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSONInterface(test).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "can't add 2 same dishes, please increment amount field instead", errorValue)
		})
}

func TestGetOrder(t *testing.T) {
	r := gofight.New()

	userRepo := repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("email", "admin@meals.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	newUser, _ := userRepo.GetByKey("email", "user1@meals.com")
	newUserID := newUser.ID.String()

	// Trying to get list of order
	// Should be success
	r.GET("/users/"+newUserID+"/orders?date=2120-06-20T00%3A00%3A00Z").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
	})

	// Trying to ger list of order with non-existing date
	// Should return an error
	r.GET("/users/"+newUserID+"/orders?date=2121-06-20T00%3A00%3A00Z").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(delivery.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := []byte(r.Body.String())
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "record not found", errorValue)
		})
}
