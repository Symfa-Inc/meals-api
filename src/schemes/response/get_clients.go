package response

import "go_api/src/domain"

// GetClients response scheme
type GetClients struct {
	Items []Client `json:"items"`
	Page  int      `json:"page"`
	Total int      `json:"total"`
} //@name GetClients

// GetCateringClients struct
type GetCateringClients struct {
	Items []CateringClient `json:"items"`
	Page  int              `json:"page"`
	Total int              `json:"total"`
} //@name GetClients

// Client response struct
type Client struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	domain.UserCatering `json:"catering"`
}

// ClientInfo response struct
type ClientInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CateringClient response struct
type CateringClient struct {
	ClientInfo   `json:"client"`
	OrdersDishes string `json:"ordersDishes" gorm:"column:orders_dishes"`
	Total        int    `json:"total" gorm:"column:total"`
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
	Total        int        `json:"total" gorm:"column:total"`
}
