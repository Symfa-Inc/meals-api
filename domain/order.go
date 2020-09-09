package domain

import (
	"time"
)

// Order struct for db
type Order struct {
	Base
	Total   *int
	Status  *string `sql:"type:order_status_types"`
	Comment *string
	Date    time.Time
}
