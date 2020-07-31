package request

// CateringUserUpdate scheme
type CateringUserUpdate struct {
	FirstName string `json:"firstName,omitempty" example:"Dmitry"`
	LastName  string `json:"lastName,omitempty" example:"Novikov"`
	Email     string `json:"email,omitempty" example:"d.novikov@wellyes.ru"`
	ClientID  string `json:"clientID"`
	Status    string `json:"status"`
}

// ClientUserUpdate scheme
type ClientUserUpdate struct {
	FirstName string `json:"firstName,omitempty" example:"Dmitry"`
	LastName  string `json:"lastName,omitempty" example:"Novikov"`
	Email     string `json:"email,omitempty" example:"d.novikov@wellyes.ru"`
	Floor     int    `json:"floor" example:"5"`
	Role      string `json:"role" example:"User"`
	Status    string `json:"status"`
}