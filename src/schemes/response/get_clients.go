package response

import "go_api/src/domain"

type GetClients struct {
	Items []domain.Client `json:"items"`
	Page  int             `json:"page"`
	Total int             `json:"total"`
} //@name GetClientsResponse
