package data

import (
	"github.com/bxcodec/faker"
	users "go_api/src/models/user"
	"go_api/src/utils"
)

// Creates Users mock data
func CreateUsersData(quantity int) []users.User {
	var usersSlice []users.User
	for i := 0; i < quantity; i++ {
		userItem := users.User{
			Password: utils.HashString("Password12!"),
			Role:     "Delivery administrator",
		}
		_ = faker.FakeData(&userItem)
		usersSlice = append(usersSlice, userItem)
	}
	return usersSlice
}
