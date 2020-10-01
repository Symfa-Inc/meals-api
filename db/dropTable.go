package main

import (
	"fmt"
	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
)

func main()  {
	config.DB.DropTableIfExists(
		&domain.UserOrders{},
		&domain.OrderDishes{},
		&domain.Order{},
		&domain.MealDish{},
		&domain.ImageDish{},
		&domain.Image{},
		&domain.Dish{},
		&domain.Category{},
		&domain.Meal{},
		&domain.Address{},
		&domain.ClientSchedule{},
		&domain.ClientUser{},
		&domain.Client{},
		&domain.CateringSchedule{},
		&domain.CateringUser{},
		&domain.Catering{},
		&domain.User{},
		&domain.Seed{},
	)

	fmt.Println("=== Tables deleted ====")
}