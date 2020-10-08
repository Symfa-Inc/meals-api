package services

// UserService struct
type UserService struct{}

// NewOrderService return pointer to order struct
// with all methods
func NewUserService() *UserService {
	return &UserService{}
}
