package response

import uuid "github.com/satori/go.uuid"

// SummaryUserOrder struct
type SummaryUserOrder struct {
	ID      uuid.UUID           `json:"-"`
	Name    string              `json:"name" gorm:"column:full_name"`
	Floor   int                 `json:"floor"`
	Items   []ItemsSummaryOrder `json:"items"`
	Comment string              `json:"comment"`
	Total   int                 `json:"total"`
}

// SummaryOrder struct
type SummaryOrder struct {
	CategorySummaryOrder `json:"category"`
	Items                []ItemsSummaryOrder `json:"items"`
}

// CategorySummaryOrder struct
type CategorySummaryOrder struct {
	ID   uuid.UUID `json:"categoryId" gorm:"column:category_id"`
	Name string    `json:"categoryName" gorm:"column:category_name"`
}

// ItemsSummaryOrder struct
type ItemsSummaryOrder struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

// SummaryOrderResult struct
type SummaryOrderResult struct {
	SummaryOrders []SummaryOrder     `json:"summary"`
	UserOrders    []SummaryUserOrder `json:"userOrders"`
	Total         int                `json:"summaryTotal"`
	Status        *string            `json:"status"`
}

// SummaryOrdersResponse struct
type SummaryOrdersResponse struct {
	Summary      SummaryOrder     `json:"summary"`
	SummaryTotal int              `json:"summaryTotal"`
	UserOrders   SummaryUserOrder `json:"userOrders"`
	Status       string           `json:"status"`
}
