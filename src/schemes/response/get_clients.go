package response

import "github.com/Aiscom-LLC/meals-api/src/domain"

// GetClients response scheme
type GetClients struct {
	Items []Client `json:"items"`
	Page  int      `json:"page"`
	Total int      `json:"total"`
} //@name GetClients

// GetClientsOrders struct
type GetClientsOrders struct {
	Items []ClientOrder `json:"items"`
	Page  int           `json:"page"`
	Total int           `json:"total"`
} //@name GetClients

// Client response struct
type Client struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	AutoApproveOrders   bool   `json:"autoApproveOrders"`
	domain.UserCatering `json:"catering"`
}

// ClientInfo response struct
type ClientInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ClientOrder response struct
type ClientOrder struct {
	ClientInfo   `json:"client"`
	OrdersDishes string  `json:"ordersDishes" gorm:"column:orders_dishes"`
	Total        float32 `json:"total" gorm:"column:total"`
}

// ClientResponse struct
type ClientResponse struct {
	ID       string              `json:"id"`
	Name     string              `json:"name"`
	Catering domain.UserCatering `json:"catering"`
}

// GetCateringClientsSwagger struct
type GetCateringClientsSwagger struct {
	Items []CateringClientSwagger `json:"items"`
	Page  int                     `json:"page"`
	Total int                     `json:"total"`
}

// CateringClientSwagger response struct
type CateringClientSwagger struct {
	ClientInfo   ClientInfo `json:"client"`
	OrdersDishes string     `json:"ordersDishes" gorm:"column:orders_dishes"`
	Total        float32    `json:"total" gorm:"column:total"`
}
