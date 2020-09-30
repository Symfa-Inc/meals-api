package request

// CateringUserUpdate scheme
type CateringUserUpdate struct {
	FirstName string `json:"firstName,omitempty" example:"Dmitry"`
	LastName  string `json:"lastName,omitempty" example:"Novikov"`
	Email     string `json:"email,omitempty" example:"d.novikov@wellyes.ru"`
}

// ClientUserUpdate scheme
type ClientUserUpdate struct {
	FirstName string `json:"firstName,omitempty" example:"Dmitry"`
	LastName  string `json:"lastName,omitempty" example:"Novikov"`
	Email     string `json:"email,omitempty" example:"d.novikov@wellyes.ru"`
	Floor     *int   `json:"floor" example:"5"`
	Role      string `json:"role" example:"User"`
}

// ForgotPassword scheme
type ForgotPassword struct {
	Email string `json:"email,omitempty" example:"admin@meals.com"`
}