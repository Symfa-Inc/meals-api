package swagger

import (
	"github.com/Aiscom-LLC/meals-api/repository/models"
	uuid "github.com/satori/go.uuid"
)

type GetClientUser struct {
	ID                uuid.UUID `json:"id"`
	models.UserClient `json:"client"`
	FirstName         string  `json:"firstName"`
	LastName          string  `json:"lastName"`
	Email             string  `json:"email"`
	Role              string  `json:"role"`
	Status            *string `json:"status"`
	Floor             int     `json:"floor"`
}
