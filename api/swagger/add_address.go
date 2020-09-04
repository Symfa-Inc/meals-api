package swagger

// AddAddress response scheme
type AddAddress struct {
	City   string `json:"city" binding:"required"`
	Street string `json:"street" binding:"required"`
	House  string `json:"house" binding:"required"`
	Floor  int    `json:"floor" binding:"required"`
} //@name AddAddressRequest
