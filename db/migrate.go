package main

import (
	"fmt"
	"go_api/db/seeds/dev"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/types"
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
	dev.CreateMeals()
	dev.CreateCategories()
	dev.CreateDishes()
	dev.CreateImages()
	dev.CreateMealDishes()
	dev.CreateImageDishes()
}

func migrate() {
	config.DB.DropTableIfExists(
		&domain.MealDish{},
		&domain.ImageDish{},
		&domain.Image{},
		&domain.Dish{},
		&domain.Category{},
		&domain.Meal{},
		&domain.User{},
		&domain.ClientSchedule{},
		&domain.Client{},
		&domain.CateringSchedule{},
		&domain.Catering{},
		&domain.Seed{},
	)

	config.DB.AutoMigrate(
		&domain.Catering{},
		&domain.CateringSchedule{},
		&domain.Client{},
		&domain.ClientSchedule{},
		&domain.User{},
		&domain.Seed{},
		&domain.Meal{},
		&domain.Category{},
		&domain.Dish{},
		&domain.ImageDish{},
		&domain.Image{},
		&domain.MealDish{},
	)

}

func addDbConstraints() {
	config.DB.Model(&domain.User{}).AddForeignKey("catering_id", "caterings(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.User{}).AddForeignKey("client_id", "clients(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.Client{}).AddForeignKey("catering_id", "caterings(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.CateringSchedule{}).AddForeignKey("catering_id", "caterings(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.ClientSchedule{}).AddForeignKey("client_id", "clients(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.Meal{}).AddForeignKey("catering_id", "caterings(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.Meal{}).AddIndex("idx_meals_date", "date")

	config.DB.Model(&domain.Category{}).AddForeignKey("catering_id", "caterings(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.Dish{}).AddForeignKey("category_id", "categories(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.MealDish{}).AddForeignKey("meal_id", "meals(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.MealDish{}).AddForeignKey("dish_id", "dishes(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.ImageDish{}).AddForeignKey("dish_id", "dishes(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.ImageDish{}).AddForeignKey("image_id", "images(id)", "CASCADE", "CASCADE")
}

func createTypes() {
	userTypesQuery := fmt.Sprintf("CREATE TYPE user_roles AS ENUM ('%s', '%s', '%s', '%s')",
		types.UserRoleEnum.SuperAdmin,
		types.UserRoleEnum.CateringAdmin,
		types.UserRoleEnum.ClientAdmin,
		types.UserRoleEnum.User,
	)

	companyTypesQuery := fmt.Sprintf("CREATE TYPE company_types AS ENUM ('%s', '%s')",
		types.CompanyTypesEnum.Catering,
		types.CompanyTypesEnum.Client,
	)

	statusTypesQuery := fmt.Sprintf("CREATE TYPE status_types AS ENUM ('%s', '%s', '%s')",
		types.StatusTypesEnum.Deleted,
		types.StatusTypesEnum.Invited,
		types.StatusTypesEnum.Active,
	)

	config.DB.Exec(userTypesQuery)
	config.DB.Exec(companyTypesQuery)
	config.DB.Exec(statusTypesQuery)
}
