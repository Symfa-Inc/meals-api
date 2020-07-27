package response

import "go_api/src/domain"

// GetClients response scheme
type GetClients struct {
	Items []ClientResponse `json:"items"`
	Page  int      `json:"page"`
	Total int      `json:"total"`
} //@name GetClientsResponse

// Client response struct
type Client struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	domain.UserCatering `json:"catering"`
}

type ClientResponse struct {
	ID       string              `json:"id"`
	Name     string              `json:"name"`
	Catering domain.UserCatering `json:"catering"`
}
