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
	dev.CreateCategories()
	dev.CreateDishes()
	dev.CreateMealDishes()
}

func migrate() {
	config.DB.DropTableIfExists(
		&domain.MealDish{},
		&domain.Dish{},
		&domain.Category{},
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
		&domain.Category{},
		&domain.Dish{},
		&domain.MealDish{},
	)

}

func addDbConstraints() {
	config.DB.Model(&domain.Meal{}).AddForeignKey("catering_id", "caterings(id)", "CASCADE", "CASCADE")
	config.DB.Model(&domain.Meal{}).AddIndex("idx_meals_date", "date")

	config.DB.Model(&domain.Category{}).AddForeignKey("catering_id", "caterings(id)", "CASCADE", "CASCADE")

	config.DB.Model(&domain.Dish{}).AddForeignKey("category_id", "categories(id)", "CASCADE", "CASCADE")

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
