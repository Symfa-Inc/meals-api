package usecase

import (
	"go_api/src/domain"
	"go_api/src/repository"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Address sturct
type Address struct{}

// NewAddress returns pointer to Address struct
// with all methods
func NewAddress() *Address {
	return &Address{}
}

var addressRepo = repository.NewAddressRepo()

// Get returns list of addresses or error
// @Summary Get list of addresses
// @Tags clients addresses
// @Produce json
// @Param id path string false "Client ID"
// @Success 200 {array} domain.Address "array of address readings"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /clients/{id}/addresses [get]
// Get returns list of addresses of passed client ID
// returns list of addresses and error
func (a Address) Get(c *gin.Context) {
	var path types.PathID

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	addressResult, code, err := addressRepo.Get(path.ID)
	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, addressResult)
}

// Add adds address for client with provided ID
// @Summary Add address for certain client
// @Tags clients addresses
// @Produce json
// @Param id path string true "Client ID"
// @Param payload body request.AddAddress false "address object"
// @Success 200 {object} domain.Address false "address object"
// @Failure 400 {object} types.Error "Error"
// @Router /clients/{id}/addresses [post]
func (a Address) Add(c *gin.Context) {
	var path types.PathID
	var body domain.Address

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	body.ClientID, _ = uuid.FromString(path.ID)
	address, err := addressRepo.Add(body)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, address)
}

// Delete soft delete of addresses reading
// @Summary Soft delete
// @Tags clients addresses
// @Produce json
// @Param id path string true "Client ID"
// @Param addressId path string true "Address ID"
// @Success 204 "Successfully deleted"
// @Failure 404 {object} types.Error "Not Found"
// @Router /clients/{id}/addresses/{addressId} [delete]
func (a Address) Delete(c *gin.Context) {
	var path types.PathAddress
	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := addressRepo.Delete(path); err != nil {
		utils.CreateError(http.StatusNotFound, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}

// Update updates addresses
// @Summary Returns 204 if success and 4xx error if failed
// @Produce json
// @Accept json
// @Tags clients addresses
// @Param id path string true "Client ID"
// @Param addressId path string true "Address ID"
// @Param payload body request.AddAddress false "address object"
// @Success 204 "Successfully updated"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /clients/{id}/addresses/{addressId} [put]
func (a Address) Update(c *gin.Context) {
	var path types.PathAddress
	var body domain.Address

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	address, err := addressRepo.Update(path, body)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, address)
}
