package types

type orderStatusEnum struct {
	Approved string
	Pending  string
	Canceled string
}

// OrderStatusTypesEnum enum
var OrderStatusTypesEnum = orderStatusEnum{
	Approved: "approved",
	Pending:  "pending",
	Canceled: "canceled",
}
