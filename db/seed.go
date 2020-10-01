package main

import (
	"github.com/Aiscom-LLC/meals-api/db/seeds/dev"
)

func main()  {

	dev.CreateCaterings()
	dev.CreateCateringSchedules()
	dev.CreateClients()
	dev.CreateClientSchedules()

	dev.CreateUsers()
	dev.CreateCateringUsers()
	dev.CreateClientUsers()
	dev.CreateMeals()
	dev.CreateCategories()
	dev.CreateDishes()

	dev.CreateImages()

	dev.CreateMealDishes()
	dev.CreateImageDishes()
	dev.CreateAddresses()
}
