package types

type statusEnum struct {
	Active  string
	Invited string
	Deleted string
}

// StatusTypesEnum enum
var StatusTypesEnum = statusEnum{
	Active:  "active",
	Invited: "invited",
	Deleted: "deleted",
}
