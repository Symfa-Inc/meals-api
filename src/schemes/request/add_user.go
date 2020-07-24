package request

// CateringUser scheme
type CateringUser struct {
	FirstName string `json:"firstName,omitempty" example:"Dmitry"`
	LastName  string `json:"lastName,omitempty" example:"Novikov"`
	Email     string `json:"email,omitempty" example:"d.novikov@wellyes.ru"`
	ClientID  string `json:"clientID"`
}

type ClientUser struct {
	FirstName string `json:"firstName,omitempty" example:"Dmitry"`
	LastName  string `json:"lastName,omitempty" example:"Novikov"`
	Email     string `json:"email,omitempty" example:"d.novikov@wellyes.ru"`
	Floor     int    `json:"floor" example:"5"`
	Role      string `json:"role" example:"User"`
}
