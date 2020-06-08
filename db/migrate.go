package main

import (
	"fmt"
	"go_api/db/seeds/dev"
	"go_api/src/config"
	"go_api/src/models"
	"go_api/src/types"
)

func main() {
	fmt.Println("=== CREATING TYPES ===")
	createTypes()
	fmt.Println("=== TYPES ARE CREATED ===")

	migrate()
	fmt.Println("=== ADD MIGRATIONS ===")

	dev.CreateAdmin()
	dev.CreateUsers();
	dev.CreateCaterings();
}

func migrate() {
	config.DB.DropTableIfExists(
		&models.Catering{},
		&models.Seed{},
		&models.User{},
	)

	config.DB.AutoMigrate(
		&models.User{},
		&models.Seed{},
		&models.Catering{},
	)
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
