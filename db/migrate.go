package main

import (
	"fmt"

	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/db/seeds/dev"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository/enums"
)

func main() {
	fmt.Println("=== CREATING TYPES ===")
	createTypes()
	fmt.Println("=== TYPES ARE CREATED ===")

	migrate()
	fmt.Println("=== ADD MIGRATIONS ===")

	addDbConstraints()
	fmt.Println("=== ADD DB CONSTRAINTS ===")

	dev.CreateCaterings()
	dev.CreateCateringSchedules()
	dev.CreateClients()
	dev.CreateClientSchedules()

	dev.CreateAdmin()

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

func migrate() {
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

	config.DB.AutoMigrate(
		&domain.Seed{},
		&domain.User{},
		&domain.Catering{},
		&domain.CateringUser{},
		&domain.CateringSchedule{},
		&domain.Client{},
		&domain.ClientUser{},
		&domain.Address{},
		&domain.ClientSchedule{},
		&domain.Meal{},
		&domain.Category{},
		&domain.Dish{},
		&domain.ImageDish{},
		&domain.Image{},
		&domain.MealDish{},
		&domain.Order{},
		&domain.OrderDishes{},
		&domain.UserOrders{},
	)

}

func addDbConstraints() {
	config.DB.Model(&domain.CateringUser{}).AddForeignKey("catering_id", "caterings(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.CateringUser{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.ClientUser{}).AddForeignKey("client_id", "clients(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.ClientUser{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.Address{}).AddForeignKey("client_id", "clients(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.Client{}).AddForeignKey("catering_id", "caterings(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.CateringSchedule{}).AddForeignKey("catering_id", "caterings(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.ClientSchedule{}).AddForeignKey("client_id", "clients(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.Meal{}).AddForeignKey("catering_id", "caterings(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.Meal{}).AddIndex("idx_meals_date", "date")

	config.DB.Model(&domain.Category{}).AddForeignKey("catering_id", "caterings(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.Category{}).AddForeignKey("client_id", "clients(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.Dish{}).AddForeignKey("category_id", "categories(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.MealDish{}).AddForeignKey("meal_id", "meals(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.MealDish{}).AddForeignKey("dish_id", "dishes(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.ImageDish{}).AddForeignKey("dish_id", "dishes(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.ImageDish{}).AddForeignKey("image_id", "images(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.OrderDishes{}).AddForeignKey("order_id", "orders(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.OrderDishes{}).AddForeignKey("dish_id", "dishes(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.UserOrders{}).AddForeignKey("order_id", "orders(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.UserOrders{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
}

func createTypes() {
	userTypesQuery := fmt.Sprintf("DROP TYPE IF EXISTS user_roles CASCADE;"+
		"CREATE TYPE user_roles AS ENUM ('%s', '%s', '%s', '%s')",
		enums.UserRoleEnum.SuperAdmin,
		enums.UserRoleEnum.CateringAdmin,
		enums.UserRoleEnum.ClientAdmin,
		enums.UserRoleEnum.User,
	)

	companyTypesQuery := fmt.Sprintf("DROP TYPE IF EXISTS company_types CASCADE;"+
		"CREATE TYPE company_types AS ENUM ('%s', '%s')",
		enums.CompanyTypesEnum.Catering,
		enums.CompanyTypesEnum.Client,
	)

	statusTypesQuery := fmt.Sprintf("DROP TYPE IF EXISTS status_types CASCADE;"+
		"CREATE TYPE status_types AS ENUM ('%s', '%s', '%s')",
		enums.StatusTypesEnum.Deleted,
		enums.StatusTypesEnum.Invited,
		enums.StatusTypesEnum.Active,
	)

	orderStatusTypesQuery := fmt.Sprintf("DROP TYPE IF EXISTS order_status_types CASCADE;"+
		"CREATE TYPE order_status_types AS ENUM ('%s', '%s', '%s')",
		enums.OrderStatusTypesEnum.Approved,
		enums.OrderStatusTypesEnum.Canceled,
		enums.OrderStatusTypesEnum.Pending,
	)

	config.DB.Exec(userTypesQuery)
	config.DB.Exec(companyTypesQuery)
	config.DB.Exec(statusTypesQuery)
	config.DB.Exec(orderStatusTypesQuery)
}
