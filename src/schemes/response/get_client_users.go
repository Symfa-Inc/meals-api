package response

import uuid "github.com/satori/go.uuid"

type GetClientUser struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    *string   `json:"status"`
	Floor     int       `json:"floor"`
}
