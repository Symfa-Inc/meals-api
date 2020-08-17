package repository

import (
	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
)

// ClientUserRepo struct
type ClientUserRepo struct{}

// NewClientUserRepo returns pointer to user repository
// with all methods
func NewClientUserRepo() *ClientUserRepo {
	return &ClientUserRepo{}
}

func (cur *ClientUserRepo) Add(clientUser domain.ClientUser) error {
	err := config.DB.
		Create(&clientUser).
		Error
	return err
}
