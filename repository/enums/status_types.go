package enums

type statusEnum struct {
	Active    string
	Invited   string
	Deleted   string
	Draft     string
	Published string
}

// StatusTypesEnum enum
var StatusTypesEnum = statusEnum{
	Active:    "active",
	Invited:   "invited",
	Deleted:   "deleted",
	Draft:     "draft",
	Published: "published",
}
