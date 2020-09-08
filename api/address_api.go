package api

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/interfaces"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

// Address struct
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
// @Success 200 {array} swagger.Address "array of address readings"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /clients/{id}/addresses [get]
// Get returns list of addresses of passed client ID
// returns list of addresses and error
func (a Address) Get(c *gin.Context) {
	var path url.PathID

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	addressResult, code, err := addressRepo.Get(path.ID)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, addressResult)
}

// Add adds address for client with provided ID
// @Summary Add address for certain client
// @Tags clients addresses
// @Produce json
// @Param id path string true "Client ID"
// @Param payload body swagger.AddAddress false "address object"
// @Success 200 {object} swagger.Address false "address object"
// @Failure 400 {object} Error "Error"
// @Router /clients/{id}/addresses [post]
func (a Address) Add(c *gin.Context) {
	var path url.PathID
	var body interfaces.Address

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	body.ClientID, _ = uuid.FromString(path.ID)
	address, err := addressRepo.Add(body)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
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
// @Failure 404 {object} Error "Not Found"
// @Router /clients/{id}/addresses/{addressId} [delete]
func (a Address) Delete(c *gin.Context) {
	var path url.PathAddress

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := addressRepo.Delete(path); err != nil {
		utils.CreateError(http.StatusNotFound, err, c)
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
// @Param payload body swagger.AddAddress false "address object"
// @Success 204 "Successfully updated"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /clients/{id}/addresses/{addressId} [put]
func (a Address) Update(c *gin.Context) {
	var path url.PathAddress
	var body interfaces.Address

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	address, err := addressRepo.Update(path, body)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	c.JSON(http.StatusOK, address)
}
