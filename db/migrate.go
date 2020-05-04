package main

import (
	"fmt"
	"go_api/src/config"
	"go_api/src/models/seeds"
	"go_api/src/models/user"
	"go_api/src/types"
)

func main() {
	fmt.Println("=== CREATING TYPES ===")

	createTypes()
	fmt.Println("=== TYPES ARE CREATED ===")

	config.DB.DropTableIfExists(&users.User{}, &seeds.Seed{})
	config.DB.AutoMigrate(
		&users.User{}, &seeds.Seed{},
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
