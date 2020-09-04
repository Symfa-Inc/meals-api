package domain

import "time"

// Order struct for db
type Order struct {
	Base
	Total   *int
	Status  *string `sql:"api_types:order_status_types"`
	Comment *string
	Date    time.Time
}
