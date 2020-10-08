package api

// User struct
type User struct{}

// NewUser returns pointer to user struct
// with all methods
func NewUser() *User {
	return &User{}
}
