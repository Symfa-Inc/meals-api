package seeds

import (
	"go_api/src/config"
	"go_api/src/types"
	"go_api/src/models/user"
	"go_api/src/models/seeds"
	"go_api/src/utils"
)

const seedName string = "init admin"

// CreateAdmin init admin user
func CreateAdmin() {

	seedExists := config.DB.Where("name = ?", seedName).First()

	if seedExists == nil {
		seed := seeds.Seed{
			name: seedName
		}
	
		superAdmin := users.User{
			FirstName: "super",
			LastName: "admin",
			Email: "admin@melas.com",
			Password: utils.HashString("Password12!"),
			Role: types.UserRoleEnum.SuperAdmin,
		}
	
	
		config.DB.Create(&superAdmin)
		config.DB.Create(&seed)	
	}
}