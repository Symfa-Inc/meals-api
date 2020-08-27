package types

type enum struct {
	SuperAdmin    string
	CateringAdmin string
	ClientAdmin   string
	User          string
}

// UserRoleEnum enum
var UserRoleEnum = enum{
	SuperAdmin:    "Super administrator",
	CateringAdmin: "Catering administrator",
	ClientAdmin:   "Client administrator",
	User:          "User",
}
