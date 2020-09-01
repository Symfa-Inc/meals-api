package tests

import (
	"github.com/Aiscom-LLC/meals-api/api"
	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/go-playground/assert/v2"
	"net/http"
	"testing"
)

func TestAddOrder(t *testing.T) {
	r := gofight.New()

	type Dish struct {
		Amount int    `json:"amount"`
		ID     string `json:"dishId"`
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
	dishRepo := repository.NewDishRepo()
	dishResult, _, _ := dishRepo.GetByKey("name", "доширак", cateringID, categoryID)
	dishID := dishResult.ID.String()
	user, _ := userRepo.GetByKey("email", "user1@meals.com")
	userID := user.ID.String()
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userID})
	var order Order
	var dish Dish
	dish.Amount = 1
	dish.ID = dishID
	order.Comment = "cool comment"
	order.Items = append(order.Items, dish)

	// Trying to create new order
	// Should be success
	r.POST("/users/"+userID+"/orders?date=2120-06-20T00%3A00%3A00Z").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSONInterface(order).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Trying to create new order with same data
	// Should return an error
	r.POST("/users/"+userID+"/orders?date=2120-06-20T00%3A00%3A00Z").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSONInterface(order).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "order for current day already created", errorValue)
		})

	// Trying to create new order with non-valid data
	// Should return an error
	dish.Amount = 1
	order.Items = append(order.Items, dish)

	r.POST("/users/"+userID+"/orders?date=2220-06-20T00%3A00%3A00Z").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSONInterface(order).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "can't add 2 same dishes, please increment amount field instead", errorValue)
		})
}

//func TestGetOrder(t *testing.T) {
//	r := gofight.New()
//
//	userRepo := repository.NewUserRepo()
//	userResult, _ := userRepo.GetByKey("email", "user1@meals.com")
//	userID := userResult.ID.String()
//	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userID})
//
//	// Trying to get list of order
//	// Should be success
//	r.GET("/users/"+userID+"/order?date=date=2120-06-20T00%3A00%3A00Z").
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//				assert.Equal(t, http.StatusOK, r.Code)
//		})
//
//	// Trying to get list of order with non-existing date
//	// Should return an error
//	r.GET("/users/"+userID+"/orders?date=2121-06-20T00%3A00%3A00Z").
//		SetCookie(gofight.H{
//			"jwt": jwt,
//		}).
//		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
//			data := r.Body.Bytes()
//			errorValue, _ := jsonparser.GetString(data, "error")
//			assert.Equal(t, http.StatusNotFound, r.Code)
//			assert.Equal(t, "record not found", errorValue)
//		})
//}
