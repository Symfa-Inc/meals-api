package repository

import (
	"errors"
	"github.com/Aiscom-LLC/meals-api/api/types"
	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/domain"
	"net/http"
)

// AddressRepo struct
type AddressRepo struct{}

// NewAddressRepo returns pointer to
// address repository with all methods
func NewAddressRepo() *AddressRepo {
	return &AddressRepo{}
}

// Add creates new address entity
// returns error or nil
func (a AddressRepo) Add(address domain.Address) (domain.Address, error) {
	if err := config.DB.Create(&address).Error; err != nil {
		return domain.Address{}, err
	}

	return address, nil
}

// Get returns list of addresses of passed client ID
// returns list of addresses and error
func (a AddressRepo) Get(id string) ([]domain.Address, int, error) {
	var addresses []domain.Address

	if isClientExist := config.DB.
		Where("id = ?", id).
		Find(&domain.Client{}).RowsAffected; isClientExist == 0 {

		return nil, http.StatusNotFound, errors.New("client with that ID is not found")
	}

	if err := config.DB.
		Where("client_id = ?", id).
		Find(&addresses).
		Error; err != nil {
		return nil, http.StatusBadRequest, err
	}

	return addresses, 0, nil
}

// Delete soft deletes address, returns err or nil
func (a AddressRepo) Delete(path types.PathAddress) error {
	if isAddressExist := config.DB.
		Where("client_id = ? AND id = ?", path.ID, path.AddressID).
		Delete(&domain.Address{}).
		RowsAffected; isAddressExist == 0 {
		return errors.New("address not found")
	}

	return nil
}

// Update updates entity
// returns error or nil and status code
func (a AddressRepo) Update(path types.PathAddress, address domain.Address) (domain.Address, error) {
	if err := config.DB.Model(&address).
		Where("id = ? AND client_id = ?", path.AddressID, path.ID).
		Update(&address).
		Scan(&address).
		Error; err != nil {
		return address, err
	}

	return address, nil
}
