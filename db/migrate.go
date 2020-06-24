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

	dev.CreateAdmin()
	dev.CreateUsers()
	dev.CreateCaterings()
	dev.CreateMeals()
	dev.CreateDishCategories()
	dev.CreateDishes()
}

func migrate() {
	config.DB.DropTableIfExists(
		&domain.MealDish{},
		&domain.Dish{},
		&domain.DishCategory{},
		&domain.Meal{},
		&domain.Catering{},
		&domain.Seed{},
		&domain.User{},
	)

	config.DB.AutoMigrate(
		&domain.User{},
		&domain.Seed{},
		&domain.Catering{},
		&domain.Meal{},
		&domain.DishCategory{},
		&domain.Dish{},
		&domain.MealDish{},
	)
}

func addDbConstraints() {
	config.DB.Model(&domain.Meal{}).AddForeignKey("catering_id", "caterings(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.DishCategory{}).AddForeignKey("catering_id", "caterings(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.Dish{}).AddForeignKey("catering_id", "caterings(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.Dish{}).AddForeignKey("dish_category_id", "dish_categories(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.MealDish{}).AddForeignKey("meal_id", "meals(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.MealDish{}).AddForeignKey("dish_id", "dishes(id)", "CASCADE", "CASCADE")
}

func createTypes() {
	userTypesQuery := fmt.Sprintf("CREATE TYPE user_roles AS ENUM ('%s', '%s', '%s', '%s')",
		types.UserRoleEnum.SuperAdmin,
		types.UserRoleEnum.CompanyAdmin,
		types.UserRoleEnum.ClientAdmin,
		types.UserRoleEnum.User,
	)

	config.DB.Exec(userTypesQuery)
}
