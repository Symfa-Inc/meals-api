package interfaces

import uuid "github.com/satori/go.uuid"

// UserOrders struct for db
type UserOrders struct {
	Base
	UserID  uuid.UUID
	OrderID uuid.UUID
}
