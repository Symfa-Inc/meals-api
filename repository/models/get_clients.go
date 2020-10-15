package models

// ClientInfo response struct
type ClientInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ClientOrder response struct
type ClientOrder struct {
	ClientInfo   `json:"client"`
	OrdersDishes string `json:"ordersDishes" gorm:"column:orders_dishes"`
	Total        int    `json:"total" gorm:"column:total"`
}

// GetClients response scheme
type GetClients struct {
	Items []Client `json:"items"`
	Page  int      `json:"page"`
	Total int      `json:"total"`
} //@name GetClients
