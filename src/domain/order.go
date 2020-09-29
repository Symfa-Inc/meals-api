package domain

import "time"

// Order struct for db
type Order struct {
	Base
	Total   *float32
	Status  *string `sql:"type:order_status_types"`
	Comment *string
	Date    time.Time
}
